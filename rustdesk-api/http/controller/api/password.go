package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

type Password struct{}

type SendResetCodeReq struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordReq struct {
	Email       string `json:"email" binding:"required"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (p *Password) SendResetCode(c *gin.Context) {
	var req SendResetCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	if err := service.AllService.PasswordResetService.SendCode(req.Email); err != nil {
		response.Error(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "verification code sent",
	})
}

func (p *Password) ResetByCode(c *gin.Context) {
	var req ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	if err := service.AllService.PasswordResetService.ResetWithCode(
		req.Email,
		req.Code,
		req.NewPassword,
	); err != nil {
		response.Error(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "password reset success",
	})
}
