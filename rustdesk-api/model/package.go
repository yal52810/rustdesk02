package model

// Package 套餐模型
type Package struct {
	IdModel
	Name        string    `json:"name" gorm:"size:100;not null;comment:套餐名称"`
	ValidDays   int       `json:"valid_days" gorm:"default:30;not null;comment:有效天数"`
	DeviceLimit int       `json:"device_limit" gorm:"default:10;not null;comment:设备数量限制"`
	Description string    `json:"description" gorm:"type:text;comment:套餐描述"`
	Price       float64   `json:"price" gorm:"default:0;comment:价格"`
	IsActive    bool      `json:"is_active" gorm:"default:true;comment:是否启用"`
	Priority    int       `json:"priority" gorm:"default:0;comment:优先级"`
	Servers     []*Server `json:"servers" gorm:"many2many:package_servers;"`
	TimeModel
}

// PackageServer 套餐-服务器关联表
type PackageServer struct {
	PackageId uint `json:"package_id" gorm:"primaryKey;comment:套餐ID"`
	ServerId  uint `json:"server_id" gorm:"primaryKey;comment:服务器ID"`
	IsPrimary bool `json:"is_primary" gorm:"default:false;comment:是否主线路"`
}

// PackageList 套餐列表响应
type PackageList struct {
	Packages []*Package `json:"list"`
	Pagination
}

func (PackageServer) TableName() string {
	return "package_servers"
}
