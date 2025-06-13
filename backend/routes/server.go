package routes

import (
	"app/controllers/server"

	"github.com/gin-gonic/gin"
)

func InitServer(r *gin.RouterGroup, e *gin.Engine) {
	const baseUrl = "/:id/servers"

	r.GET(baseUrl + "/", server.Index)
	r.POST(baseUrl + "/", server.Create)
	r.PUT(baseUrl + "/:serverId", server.Update)
	r.DELETE(baseUrl + "/:serverId", server.Delete)
}
