package service

import (
	"os"
	"sort"

	"github.com/lejianwen/rustdesk-api/v2/model"
)

type ServerConfigService struct {
	*BaseService
}

type ServerConfigResult struct {
	IdServer    string `json:"id_server"`
	RelayServer string `json:"relay_server"`
	Key         string `json:"key"`
	ApiServer   string `json:"api_server,omitempty"`
	WsHost      string `json:"ws_host,omitempty"`
}

func (s *ServerConfigService) GetServerConfig(user *model.User) (idServer, relayServer, key string) {
	result := s.GetServerConfigSmart(user, "", false)
	return result.IdServer, result.RelayServer, result.Key
}

func (s *ServerConfigService) GetServerConfigSmart(user *model.User, clientIP string, isRestricted bool) *ServerConfigResult {
	if user != nil && user.CustomIdServer != "" {
		return &ServerConfigResult{
			IdServer:    user.CustomIdServer,
			RelayServer: user.CustomRelayServer,
			Key:         user.CustomKey,
		}
	}

	if assigned := s.getAssignedServerConfig(user, isRestricted); assigned != nil {
		return assigned
	}

	if clientIP == "" {
		return s.getDefaultConfig()
	}

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
		return s.getDefaultConfig()
	}

	config := serverToConfig(result.PrimaryServer)
	if result.Protocol == "wss" && result.PrimaryServer.WsHost != "" {
		config.WsHost = result.PrimaryServer.WsHost
	}
	return config
}

func (s *ServerConfigService) getAssignedServerConfig(user *model.User, isRestricted bool) *ServerConfigResult {
	if user == nil {
		return nil
	}

	candidates := make([]*model.Server, 0, 4)
	if user.PrimaryServer != nil {
		candidates = append(candidates, user.PrimaryServer)
	}
	if user.BackupServer != nil {
		candidates = append(candidates, user.BackupServer)
	}
	if user.Package != nil && len(user.Package.Servers) > 0 {
		servers := append([]*model.Server{}, user.Package.Servers...)
		sort.SliceStable(servers, func(i, j int) bool {
			if servers[i].Priority == servers[j].Priority {
				return servers[i].Id < servers[j].Id
			}
			return servers[i].Priority > servers[j].Priority
		})
		candidates = append(candidates, servers...)
	}

	filtered := make([]*model.Server, 0, len(candidates))
	seen := make(map[uint]struct{}, len(candidates))
	for _, server := range candidates {
		if server == nil || !server.IsActive {
			continue
		}
		if _, ok := seen[server.Id]; ok {
			continue
		}
		seen[server.Id] = struct{}{}
		filtered = append(filtered, server)
	}
	if len(filtered) == 0 {
		return nil
	}

	if isRestricted {
		for _, server := range filtered {
			if server.SupportWSS && server.WsHost != "" {
				config := serverToConfig(server)
				config.WsHost = server.WsHost
				return config
			}
		}
	}

	return serverToConfig(filtered[0])
}

func serverToConfig(server *model.Server) *ServerConfigResult {
	if server == nil {
		return nil
	}
	return &ServerConfigResult{
		IdServer:    server.IdServer,
		RelayServer: server.RelayServer,
		Key:         server.Key,
		ApiServer:   server.ApiServer,
		WsHost:      server.WsHost,
	}
}

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
