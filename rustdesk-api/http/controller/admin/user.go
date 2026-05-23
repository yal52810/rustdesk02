package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	adResp "github.com/lejianwen/rustdesk-api/v2/http/response/admin"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"github.com/lejianwen/rustdesk-api/v2/utils"
	"gorm.io/gorm"
)

type User struct {
}

// Detail 管理员
// @Tags 用户
// @Summary 管理员详情
// @Description 管理员详情
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 500 {object} response.Response
// @Router /admin/user/detail/{id} [get]
// @Security token
func (ct *User) Detail(c *gin.Context) {
	id := c.Param("id")
	iid, _ := strconv.Atoi(id)
	u := service.AllService.UserService.InfoById(uint(iid))
	if u.Id > 0 {
		response.Success(c, u)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
	return
}

// Create 管理员
// @Tags 用户
// @Summary 创建管理员
// @Description 创建管理员
// @Accept  json
// @Produce  json
// @Param body body admin.UserForm true "管理员信息"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 500 {object} response.Response
// @Router /admin/user/create [post]
// @Security token
func (ct *User) Create(c *gin.Context) {
	f := &admin.UserForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := f.ToUser()
	err := service.AllService.UserService.Create(u)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// List 列表
// @Tags 用户
// @Summary 管理员列表
// @Description 管理员列表
// @Accept  json
// @Produce  json
// @Param page query int false "页码"
// @Param page_size query int false "页大小"
// @Param username query int false "账户"
// @Success 200 {object} response.Response{data=model.UserList}
// @Failure 500 {object} response.Response
// @Router /admin/user/list [get]
// @Security token
func (ct *User) List(c *gin.Context) {
	query := &admin.UserQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.UserService.List(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.Username != "" {
			tx.Where("username like ?", "%"+query.Username+"%")
		}
	})
	response.Success(c, res)
}

// Update 编辑
// @Tags 用户
// @Summary 管理员编辑
// @Description 管理员编辑
// @Accept  json
// @Produce  json
// @Param body body admin.UserForm true "用户信息"
// @Success 200 {object} response.Response{data=model.User}
// @Failure 500 {object} response.Response
// @Router /admin/user/update [post]
// @Security token
func (ct *User) Update(c *gin.Context) {
	f := &admin.UserForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if f.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := f.ToUser()
	err := service.AllService.UserService.Update(u)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete 删除
// @Tags 用户
// @Summary 管理员删除
// @Description 管理员编删除
// @Accept  json
// @Produce  json
// @Param body body admin.UserForm true "用户信息"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/user/delete [post]
// @Security token
func (ct *User) Delete(c *gin.Context) {
	f := &admin.UserForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	id := f.Id
	errList := global.Validator.ValidVar(c, id, "required,gt=0")
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := service.AllService.UserService.InfoById(f.Id)
	if u.Id > 0 {
		err := service.AllService.UserService.Delete(u)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// UpdatePassword 修改密码
// @Tags 用户
// @Summary 修改密码
// @Description 修改密码
// @Accept  json
// @Produce  json
// @Param body body admin.UserPasswordForm true "用户信息"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/user/updatePassword [post]
// @Security token
func (ct *User) UpdatePassword(c *gin.Context) {
	f := &admin.UserPasswordForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := service.AllService.UserService.InfoById(f.Id)
	if u.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	err := service.AllService.UserService.UpdatePassword(u, f.Password)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Current 当前用户
// @Tags 用户
// @Summary 当前用户
// @Description 当前用户
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=adResp.LoginPayload}
// @Failure 500 {object} response.Response
// @Router /admin/user/current [get]
// @Security token
func (ct *User) Current(c *gin.Context) {
	u := service.AllService.UserService.CurUser(c)
	token, _ := c.Get("token")
	t := token.(string)
	responseLoginSuccess(c, u, t)
}

// ChangeCurPwd 修改当前用户密码
// @Tags 用户
// @Summary 修改当前用户密码
// @Description 修改当前用户密码
// @Accept  json
// @Produce  json
// @Param body body admin.ChangeCurPasswordForm true "用户信息"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/user/changeCurPwd [post]
// @Security token
func (ct *User) ChangeCurPwd(c *gin.Context) {
	f := &admin.ChangeCurPasswordForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := service.AllService.UserService.CurUser(c)
	// Verify the old password only when the account already has one set
	if !service.AllService.UserService.IsPasswordEmptyByUser(u) {
		ok, _, err := utils.VerifyPassword(u.Password, f.OldPassword)
		if err != nil || !ok {
			response.Fail(c, 101, response.TranslateMsg(c, "OldPasswordError"))
			return
		}
	}
	err := service.AllService.UserService.UpdatePassword(u, f.NewPassword)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// MyOauth
// @Tags 用户
// @Summary 我的授权
// @Description 我的授权
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response{data=[]adResp.UserOauthItem}
// @Failure 500 {object} response.Response
// @Router /admin/user/myOauth [get]
// @Security token
func (ct *User) MyOauth(c *gin.Context) {
	u := service.AllService.UserService.CurUser(c)
	oal := service.AllService.OauthService.List(1, 100, nil)
	ops := make([]string, 0)
	for _, oa := range oal.Oauths {
		ops = append(ops, oa.Op)
	}
	uts := service.AllService.UserService.UserThirdsByUserId(u.Id)
	var res []*adResp.UserOauthItem
	for _, oa := range oal.Oauths {
		item := &adResp.UserOauthItem{
			Op: oa.Op,
		}
		for _, ut := range uts {
			if ut.Op == oa.Op {
				item.Status = 1
				break
			}
		}
		res = append(res, item)
	}
	response.Success(c, res)
}

// groupUsers
func (ct *User) GroupUsers(c *gin.Context) {
	aG := service.AllService.GroupService.List(1, 999, nil)
	aU := service.AllService.UserService.List(1, 9999, nil)
	response.Success(c, gin.H{
		"groups": aG.Groups,
		"users":  aU.Users,
	})
}

// Register
func (ct *User) Register(c *gin.Context) {
	if !global.Config.App.Register {
		response.Fail(c, 101, response.TranslateMsg(c, "RegisterClosed"))
		return
	}
	f := &admin.RegisterForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	regStatus := model.StatusCode(global.Config.App.RegisterStatus)
	// 注册状态可能未配置，默认启用
	if regStatus != model.COMMON_STATUS_DISABLED && regStatus != model.COMMON_STATUS_ENABLE {
		regStatus = model.COMMON_STATUS_ENABLE
	}

	u, err := service.AllService.UserService.RegisterWithActivationCode(f.Username, f.Email, f.Password, regStatus, f.ActivationCode)
	if err != nil {
		response.Fail(c, 101, err.Error())
		return
	}
	if u == nil || u.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed"))
		return
	}
	if regStatus == model.COMMON_STATUS_DISABLED {
		// 需要管理员审核
		response.Fail(c, 101, response.TranslateMsg(c, "RegisterSuccessWaitAdminConfirm"))
		return
	}
	// 注册成功后自动登录
	ut := service.AllService.UserService.Login(u, &model.LoginLog{
		UserId: u.Id,
		Client: model.LoginLogClientWebAdmin,
		Uuid:   "",
		Ip:     c.ClientIP(),
		Type:   model.LoginLogTypeAccount,
	})
	responseLoginSuccess(c, u, ut.Token)
}

// BatchCreate 批量创建用户
// @Tags 用户
// @Summary 批量创建用户
// @Description 批量创建用户
// @Accept  json
// @Produce  json
// @Param body body admin.BatchUserCreateForm true "批量用户信息"
// @Success 200 {object} response.Response{data=gin.H}
// @Failure 500 {object} response.Response
// @Router /admin/user/batchCreate [post]
// @Security token
func (ct *User) BatchCreate(c *gin.Context) {
	f := &admin.BatchUserCreateForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	
	// 转换表单数据为用户对象
	users := make([]*model.User, 0, len(f.Users))
	for _, item := range f.Users {
		user := &model.User{
			Username:    item.Username,
			Email:       item.Email,
			Password:    item.Password,
			Nickname:    item.Nickname,
			GroupId:     item.GroupId,
			Status:      item.Status,
			Remark:      item.Remark,
			ValidDays:   item.ValidDays,
			DeviceLimit: item.DeviceLimit,
			IsAdmin:     &[]bool{false}[0], // 默认非管理员
		}
		users = append(users, user)
	}
	
	// 批量创建用户
	successCount, errors := service.AllService.UserService.BatchCreate(users)
	
	if len(errors) > 0 {
		response.Success(c, gin.H{
			"success_count": successCount,
			"error_count":   len(errors),
			"errors":        errors,
			"message":       "部分用户创建失败",
		})
		return
	}
	
	response.Success(c, gin.H{
		"success_count": successCount,
		"error_count":   0,
		"message":       "所有用户创建成功",
	})
}

func (ct *User) QuickBatchCreate(c *gin.Context) {
	f := &admin.QuickBatchUserCreateForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}

	users := make([]*model.User, 0, f.Count)
	userPasswords := make([]map[string]string, 0, f.Count)

	digits := len(fmt.Sprintf("%d", f.Count))
	formatStr := fmt.Sprintf("%%s%%0%dd", digits)

	for i := 1; i <= f.Count; i++ {
		username := fmt.Sprintf(formatStr, f.UsernamePrefix, i)
		password := utils.RandomNumericPassword(f.PasswordLength)

		user := &model.User{
			Username:    username,
			Password:    password, // 不加密，让 BatchCreate 服务函数处理
			GroupId:     f.GroupId,
			Status:      f.Status,
			ValidDays:   f.ValidDays,
			DeviceLimit: f.DeviceLimit,
			IsAdmin:     &[]bool{false}[0],
		}
		users = append(users, user)
		userPasswords = append(userPasswords, map[string]string{
			"username": username,
			"password": password,
		})
	}

	successCount, errors := service.AllService.UserService.BatchCreate(users)

	var exportText, exportCSV strings.Builder
	exportText.WriteString("账号\t密码\n")
	exportText.WriteString("--------------------\n")
	exportCSV.WriteString("账号,密码\n")

	for _, up := range userPasswords {
		exportText.WriteString(fmt.Sprintf("%s\t%s\n", up["username"], up["password"]))
		exportCSV.WriteString(fmt.Sprintf("%s,%s\n", up["username"], up["password"]))
	}

	if len(errors) > 0 {
		response.Success(c, gin.H{
			"success_count": successCount,
			"error_count":   len(errors),
			"errors":        errors,
			"users":         userPasswords,
			"export_text":   exportText.String(),
			"export_csv":    exportCSV.String(),
			"message":       "部分用户创建失败",
		})
		return
	}

	response.Success(c, gin.H{
		"success_count": successCount,
		"error_count":   0,
		"users":         userPasswords,
		"export_text":   exportText.String(),
		"export_csv":    exportCSV.String(),
		"message":       "所有用户创建成功",
	})
}

