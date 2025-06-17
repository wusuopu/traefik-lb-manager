package routes

import (
	"app/config"
	"app/middlewares"
	"embed"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup, engine *gin.Engine, embededFiles embed.FS) {
	engine.GET(config.Config.Server.BaseUrl + "/", middlewares.BasicAuthMiddleware, func(ctx *gin.Context) {
		data, _ := embededFiles.ReadFile("assets/index.html")
		ctx.Data(200, "text/html", data)
	})

	engine.GET("_health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	// router.GET("/dashboard", dashboard.Index)

	workspaceGroup := router.Group("/workspaces")
	InitWorkspace(workspaceGroup, engine)
	InitService(workspaceGroup, engine)
	InitMiddleware(workspaceGroup, engine)
	InitCertificate(workspaceGroup, engine)
	InitServer(workspaceGroup, engine)
	InitAcme(engine)
}