package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/config"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"io"
	"os"
	"strings"
)

type Config struct {
}

type MailConfigPayload struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	From        string `json:"from"`
	FromName    string `json:"from_name"`
	UseSSL      bool   `json:"use_ssl"`
	SkipVerify  bool   `json:"skip_verify"`
	PasswordSet bool   `json:"password_set"`
}

type UpdateMailConfigReq struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	From       string `json:"from" binding:"required"`
	FromName   string `json:"from_name"`
	UseSSL     bool   `json:"use_ssl"`
	SkipVerify bool   `json:"skip_verify"`
}

type TestMailConfigReq struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	From       string `json:"from"`
	FromName   string `json:"from_name"`
	UseSSL     bool   `json:"use_ssl"`
	SkipVerify bool   `json:"skip_verify"`
	To         string `json:"to"`
}

// ServerConfig RUSTDESK服务配置
// @Tags ADMIN
// @Summary RUSTDESK服务配置
// @Description 服务配置,给webclient提供api-server
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config/server [get]
// @Security token
func (co *Config) ServerConfig(c *gin.Context) {
	cf := &response.ServerConfigResponse{
		IdServer:    global.Config.Rustdesk.IdServer,
		Key:         global.Config.Rustdesk.Key,
		RelayServer: global.Config.Rustdesk.RelayServer,
		ApiServer:   global.Config.Rustdesk.ApiServer,
		WsHost:      global.Config.Rustdesk.WsHost,
			CardShopUrl: global.Config.Rustdesk.CardShopUrl,
	}
	response.Success(c, cf)
}

type updateServerKeyReq struct {
	IdServer    string `json:"id_server"`
	RelayServer string `json:"relay_server"`
	ApiServer   string `json:"api_server"`
	Key         string `json:"key"`
	WsHost      string `json:"ws_host"`
	CardShopUrl string `json:"card_shop_url"`
}

// UpdateServerKey 更新服务器全局配置
func (co *Config) UpdateServerKey(c *gin.Context) {
	req := &updateServerKeyReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	if global.Viper == nil {
		response.Fail(c, 101, "config writer not initialized")
		return
	}

	if req.IdServer != "" {
		global.Viper.Set("rustdesk.id-server", req.IdServer)
		global.Config.Rustdesk.IdServer = req.IdServer
		if service.Config != nil {
			service.Config.Rustdesk.IdServer = req.IdServer
		}
	}
	if req.RelayServer != "" {
		global.Viper.Set("rustdesk.relay-server", req.RelayServer)
		global.Config.Rustdesk.RelayServer = req.RelayServer
		if service.Config != nil {
			service.Config.Rustdesk.RelayServer = req.RelayServer
		}
	}
	if req.ApiServer != "" {
		global.Viper.Set("rustdesk.api-server", req.ApiServer)
		global.Config.Rustdesk.ApiServer = req.ApiServer
		if service.Config != nil {
			service.Config.Rustdesk.ApiServer = req.ApiServer
		}
	}
	if req.Key != "" {
		global.Viper.Set("rustdesk.key", req.Key)
		global.Config.Rustdesk.Key = req.Key
		if service.Config != nil {
			service.Config.Rustdesk.Key = req.Key
		}
	}
	if req.WsHost != "" {
		global.Viper.Set("rustdesk.ws-host", req.WsHost)
		global.Config.Rustdesk.WsHost = req.WsHost
		if service.Config != nil {
			service.Config.Rustdesk.WsHost = req.WsHost
		}
	}
	// Card shop URL can be cleared (empty is valid)
	global.Viper.Set("rustdesk.card-shop-url", req.CardShopUrl)
	global.Config.Rustdesk.CardShopUrl = req.CardShopUrl
	if service.Config != nil {
		service.Config.Rustdesk.CardShopUrl = req.CardShopUrl
	}

	if err := global.Viper.WriteConfig(); err != nil {
		service.Logger.Warn("failed to persist config to file (may be read-only fs): ", err)
	}

	response.Success(c, gin.H{
		"id_server":     global.Config.Rustdesk.IdServer,
		"relay_server":  global.Config.Rustdesk.RelayServer,
		"api_server":    global.Config.Rustdesk.ApiServer,
		"key":           global.Config.Rustdesk.Key,
		"ws_host":       global.Config.Rustdesk.WsHost,
		"card_shop_url": global.Config.Rustdesk.CardShopUrl,
	})
}

// AppConfig APP服务配置
// @Tags ADMIN
// @Summary APP服务配置
// @Description APP服务配置
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config/app [get]
// @Security token
func (co *Config) AppConfig(c *gin.Context) {
	response.Success(c, &gin.H{
		"web_client": global.Config.App.WebClient,
	})
}

// AdminConfig ADMIN服务配置
// @Tags ADMIN
// @Summary ADMIN服务配置
// @Description ADMIN服务配置
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config/admin [get]
// @Security token
func (co *Config) AdminConfig(c *gin.Context) {

	u := &model.User{}
	token := c.GetHeader("api-token")
	if token != "" {
		u, _ = service.AllService.UserService.InfoByAccessToken(token)
		if !service.AllService.UserService.CheckUserEnable(u) {
			u.Id = 0
		}
	}

	if u.Id == 0 {
		response.Success(c, &gin.H{
			"title": global.Config.Admin.Title,
		})
		return
	}

	hello := global.Config.Admin.Hello
	if hello == "" {
		helloFile := global.Config.Admin.HelloFile
		if helloFile != "" {
			b, err := os.ReadFile(helloFile)
			if err == nil && len(b) > 0 {
				hello = string(b)
			}
		}
	}

	//replace {{username}} to username
	hello = strings.Replace(hello, "{{username}}", u.Username, -1)
	response.Success(c, &gin.H{
		"title": global.Config.Admin.Title,
		"hello": hello,
	})
}

func (co *Config) MailConfig(c *gin.Context) {
	response.Success(c, &MailConfigPayload{
		Host:        global.Config.Mail.Host,
		Port:        global.Config.Mail.Port,
		Username:    global.Config.Mail.Username,
		From:        global.Config.Mail.From,
		FromName:    global.Config.Mail.FromName,
		UseSSL:      global.Config.Mail.UseSSL,
		SkipVerify:  global.Config.Mail.SkipVerify,
		PasswordSet: strings.TrimSpace(global.Config.Mail.Password) != "",
	})
}

func (co *Config) UpdateMailConfig(c *gin.Context) {
	req := &UpdateMailConfigReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	newCfg := global.Config.Mail
	newCfg.Host = strings.TrimSpace(req.Host)
	newCfg.Port = req.Port
	newCfg.Username = strings.TrimSpace(req.Username)
	newCfg.From = strings.TrimSpace(req.From)
	newCfg.FromName = strings.TrimSpace(req.FromName)
	newCfg.UseSSL = req.UseSSL
	newCfg.SkipVerify = req.SkipVerify
	if strings.TrimSpace(req.Password) != "" {
		newCfg.Password = req.Password
	}

	if err := persistMailConfig(newCfg); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	global.Config.Mail = newCfg
	if service.Config != nil {
		service.Config.Mail = newCfg
	}

	response.Success(c, &MailConfigPayload{
		Host:        newCfg.Host,
		Port:        newCfg.Port,
		Username:    newCfg.Username,
		From:        newCfg.From,
		FromName:    newCfg.FromName,
		UseSSL:      newCfg.UseSSL,
		SkipVerify:  newCfg.SkipVerify,
		PasswordSet: strings.TrimSpace(newCfg.Password) != "",
	})
}

func (co *Config) TestMailConfig(c *gin.Context) {
	req := &TestMailConfigReq{}
	if err := c.ShouldBindJSON(req); err != nil && !errors.Is(err, io.EOF) {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}

	mailCfg := global.Config.Mail
	if host := strings.TrimSpace(req.Host); host != "" {
		mailCfg.Host = host
	}
	if req.Port > 0 {
		mailCfg.Port = req.Port
	}
	if username := strings.TrimSpace(req.Username); username != "" {
		mailCfg.Username = username
	}
	if password := strings.TrimSpace(req.Password); password != "" {
		mailCfg.Password = password
	}
	if from := strings.TrimSpace(req.From); from != "" {
		mailCfg.From = from
	}
	if fromName := strings.TrimSpace(req.FromName); fromName != "" {
		mailCfg.FromName = fromName
	}
	mailCfg.UseSSL = req.UseSSL
	mailCfg.SkipVerify = req.SkipVerify

	to := strings.TrimSpace(req.To)
	if to == "" {
		to = strings.TrimSpace(mailCfg.From)
	}
	if to == "" {
		response.Fail(c, 101, "please provide a sender email or test recipient")
		return
	}

	if err := service.AllService.MailService.SendWithConfig(mailCfg, to, "RustDesk mail test", "This is a test email sent from the RustDesk admin panel."); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Test email sent"})
}

func persistMailConfig(mailCfg config.Mail) error {
	if global.Viper == nil {
		return errors.New("config writer not initialized")
	}

	global.Viper.Set("mail.host", mailCfg.Host)
	global.Viper.Set("mail.port", mailCfg.Port)
	global.Viper.Set("mail.username", mailCfg.Username)
	global.Viper.Set("mail.password", mailCfg.Password)
	global.Viper.Set("mail.from", mailCfg.From)
	global.Viper.Set("mail.from-name", mailCfg.FromName)
	global.Viper.Set("mail.use-ssl", mailCfg.UseSSL)
	global.Viper.Set("mail.skip-verify", mailCfg.SkipVerify)

	if err := global.Viper.WriteConfig(); err != nil {
		service.Logger.Warn("failed to persist mail config to file (may be read-only fs): ", err)
	}
	return nil
}
