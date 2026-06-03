package service

import (
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
)

type HealthCheckService struct{}

func (s *HealthCheckService) CheckTCPAddress(address string, defaultPort string) bool {
	host, port := parseServerAddress(address, defaultPort)
	if host == "" {
		return false
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (s *HealthCheckService) StartHealthCheck() {
	go s.CheckAllServers()

	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			s.CheckAllServers()
		}
	}()
}

func (s *HealthCheckService) CheckAllServers() {
	servers, err := AllService.ServerService.GetActiveServers()
	if err != nil {
		Logger.Error("failed to get active servers: ", err)
		return
	}

	for _, server := range servers {
		isOnline := s.CheckServerLine(server)
		if server.IsOnline == isOnline {
			continue
		}
		if err := AllService.ServerService.UpdateOnlineStatus(server.Id, isOnline); err != nil {
			Logger.Error("failed to update server online status: ", err)
			continue
		}
		if isOnline {
			Logger.Infof("server %s is now online", server.Name)
		} else {
			Logger.Warnf("server %s is now offline", server.Name)
		}
	}
}

func (s *HealthCheckService) CheckServerLine(server *model.Server) bool {
	// Try port 21115 (hbbs admin) first, then 21116 (rendezvous) as fallback
	idOK := s.CheckTCPAddress(server.IdServer, "21115") || s.CheckTCPAddress(server.IdServer, "21116")
	relayOK := s.CheckTCPAddress(server.RelayServer, "21117")

	if server.SupportWSS && server.WsHost != "" {
		wssOK := s.CheckTCPAddress(server.WsHost, "443")
		return idOK && (relayOK || wssOK)
	}

	return idOK && relayOK
}

func parseServerAddress(address string, defaultPort string) (string, string) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", ""
	}

	if strings.Contains(address, "://") {
		u, err := url.Parse(address)
		if err == nil {
			host := u.Hostname()
			port := u.Port()
			if port == "" {
				switch strings.ToLower(u.Scheme) {
				case "https", "wss":
					port = "443"
				case "http", "ws":
					port = "80"
				default:
					port = defaultPort
				}
			}
			return host, port
		}
	}

	host, port, err := net.SplitHostPort(address)
	if err == nil {
		return host, port
	}

	return address, defaultPort
}
