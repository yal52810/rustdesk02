package model

import (
	"time"
)

type ActivationCode struct {
	Id              uint       `gorm:"primarykey" json:"id"`
	Code            string     `gorm:"uniqueIndex;size:32;not null" json:"code"`
	PackageId       *uint      `gorm:"comment:套餐ID" json:"package_id"`
	Package         *Package   `json:"package,omitempty" gorm:"foreignKey:PackageId"`
	ValidDays       int        `gorm:"default:365;not null" json:"valid_days"`
	DeviceLimit     int        `gorm:"default:10;not null" json:"device_limit"`
	ExpiresAt       *time.Time `gorm:"index" json:"expires_at"`
	UsedBy          uint       `gorm:"index" json:"used_by"`
	UsedAt          *time.Time `json:"used_at"`
	Remark          string     `gorm:"size:255" json:"remark"`
	PrimaryServerId *uint      `gorm:"comment:主线路ID" json:"primary_server_id"`
	BackupServerId  *uint      `gorm:"comment:备用线路ID" json:"backup_server_id"`
	AddDays         int        `gorm:"default:0;not null;comment:增加天数" json:"add_days"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (ActivationCode) TableName() string {
	return "activation_codes"
}
