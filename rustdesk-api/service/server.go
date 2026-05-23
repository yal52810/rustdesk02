package service

import (
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

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
	return s.db.Create(server).Error
}

// Update 更新服务器
func (s *ServerService) Update(server *model.Server) error {
	return s.db.Save(server).Error
}

// Delete 删除服务器
func (s *ServerService) Delete(id uint) error {
	return s.db.Delete(&model.Server{}, id).Error
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

// UpdateOnlineStatus 更新服务器在线状态
func (s *ServerService) UpdateOnlineStatus(id uint, isOnline bool) error {
	return s.db.Model(&model.Server{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_online":     isOnline,
		"last_check_at": gorm.Expr("NOW()"),
	}).Error
}
