package routes

import (
	"app/controllers/workspace"

	"github.com/gin-gonic/gin"
)

func InitWorkspace(r *gin.RouterGroup) {
	r.GET("/", workspace.Index)
	r.POST("/", workspace.Create)
	r.PUT("/:id", workspace.Update)
	r.DELETE("/:id", workspace.Delete)
}
