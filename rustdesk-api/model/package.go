package model

// Package stores plan metadata and the available server lines bound to it.
type Package struct {
	IdModel
	Name                string    `json:"name" gorm:"size:100;not null;comment:套餐名称"`
	ValidDays           int       `json:"valid_days" gorm:"default:30;not null;comment:有效天数"`
	DeviceLimit         int       `json:"device_limit" gorm:"default:10;not null;comment:设备数量限制"`
	FileTransferLimitMB int       `json:"file_transfer_limit_mb" gorm:"default:100;not null;comment:文件传输上限(MB)"`
	Description         string    `json:"description" gorm:"type:text;comment:套餐描述"`
	Price               float64   `json:"price" gorm:"default:0;comment:价格"`
	IsActive            bool      `json:"is_active" gorm:"default:true;comment:是否启用"`
	IsDefaultNewUser    bool      `json:"is_default_new_user" gorm:"default:false;comment:新用户注册默认套餐"`
	Priority            int       `json:"priority" gorm:"default:0;comment:优先级"`
	Servers             []*Server `json:"servers" gorm:"many2many:package_servers;"`
	TimeModel
}

// PackageServer stores package-server bindings.
type PackageServer struct {
	PackageId uint `json:"package_id" gorm:"primaryKey;comment:套餐ID"`
	ServerId  uint `json:"server_id" gorm:"primaryKey;comment:服务器ID"`
	IsPrimary bool `json:"is_primary" gorm:"default:false;comment:是否主线路"`
}

// PackageList is the package list response.
type PackageList struct {
	Packages []*Package `json:"list"`
	Pagination
}

func (PackageServer) TableName() string {
	return "package_servers"
}
