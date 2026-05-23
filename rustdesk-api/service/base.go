package service

import "gorm.io/gorm"

// BaseService provides common database access for all services
type BaseService struct {
	db *gorm.DB
}

// NewBaseService creates a new BaseService instance
func NewBaseService() *BaseService {
	return &BaseService{
		db: DB,
	}
}
