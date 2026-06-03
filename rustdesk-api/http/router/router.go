package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/controller/web"
)

func WebInit(g *gin.Engine) {
	i := &web.Index{}
	g.GET("/", i.Index)

	// 用户中心快捷入口
	g.GET("/portal", func(c *gin.Context) { c.Redirect(302, "/_admin/#/portal") })
	g.GET("/login", func(c *gin.Context) { c.Redirect(302, "/_admin/#/login") })
	g.GET("/register", func(c *gin.Context) { c.Redirect(302, "/_admin/#/register") })

	if global.Config.App.WebClient == 1 {
		g.GET("/webclient-config/index.js", i.ConfigJs)
		g.StaticFS("/webclient", http.Dir(global.Config.Gin.ResourcesPath+"/web"))
		// v2 客户端不存在时回退到 v1
		web2Path := global.Config.Gin.ResourcesPath + "/web2"
		if _, err := os.Stat(web2Path); err == nil {
			g.StaticFS("/webclient2", http.Dir(web2Path))
		} else {
			g.StaticFS("/webclient2", http.Dir(global.Config.Gin.ResourcesPath+"/web"))
		}
	}
	g.StaticFS("/_admin", http.Dir(global.Config.Gin.ResourcesPath+"/admin"))
}
