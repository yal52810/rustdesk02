package service

import (
	"errors"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"sort"
)

// ServerSchedulerService 服务器智能调度服务
type ServerSchedulerService struct{}

// ScheduleRequest 调度请求
type ScheduleRequest struct {
	ClientIP       string     // 客户端 IP
	ClientRegion   RegionType // 客户端地区
	PreferProtocol string     // 偏好协议 (tcp/wss/auto)
	IsRestricted   bool       // 是否受限环境（只能用 443 端口）
}

// ScheduleResult 调度结果
type ScheduleResult struct {
	PrimaryServer   *model.Server   // 主服务器
	FallbackServers []*model.Server // 备用服务器列表
	Protocol        string          // 推荐协议 (tcp/wss)
}

// Schedule 智能调度服务器
func (s *ServerSchedulerService) Schedule(req *ScheduleRequest) (*ScheduleResult, error) {
	// 1. 获取所有在线且启用的服务器
	servers, err := s.getAvailableServers()
	if err != nil {
		return nil, err
	}

	if len(servers) == 0 {
		return nil, errors.New("no available servers")
	}

	// 2. 根据客户端地区和协议支持过滤服务器
	candidates := s.filterServers(servers, req)
	if len(candidates) == 0 {
		// 如果没有匹配的，使用所有可用服务器
		candidates = servers
	}

	// 3. 根据优先级和成本排序
	s.sortServers(candidates, req)

	// 4. 选择主服务器和备用服务器
	result := &ScheduleResult{
		PrimaryServer:   candidates[0],
		FallbackServers: make([]*model.Server, 0),
		Protocol:        s.determineProtocol(candidates[0], req),
	}

	// 添加备用服务器（最多3个）
	for i := 1; i < len(candidates) && i <= 3; i++ {
		result.FallbackServers = append(result.FallbackServers, candidates[i])
	}

	return result, nil
}

// getAvailableServers 获取所有可用服务器
func (s *ServerSchedulerService) getAvailableServers() ([]*model.Server, error) {
	var servers []*model.Server
	err := DB.Where("is_active = ? AND is_online = ?", true, true).
		Order("priority DESC, cost_weight ASC").
		Find(&servers).Error
	return servers, err
}

// filterServers 根据条件过滤服务器
func (s *ServerSchedulerService) filterServers(servers []*model.Server, req *ScheduleRequest) []*model.Server {
	filtered := make([]*model.Server, 0)

	for _, server := range servers {
		// 检查协议支持
		if req.IsRestricted {
			// 受限环境只能用 WSS
			if !server.SupportWSS || server.WsHost == "" {
				continue
			}
		} else if req.PreferProtocol == "tcp" {
			if !server.SupportTCP {
				continue
			}
		} else if req.PreferProtocol == "wss" {
			if !server.SupportWSS || server.WsHost == "" {
				continue
			}
		}

		// 检查地区匹配
		if s.isRegionMatch(server.Region, req.ClientRegion) {
			filtered = append(filtered, server)
		}
	}

	return filtered
}

// isRegionMatch 检查服务器地区是否匹配客户端地区
func (s *ServerSchedulerService) isRegionMatch(serverRegion string, clientRegion RegionType) bool {
	// 中国大陆客户端优先匹配 CN 开头的服务器
	if clientRegion == RegionCN {
		return len(serverRegion) >= 2 && serverRegion[:2] == "CN"
	}

	// 香港、台湾、澳门客户端匹配对应地区
	if clientRegion == RegionHK {
		return serverRegion == "HK"
	}
	if clientRegion == RegionTW {
		return serverRegion == "TW"
	}
	if clientRegion == RegionMO {
		return serverRegion == "MO"
	}

	// 美国客户端匹配 US 开头的服务器
	if clientRegion == RegionUS {
		return len(serverRegion) >= 2 && serverRegion[:2] == "US"
	}

	// 欧洲客户端匹配 EU 开头的服务器
	if clientRegion == RegionEU {
		return len(serverRegion) >= 2 && serverRegion[:2] == "EU"
	}

	// 其他情况返回 true（允许所有服务器）
	return true
}

// sortServers 对服务器进行排序
func (s *ServerSchedulerService) sortServers(servers []*model.Server, req *ScheduleRequest) {
	sort.Slice(servers, func(i, j int) bool {
		// 1. 优先级高的优先
		if servers[i].Priority != servers[j].Priority {
			return servers[i].Priority > servers[j].Priority
		}

		// 2. 成本低的优先（除非是中国大陆客户端，优先质量）
		if req.ClientRegion == RegionCN {
			// 中国大陆客户端优先选择国内服务器，不考虑成本
			isCN_i := len(servers[i].Region) >= 2 && servers[i].Region[:2] == "CN"
			isCN_j := len(servers[j].Region) >= 2 && servers[j].Region[:2] == "CN"
			if isCN_i != isCN_j {
				return isCN_i
			}
		}

		// 3. 成本权重低的优先
		return servers[i].CostWeight < servers[j].CostWeight
	})
}

// determineProtocol 确定使用的协议
func (s *ServerSchedulerService) determineProtocol(server *model.Server, req *ScheduleRequest) string {
	// 受限环境强制使用 WSS
	if req.IsRestricted {
		return "wss"
	}

	// 用户指定协议
	if req.PreferProtocol == "tcp" && server.SupportTCP {
		return "tcp"
	}
	if req.PreferProtocol == "wss" && server.SupportWSS {
		return "wss"
	}

	// 自动选择：优先 TCP（性能好、成本低）
	if server.SupportTCP {
		return "tcp"
	}
	if server.SupportWSS {
		return "wss"
	}

	return "tcp" // 默认
}

// GetServerConfig 获取服务器配置（用于下发给客户端）
func (s *ServerSchedulerService) GetServerConfig(server *model.Server, protocol string) map[string]interface{} {
	config := map[string]interface{}{
		"id_server":    server.IdServer,
		"relay_server": server.RelayServer,
		"api_server":   server.ApiServer,
		"key":          server.Key,
	}

	// 如果使用 WSS 协议，添加 ws_host
	if protocol == "wss" && server.WsHost != "" {
		config["ws_host"] = server.WsHost
	}

	return config
}
