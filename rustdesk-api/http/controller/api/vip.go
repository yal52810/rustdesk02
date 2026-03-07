package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	apiResp "github.com/lejianwen/rustdesk-api/v2/http/response/api"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

type Vip struct{}

// Servers returns the list of active servers
// @Tags VIP
// @Summary 鑾峰彇鏈嶅姟鍣ㄥ垪琛?
// @Description 鑾峰彇鍙敤鐨勬湇鍔″櫒鑺傜偣鍒楄〃
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ServerList
// @Router /vip/servers [get]
func (v *Vip) Servers(c *gin.Context) {
	servers, err := service.AllService.ServerService.GetActiveServers()
	if err != nil {
		response.Error(c, response.TranslateMsg(c, "SystemError")+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": servers})
}

type RedeemReq struct {
	Code string `json:"code" binding:"required"`
}

// Redeem redeems an activation code for the current user
// @Tags VIP
// @Summary 婵€娲荤爜鍏呭€?
// @Description 浣跨敤婵€娲荤爜寤堕暱鏈嶅姟鏃堕棿
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router /vip/redeem [post]
// @Security token
func (v *Vip) Redeem(c *gin.Context) {
	var req RedeemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	user := service.AllService.UserService.CurUser(c)
	if user == nil {
		response.Error(c, response.TranslateMsg(c, "UserNotFound"))
		return
	}

	// Calculate added days
	ac, err := service.AllService.ActivationCodeService.ValidateAndUse(req.Code, user.Id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	addedDays := ac.ValidDays

	// Update user info
	if user.ValidDays != -1 {
		if ac.ValidDays == -1 {
			user.ValidDays = -1
		} else {
			user.ValidDays += ac.ValidDays
		}
	}

	if ac.DeviceLimit > user.DeviceLimit {
		user.DeviceLimit = ac.DeviceLimit
	}

	if err := service.AllService.UserService.Update(user); err != nil {
		response.Error(c, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"valid_days": addedDays,
		"message":    "Redeem success",
	})
}

type RegisterReq struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email"`
	ActivationCode string `json:"activation_code"`
}

// Register registers a new user
// @Tags VIP
// @Summary 鐢ㄦ埛娉ㄥ唽
// @Description 鐢ㄦ埛娉ㄥ唽锛屽彲閫夋縺娲荤爜
// @Accept  json
// @Produce  json
// @Success 200 {object} apiResp.UserPayload
// @Router /register [post]
func (v *Vip) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	if len(req.Password) < 6 {
		response.Error(c, "Password must be at least 6 characters")
		return
	}

	var user *model.User
	var err error

	if req.ActivationCode != "" {
		user, err = service.AllService.UserService.RegisterWithActivationCode(
			req.Username, req.Email, req.Password, model.COMMON_STATUS_ENABLE, req.ActivationCode,
		)
	} else {
		user = service.AllService.UserService.Register(
			req.Username, req.Email, req.Password, model.COMMON_STATUS_ENABLE,
		)
		if user == nil {
			err = errors.New("Register failed, username might exist")
		}
	}

	if err != nil {
		response.Error(c, err.Error())
		return
	}

	// Auto login? Or just return user info.
	// Client typically expects user info or success.
	// Let's return UserPayload

	if user.Email != "" && service.AllService.MailService != nil && service.AllService.MailService.IsConfigured() {
		if err := service.AllService.MailService.SendRegisterSuccess(user.Email, user.Username); err != nil {
			service.Logger.Warn("send register success mail failed: ", err)
		}
	}
	up := (&apiResp.UserPayload{}).FromUser(user)
	c.JSON(http.StatusOK, up)
}
