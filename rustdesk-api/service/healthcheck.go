package service

import (
	"net"
	"time"
)

type HealthCheckService struct{}

// CheckServerHealth 检查服务器健康状态
func (s *HealthCheckService) CheckServerHealth(host string, port string) bool {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// StartHealthCheck 启动健康检查定时任务
func (s *HealthCheckService) StartHealthCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			s.CheckAllServers()
		}
	}()
}

// CheckAllServers 检查所有服务器
func (s *HealthCheckService) CheckAllServers() {
	servers, err := AllService.ServerService.GetActiveServers()
	if err != nil {
		Logger.Error("Failed to get active servers:", err)
		return
	}

	for _, server := range servers {
		// 从 id_server 或 relay_server 提取主机和端口
		host, port := parseServerAddress(server.IdServer)
		if host == "" {
			host, port = parseServerAddress(server.RelayServer)
		}

		if host == "" {
			continue
		}

		isOnline := s.CheckServerHealth(host, port)
		if server.IsOnline != isOnline {
			err := AllService.ServerService.UpdateOnlineStatus(server.Id, isOnline)
			if err != nil {
				Logger.Error("Failed to update server status:", err)
			} else {
				if isOnline {
					Logger.Infof("Server %s is now online", server.Name)
				} else {
					Logger.Warnf("Server %s is now offline", server.Name)
				}
			}
		}
	}
}

// parseServerAddress 解析服务器地址，提取主机和端口
func parseServerAddress(address string) (string, string) {
	if address == "" {
		return "", ""
	}

	// 简单解析，假设格式为 host:port 或 host
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		// 如果没有端口，使用默认端口
		return address, "21116"
	}
	return host, port
}
