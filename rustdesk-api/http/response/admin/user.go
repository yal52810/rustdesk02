package admin

import "github.com/lejianwen/rustdesk-api/v2/model"

type LoginPayload struct {
	Username   string          `json:"username"`
	Email      string          `json:"email"`
	Avatar     string          `json:"avatar"`
	Token      string          `json:"token"`
	RouteNames []string        `json:"route_names"`
	Nickname   string          `json:"nickname"`
	ValidDays  int             `json:"valid_days"`
	DeviceLimit int            `json:"device_limit"`
	ExpiredAt  *string         `json:"expired_at"`
	Package    *model.Package  `json:"package,omitempty"`
}

func (lp *LoginPayload) FromUser(user *model.User) {
	lp.Username = user.Username
	lp.Email = user.Email
	lp.Avatar = user.Avatar
	lp.Nickname = user.Nickname
	lp.ValidDays = user.ValidDays
	lp.DeviceLimit = user.DeviceLimit
	// 计算到期时间
	if user.FirstLoginAt != nil && user.ValidDays > 0 {
		expiredAt := user.FirstLoginAt.AddDate(0, 0, user.ValidDays)
		expiredStr := expiredAt.Format("2006-01-02 15:04:05")
		lp.ExpiredAt = &expiredStr
	}
	lp.Package = user.Package
}

type UserOauthItem struct {
	Op     string `json:"op"`
	Status int    `json:"status"`
}

