package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"strconv"
)

type Server struct {
}

// List 服务器列表
func (ct *Server) List(c *gin.Context) {
	query := &admin.PageQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.ServerService.List(query.Page, query.PageSize, nil)
	response.Success(c, res)
}

// Detail 服务器详情
func (ct *Server) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	server, err := service.AllService.ServerService.GetById(uint(id))
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "DataNotFound"))
		return
	}
	response.Success(c, server)
}

// Create 创建服务器
func (ct *Server) Create(c *gin.Context) {
	form := &admin.ServerForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	server := &model.Server{
		Name:          form.Name,
		Region:        form.Region,
		IdServer:      form.IdServer,
		RelayServer:   form.RelayServer,
		Key:           form.Key,
		ApiServer:     form.ApiServer,
		WsHost:        form.WsHost,
		TopologyGroup: form.TopologyGroup,
		SupportTCP:    form.SupportTCP,
		SupportWSS:    form.SupportWSS,
		CostWeight:    form.CostWeight,
		IsDefault:     form.IsDefault,
		IsActive:      form.IsActive,
		Priority:      form.Priority,
		Description:   form.Description,
	}

	if err := service.AllService.ServerService.Create(server); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, server)
}

// Update 更新服务器
func (ct *Server) Update(c *gin.Context) {
	form := &admin.ServerForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	server := &model.Server{
		IdModel:       model.IdModel{Id: form.Id},
		Name:          form.Name,
		Region:        form.Region,
		IdServer:      form.IdServer,
		RelayServer:   form.RelayServer,
		Key:           form.Key,
		ApiServer:     form.ApiServer,
		WsHost:        form.WsHost,
		TopologyGroup: form.TopologyGroup,
		SupportTCP:    form.SupportTCP,
		SupportWSS:    form.SupportWSS,
		CostWeight:    form.CostWeight,
		IsDefault:     form.IsDefault,
		IsActive:      form.IsActive,
		Priority:      form.Priority,
		Description:   form.Description,
	}

	if err := service.AllService.ServerService.Update(server); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, server)
}

// Delete 删除服务器
func (ct *Server) Delete(c *gin.Context) {
	form := &admin.ServerForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	if err := service.AllService.ServerService.Delete(form.Id); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, nil)
}
