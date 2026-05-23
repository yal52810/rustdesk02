package model

import "time"

type User struct {
	IdModel
	Username string `json:"username" gorm:"default:'';not null;uniqueIndex"`
	Email    string `json:"email" gorm:"default:'';not null;index"`
	// Email	string     	`json:"email" `
	Password          string     `json:"-" gorm:"default:'';not null;"`
	Nickname          string     `json:"nickname" gorm:"default:'';not null;"`
	Avatar            string     `json:"avatar" gorm:"default:'';not null;"`
	GroupId           uint       `json:"group_id" gorm:"default:0;not null;index"`
	IsAdmin           *bool      `json:"is_admin" gorm:"default:0;not null;"`
	Status            StatusCode `json:"status" gorm:"default:1;not null;"`
	Remark            string     `json:"remark" gorm:"default:'';not null;"`
	FirstLoginAt      *time.Time `json:"first_login_at" gorm:"type:timestamp;default:null"`
	ValidDays         int        `json:"valid_days" gorm:"default:0;not null"`
	DeviceLimit       int        `json:"device_limit" gorm:"default:10;not null;"`
	CustomIdServer    string     `json:"custom_id_server" gorm:"default:'';not null;"`
	CustomRelayServer string     `json:"custom_relay_server" gorm:"default:'';not null;"`
	CustomKey         string     `json:"custom_key" gorm:"default:'';not null;"`
	RelayServerId     *uint      `json:"relay_server_id" gorm:"comment:绑定的中继服务器ID"`
	PackageId         *uint      `json:"package_id" gorm:"comment:套餐ID"`
	Package           *Package   `json:"package,omitempty" gorm:"foreignKey:PackageId"`
	PrimaryServerId   *uint      `json:"primary_server_id" gorm:"comment:主线路ID"`
	BackupServerId    *uint      `json:"backup_server_id" gorm:"comment:备用线路ID"`
	PrimaryServer     *Server    `json:"primary_server,omitempty" gorm:"foreignKey:PrimaryServerId"`
	BackupServer      *Server    `json:"backup_server,omitempty" gorm:"foreignKey:BackupServerId"`
	TimeModel
}

// BeforeSave 钩子用于确保 email 字段有合理的默认值
//func (u *User) BeforeSave(tx *gorm.DB) (err error) {
//	// 如果 email 为空，设置为默认值
//	if u.Email == "" {
//		u.Email = fmt.Sprintf("%s@example.com", u.Username)
//	}
//	return nil
//}

type UserList struct {
	Users []*User `json:"list,omitempty"`
	Pagination
}

var UserRouteNames = []string{
	"MyTagList", "MyAddressBookList", "MyInfo", "MyAddressBookCollection", "MyPeer", "MyShareRecordList", "MyLoginLog",
}
var AdminRouteNames = []string{"*"}
