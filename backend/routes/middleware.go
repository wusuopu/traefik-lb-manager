package routes

import (
	"app/controllers/middleware"
	"app/middlewares"

	"github.com/gin-gonic/gin"
)

func InitMiddleware(r *gin.RouterGroup, e *gin.Engine) {
	const baseUrl = "/:id/middlewares"

	r.GET(baseUrl + "/", middlewares.CheckWorkspace, middleware.Index)
	r.POST(baseUrl + "/", middlewares.CheckWorkspace, middleware.Create)
	r.PUT(baseUrl + "/:middlewareId", middlewares.CheckWorkspace, middleware.Update)
	r.DELETE(baseUrl + "/:middlewareId", middlewares.CheckWorkspace, middleware.Delete)
}
