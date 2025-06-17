package routes

import (
	"app/controllers/service"
	"app/middlewares"

	"github.com/gin-gonic/gin"
)

func InitService(r *gin.RouterGroup, e *gin.Engine) {
	const baseUrl = "/:id/services"

	r.GET(baseUrl + "/", middlewares.CheckWorkspace, service.Index)
	r.POST(baseUrl + "/", middlewares.CheckWorkspace, service.Create)
	r.GET(baseUrl + "/external", service.ExternalIndex)
	r.PUT(baseUrl + "/:serviceId", middlewares.CheckWorkspace, service.Update)
	r.DELETE(baseUrl + "/:serviceId", middlewares.CheckWorkspace, service.Delete)
}
