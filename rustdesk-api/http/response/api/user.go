package api

import "github.com/lejianwen/rustdesk-api/v2/model"

type UserPayload struct {
	Name         string                 `json:"name"`
	Email        string                 `json:"email"`
	Note         string                 `json:"note"`
	IsAdmin      *bool                  `json:"is_admin"`
	Status       int                    `json:"status"`
	ValidDays    int                    `json:"valid_days"`
	FirstLoginAt string                 `json:"first_login_at"`
	DeviceLimit  int                    `json:"device_limit"`
	Info         map[string]interface{} `json:"info"`
}

func (up *UserPayload) FromUser(user *model.User) *UserPayload {
	up.Name = user.Username
	up.Email = user.Email
	up.IsAdmin = user.IsAdmin
	up.Status = int(user.Status)
	up.ValidDays = user.ValidDays
	if user.FirstLoginAt != nil {
		up.FirstLoginAt = user.FirstLoginAt.Format("2006-01-02 15:04:05")
	}
	up.DeviceLimit = user.DeviceLimit
	up.Info = map[string]interface{}{
		"package_id":             user.PackageId,
		"primary_server_id":      user.PrimaryServerId,
		"backup_server_id":       user.BackupServerId,
		"relay_server_id":        user.RelayServerId,
		"file_transfer_limit_mb": 100,
	}
	if user.Package != nil {
		up.Info["package_name"] = user.Package.Name
		if user.Package.FileTransferLimitMB > 0 {
			up.Info["file_transfer_limit_mb"] = user.Package.FileTransferLimitMB
		}
	}
	if user.PrimaryServer != nil {
		up.Info["primary_server_name"] = user.PrimaryServer.Name
		up.Info["primary_server_region"] = user.PrimaryServer.Region
		up.Info["primary_server_supports_websocket"] = user.PrimaryServer.SupportWSS
	}
	if user.BackupServer != nil {
		up.Info["backup_server_name"] = user.BackupServer.Name
		up.Info["backup_server_region"] = user.BackupServer.Region
		up.Info["backup_server_supports_websocket"] = user.BackupServer.SupportWSS
	}
	return up
}

type LoginRes struct {
	Type        string      `json:"type"`
	AccessToken string      `json:"access_token"`
	User        UserPayload `json:"user"`
	Secret      string      `json:"secret,omitempty"`
	TfaType     string      `json:"tfa_type,omitempty"`
	IdServer    string      `json:"id_server"`
	RelayServer string      `json:"relay_server"`
	Key         string      `json:"key"`
	ApiServer   string      `json:"api_server,omitempty"`
	WsHost      string      `json:"ws_host,omitempty"`
}
