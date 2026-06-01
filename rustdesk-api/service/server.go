package service

import (
	"context"
	"encoding/json"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

// ServerTopologyEntry 服务器拓扑条目（供 Rust hbbs 读取）
type ServerTopologyEntry struct {
	Name        string `json:"name"`
	Region      string `json:"region"`
	RelayServer string `json:"relay_server"`
	Priority    int    `json:"priority"`
	CostWeight  int    `json:"cost_weight"`
	SupportWSS  bool   `json:"support_wss"`
	WsHost      string `json:"ws_host,omitempty"`
}

// RelayServerEntry API 响应中的中继服务器条目
type RelayServerEntry struct {
	Name        string `json:"name"`
	Region      string `json:"region"`
	RelayServer string `json:"relay_server"`
	WsHost      string `json:"ws_host,omitempty"`
	Priority    int    `json:"priority"`
	CostWeight  int    `json:"cost_weight"`
	SupportWSS  bool   `json:"support_wss"`
}

type ServerService struct {
	*BaseService
}

// List 获取服务器列表
func (s *ServerService) List(page, pageSize uint, where func(tx *gorm.DB)) (res *model.ServerList) {
	res = &model.ServerList{}
	res.Page = int64(page)
	res.PageSize = int64(pageSize)
	tx := s.db.Model(&model.Server{})
	if where != nil {
		where(tx)
	}
	tx.Count(&res.Total)
	tx.Scopes(Paginate(page, pageSize))
	tx.Order("priority DESC, id ASC").Find(&res.Servers)
	return
}

// Create 创建服务器
func (s *ServerService) Create(server *model.Server) error {
	err := s.db.Create(server).Error
	if err == nil {
		s.PublishServerTopologyToRedis()
	}
	return err
}

// Update 更新服务器
func (s *ServerService) Update(server *model.Server) error {
	err := s.db.Save(server).Error
	if err == nil {
		s.PublishServerTopologyToRedis()
	}
	return err
}

// Delete 删除服务器
func (s *ServerService) Delete(id uint) error {
	err := s.db.Delete(&model.Server{}, id).Error
	if err == nil {
		s.PublishServerTopologyToRedis()
	}
	return err
}

// GetById 根据ID获取服务器
func (s *ServerService) GetById(id uint) (*model.Server, error) {
	var server model.Server
	err := s.db.First(&server, id).Error
	return &server, err
}

// GetActiveServers 获取所有启用的服务器
func (s *ServerService) GetActiveServers() ([]*model.Server, error) {
	var servers []*model.Server
	err := s.db.Where("is_active = ?", true).Order("priority DESC, id ASC").Find(&servers).Error
	return servers, err
}

// GetDefaultServer 获取默认服务器
func (s *ServerService) GetDefaultServer() (*model.Server, error) {
	var server model.Server
	err := s.db.Where("is_default = ? AND is_active = ?", true, true).First(&server).Error
	return &server, err
}

// GetRelayServerEntries 获取所有可用中继服务器条目（供 API 下发给客户端）
func (s *ServerService) GetRelayServerEntries() []RelayServerEntry {
	servers, err := s.GetActiveServers()
	if err != nil || len(servers) == 0 {
		return nil
	}
	entries := make([]RelayServerEntry, 0, len(servers))
	for _, server := range servers {
		entries = append(entries, RelayServerEntry{
			Name:        server.Name,
			Region:      server.Region,
			RelayServer: server.RelayServer,
			WsHost:      server.WsHost,
			Priority:    server.Priority,
			CostWeight:  server.CostWeight,
			SupportWSS:  server.SupportWSS,
		})
	}
	return entries
}

// UpdateOnlineStatus 更新服务器在线状态
func (s *ServerService) UpdateOnlineStatus(id uint, isOnline bool) error {
	err := s.db.Model(&model.Server{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_online":     isOnline,
		"last_check_at": gorm.Expr("NOW()"),
	}).Error
	if err == nil {
		s.PublishServerTopologyToRedis()
	}
	return err
}

// PublishServerTopologyToRedis 将服务器拓扑发布到 Redis，供 Rust hbbs 读取
func (s *ServerService) PublishServerTopologyToRedis() {
	if global.Redis == nil {
		return
	}

	servers, err := s.GetActiveServers()
	if err != nil {
		return
	}

	entries := make([]ServerTopologyEntry, 0, len(servers))
	for _, server := range servers {
		entries = append(entries, ServerTopologyEntry{
			Name:        server.Name,
			Region:      server.Region,
			RelayServer: server.RelayServer,
			Priority:    server.Priority,
			CostWeight:  server.CostWeight,
			SupportWSS:  server.SupportWSS,
			WsHost:      server.WsHost,
		})
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return
	}

	global.Redis.Set(context.Background(), "rustdesk:server_topology", string(data), 0)
}
