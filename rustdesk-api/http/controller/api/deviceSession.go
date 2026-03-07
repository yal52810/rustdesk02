package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

type DeviceSessionController struct{}

type ConnectForm struct {
	DeviceId string `json:"device_id" binding:"required"`
}

type HeartbeatForm struct {
	DeviceId string `json:"device_id" binding:"required"`
}

type DisconnectForm struct {
	DeviceId string `json:"device_id" binding:"required"`
}

func (d *DeviceSessionController) Connect(c *gin.Context) {
	var form ConnectForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Error(c, "Invalid parameters: "+err.Error())
		return
	}

	userId := c.GetUint("user_id")
	if userId == 0 {
		response.Error(c, "Unauthorized")
		return
	}

	err := service.AllService.DeviceSessionService.Connect(userId, form.DeviceId)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Connected successfully"})
}

func (d *DeviceSessionController) Heartbeat(c *gin.Context) {
	var form HeartbeatForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Error(c, "Invalid parameters: "+err.Error())
		return
	}

	userId := c.GetUint("user_id")
	if userId == 0 {
		response.Error(c, "Unauthorized")
		return
	}

	err := service.AllService.DeviceSessionService.Heartbeat(userId, form.DeviceId)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Heartbeat updated"})
}

func (d *DeviceSessionController) Disconnect(c *gin.Context) {
	var form DisconnectForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Error(c, "Invalid parameters: "+err.Error())
		return
	}

	userId := c.GetUint("user_id")
	if userId == 0 {
		response.Error(c, "Unauthorized")
		return
	}

	err := service.AllService.DeviceSessionService.Disconnect(userId, form.DeviceId)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Disconnected successfully"})
}

func (d *DeviceSessionController) GetSessions(c *gin.Context) {
	userId := c.GetUint("user_id")
	if userId == 0 {
		response.Error(c, "Unauthorized")
		return
	}

	sessions, err := service.AllService.DeviceSessionService.GetUserSessions(userId)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}
