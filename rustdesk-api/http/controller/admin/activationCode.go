package admin

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

type ActivationCodeForm struct {
	PackageId   *uint      `json:"package_id"`
	ValidDays   int        `json:"valid_days" binding:"required,min=1"`
	DeviceLimit int        `json:"device_limit" binding:"required,min=1"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Remark      string     `json:"remark"`
}

type BatchActivationCodeForm struct {
	PackageId   *uint      `json:"package_id"`
	Count       int        `json:"count" binding:"required,min=1,max=1000"`
	ValidDays   int        `json:"valid_days" binding:"required,min=1"`
	DeviceLimit int        `json:"device_limit" binding:"required,min=1"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Remark      string     `json:"remark"`
}

type ActivationCodeController struct{}

func (ac *ActivationCodeController) Create(c *gin.Context) {
	var form ActivationCodeForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	code, err := service.AllService.ActivationCodeService.Create(
		form.PackageId,
		form.ValidDays,
		form.DeviceLimit,
		form.ExpiresAt,
		form.Remark,
	)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, code)
}

func (ac *ActivationCodeController) BatchCreate(c *gin.Context) {
	var form BatchActivationCodeForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Error(c, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	codes, err := service.AllService.ActivationCodeService.BatchCreate(
		form.PackageId,
		form.Count,
		form.ValidDays,
		form.DeviceLimit,
		form.ExpiresAt,
		form.Remark,
	)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, codes)
}

func (ac *ActivationCodeController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	codes, total, err := service.AllService.ActivationCodeService.List(page, pageSize)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  codes,
		"total": total,
	})
}

func (ac *ActivationCodeController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		response.Error(c, response.TranslateMsg(c, "ParamsError"))
		return
	}

	err := service.AllService.ActivationCodeService.Delete(uint(id))
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, nil)
}
