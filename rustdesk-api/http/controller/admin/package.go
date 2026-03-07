package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"strconv"
)

type Package struct {
}

// List 套餐列表
func (ct *Package) List(c *gin.Context) {
	query := &admin.PageQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.PackageService.List(query.Page, query.PageSize, nil)
	response.Success(c, res)
}

// Detail 套餐详情
func (ct *Package) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	pkg, err := service.AllService.PackageService.GetById(uint(id))
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "DataNotFound"))
		return
	}
	response.Success(c, pkg)
}

// Create 创建套餐
func (ct *Package) Create(c *gin.Context) {
	form := &admin.PackageForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	pkg := &model.Package{
		Name:        form.Name,
		ValidDays:   form.ValidDays,
		DeviceLimit: form.DeviceLimit,
		Description: form.Description,
		Price:       form.Price,
		IsActive:    form.IsActive,
		Priority:    form.Priority,
	}

	if err := service.AllService.PackageService.Create(pkg, form.ServerIds); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, pkg)
}

// Update 更新套餐
func (ct *Package) Update(c *gin.Context) {
	form := &admin.PackageForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	pkg := &model.Package{
		IdModel:     model.IdModel{Id: form.Id},
		Name:        form.Name,
		ValidDays:   form.ValidDays,
		DeviceLimit: form.DeviceLimit,
		Description: form.Description,
		Price:       form.Price,
		IsActive:    form.IsActive,
		Priority:    form.Priority,
	}

	if err := service.AllService.PackageService.Update(pkg, form.ServerIds); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, pkg)
}

// Delete 删除套餐
func (ct *Package) Delete(c *gin.Context) {
	form := &admin.PackageForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	if err := service.AllService.PackageService.Delete(form.Id); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, nil)
}
