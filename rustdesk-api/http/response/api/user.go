package api

import "github.com/lejianwen/rustdesk-api/v2/model"

/*
	pub enum UserStatus {
	    Disabled = 0,
	    Normal = 1,
	    Unverified = -1,
	}
*/

/*
UserPayload
String name = ”;
String email = ”;
String note = ”;
UserStatus status;
bool isAdmin = false;
*/
type UserPayload struct {
	Name    string                 `json:"name"`
	Email   string                 `json:"email"`
	Note    string                 `json:"note"`
	IsAdmin *bool                  `json:"is_admin"`
	Status  int                    `json:"status"`
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
	up.Info = map[string]interface{}{}
	return up
}

/*
	class HttpType {
	  static const kAuthReqTypeAccount = "account";
	  static const kAuthReqTypeMobile = "mobile";
	  static const kAuthReqTypeSMSCode = "sms_code";
	  static const kAuthReqTypeEmailCode = "email_code";
	  static const kAuthReqTypeTfaCode = "tfa_code";

	  static const kAuthResTypeToken = "access_token";
	  static const kAuthResTypeEmailCheck = "email_check";
	  static const kAuthResTypeTfaCheck = "tfa_check";
	}
*/
type LoginRes struct {
	Type        string      `json:"type"`
	AccessToken string      `json:"access_token"`
	User        UserPayload `json:"user"`
	Secret      string      `json:"secret,omitempty"`
	TfaType     string      `json:"tfa_type,omitempty"`
	IdServer    string      `json:"id_server"`
	RelayServer string      `json:"relay_server"`
	Key         string      `json:"key"`
}
