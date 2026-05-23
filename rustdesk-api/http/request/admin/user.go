package admin

import (
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
)

type UserForm struct {
	Id       uint   `json:"id"`
	Username string `json:"username" validate:"required,gte=2,lte=32"`
	Email    string `json:"email"` //validate:"required,email" email不强制
	//Password string           `json:"password" validate:"required,gte=4,lte=20"`
	Nickname string           `json:"nickname"`
	Avatar   string           `json:"avatar"`
	GroupId  uint             `json:"group_id" validate:"required"`
	IsAdmin  *bool            `json:"is_admin" `
	Status   model.StatusCode `json:"status" validate:"required,gte=0"`
	Remark   string           `json:"remark"`
	FirstLoginAt *time.Time   `json:"first_login_at"`
	ValidDays    int          `json:"valid_days" validate:"gte=-1"`
	DeviceLimit  int          `json:"device_limit" validate:"gte=0"`
	PackageId    *uint        `json:"package_id"`
	CustomIdServer    string  `json:"custom_id_server"`
	CustomRelayServer string  `json:"custom_relay_server"`
	CustomKey         string  `json:"custom_key"`
}

func (uf *UserForm) FromUser(user *model.User) *UserForm {
	uf.Id = user.Id
	uf.Username = user.Username
	uf.Nickname = user.Nickname
	uf.Email = user.Email
	uf.Avatar = user.Avatar
	uf.GroupId = user.GroupId
	uf.IsAdmin = user.IsAdmin
	uf.Status = user.Status
	uf.Remark = user.Remark
	uf.FirstLoginAt = user.FirstLoginAt
	uf.ValidDays = user.ValidDays
	uf.DeviceLimit = user.DeviceLimit
	uf.PackageId = user.PackageId
	uf.CustomIdServer = user.CustomIdServer
	uf.CustomRelayServer = user.CustomRelayServer
	uf.CustomKey = user.CustomKey
	return uf
}
func (uf *UserForm) ToUser() *model.User {
	user := &model.User{}
	user.Id = uf.Id
	user.Username = uf.Username
	user.Nickname = uf.Nickname
	user.Email = uf.Email
	user.Avatar = uf.Avatar
	user.GroupId = uf.GroupId
	user.IsAdmin = uf.IsAdmin
	user.Status = uf.Status
	user.Remark = uf.Remark
	user.FirstLoginAt = uf.FirstLoginAt
	user.ValidDays = uf.ValidDays
	user.DeviceLimit = uf.DeviceLimit
	user.PackageId = uf.PackageId
	user.CustomIdServer = uf.CustomIdServer
	user.CustomRelayServer = uf.CustomRelayServer
	user.CustomKey = uf.CustomKey
	return user
}

type PageQuery struct {
	Page     uint `form:"page"`
	PageSize uint `form:"page_size"`
}

type UserQuery struct {
	PageQuery
	Username string `form:"username"`
}
type UserPasswordForm struct {
	Id       uint   `json:"id" validate:"required"`
	Password string `json:"password" validate:"required,gte=4,lte=32"`
}

type ChangeCurPasswordForm struct {
	OldPassword string `json:"old_password" validate:"required,gte=4,lte=32"`
	NewPassword string `json:"new_password" validate:"required,gte=4,lte=32"`
}
type GroupUsersQuery struct {
	IsMy   int  `json:"is_my"`
	UserId uint `json:"user_id"`
}

type RegisterForm struct {
	Username        string `json:"username" validate:"required,gte=2,lte=32"`
	Email           string `json:"email"` // validate:"required,email"
	Password        string `json:"password" validate:"required,gte=4,lte=32"`
	ConfirmPassword string `json:"confirm_password" validate:"required,gte=4,lte=32"`
	ActivationCode  string `json:"activation_code"`
}

type UserTokenBatchDeleteForm struct {
	Ids []uint `json:"ids" validate:"required"`
}

type BatchUserCreateForm struct {
	Users []BatchUserItem `json:"users" validate:"required,min=1,max=100"`
}

type BatchUserItem struct {
	Username    string           `json:"username" validate:"required,gte=2,lte=32"`
	Email       string           `json:"email"`
	Password    string           `json:"password" validate:"required,gte=4,lte=32"`
	Nickname    string           `json:"nickname"`
	GroupId     uint             `json:"group_id" validate:"required"`
	Status      model.StatusCode `json:"status" validate:"required,gte=0"`
	Remark      string           `json:"remark"`
	ValidDays   int              `json:"valid_days" validate:"gte=-1"`
	DeviceLimit int              `json:"device_limit" validate:"gte=0"`
}

type QuickBatchUserCreateForm struct {
	Count          int              `json:"count" validate:"required,min=1,max=100"`
	UsernamePrefix string           `json:"username_prefix" validate:"required,gte=2,lte=20"`
	PasswordLength int              `json:"password_length" validate:"required,min=6,max=32"`
	GroupId        uint             `json:"group_id" validate:"required"`
	Status         model.StatusCode `json:"status" validate:"required,gte=0"`
	ValidDays      int              `json:"valid_days" validate:"gte=-1"`
	DeviceLimit    int              `json:"device_limit" validate:"gte=0"`
}

