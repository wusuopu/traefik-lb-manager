package routes

import (
	"app/controllers/server"
	"app/middlewares"

	"github.com/gin-gonic/gin"
)

func InitServer(r *gin.RouterGroup, e *gin.Engine) {
	const baseUrl = "/:id/servers"

	r.GET(baseUrl + "/", middlewares.CheckWorkspace, server.Index)
	r.POST(baseUrl + "/", middlewares.CheckWorkspace, server.Create)
	r.PUT(baseUrl + "/:serverId", middlewares.CheckWorkspace, server.Update)
	r.DELETE(baseUrl + "/:serverId", middlewares.CheckWorkspace, server.Delete)
}
