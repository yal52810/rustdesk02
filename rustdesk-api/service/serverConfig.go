package service

import (
	"os"

	"github.com/lejianwen/rustdesk-api/v2/model"
)

type ServerConfigService struct {
	*BaseService
}

// ServerConfigResult 服务器配置结果
type ServerConfigResult struct {
	IdServer    string `json:"id_server"`
	RelayServer string `json:"relay_server"`
	Key         string `json:"key"`
	ApiServer   string `json:"api_server,omitempty"`
	WsHost      string `json:"ws_host,omitempty"`
}

// GetServerConfig 获取服务器配置（兼容旧接口）
func (s *ServerConfigService) GetServerConfig(user *model.User) (idServer, relayServer, key string) {
	result := s.GetServerConfigSmart(user, "", false)
	return result.IdServer, result.RelayServer, result.Key
}

// GetServerConfigSmart 智能获取服务器配置
func (s *ServerConfigService) GetServerConfigSmart(user *model.User, clientIP string, isRestricted bool) *ServerConfigResult {
	// 1. 检查用户是否有自定义配置
	if user.CustomIdServer != "" {
		return &ServerConfigResult{
			IdServer:    user.CustomIdServer,
			RelayServer: user.CustomRelayServer,
			Key:         user.CustomKey,
		}
	}

	// 2. 如果没有提供客户端 IP，使用默认配置
	if clientIP == "" {
		return s.getDefaultConfig()
	}

	// 3. 使用智能调度
	geoIPService := &GeoIPService{}
	clientRegion := geoIPService.GetRegionByIP(clientIP)

	schedulerService := &ServerSchedulerService{}
	scheduleReq := &ScheduleRequest{
		ClientIP:       clientIP,
		ClientRegion:   clientRegion,
		PreferProtocol: "auto",
		IsRestricted:   isRestricted,
	}

	result, err := schedulerService.Schedule(scheduleReq)
	if err != nil || result.PrimaryServer == nil {
		// 调度失败，使用默认配置
		return s.getDefaultConfig()
	}

	// 4. 返回调度结果
	server := result.PrimaryServer
	config := &ServerConfigResult{
		IdServer:    server.IdServer,
		RelayServer: server.RelayServer,
		Key:         server.Key,
		ApiServer:   server.ApiServer,
	}

	// 如果推荐使用 WSS，添加 ws_host
	if result.Protocol == "wss" && server.WsHost != "" {
		config.WsHost = server.WsHost
	}

	return config
}

// getDefaultConfig 获取默认配置
func (s *ServerConfigService) getDefaultConfig() *ServerConfigResult {
	idServer := os.Getenv("RUSTDESK_ID_SERVER")
	if idServer == "" {
		idServer = Config.Rustdesk.IdServer
	}

	relayServer := os.Getenv("RUSTDESK_RELAY_SERVER")
	if relayServer == "" {
		relayServer = Config.Rustdesk.RelayServer
	}

	key := os.Getenv("RUSTDESK_KEY")
	if key == "" {
		key = Config.Rustdesk.Key
	}

	apiServer := os.Getenv("RUSTDESK_API_SERVER")
	if apiServer == "" {
		apiServer = Config.Rustdesk.ApiServer
	}

	wsHost := os.Getenv("RUSTDESK_WS_HOST")
	if wsHost == "" {
		wsHost = Config.Rustdesk.WsHost
	}

	return &ServerConfigResult{
		IdServer:    idServer,
		RelayServer: relayServer,
		Key:         key,
		ApiServer:   apiServer,
		WsHost:      wsHost,
	}
}
