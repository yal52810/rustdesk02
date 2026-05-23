package model

import (
	"time"
)

type DeviceSession struct {
	Id           uint      `gorm:"primarykey" json:"id"`
	UserId       uint      `gorm:"index;not null" json:"user_id"`
	DeviceId     string    `gorm:"index;size:64;not null" json:"device_id"`
	Status       int       `gorm:"default:1;not null" json:"status"`
	LastActiveAt time.Time `gorm:"index;not null" json:"last_active_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (DeviceSession) TableName() string {
	return "device_sessions"
}

const (
	DeviceSessionStatusOffline = 0
	DeviceSessionStatusOnline  = 1
)
