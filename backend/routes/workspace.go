package routes

import (
	"app/config"
	"app/controllers/workspace"

	"github.com/gin-gonic/gin"
)

func InitWorkspace(r *gin.RouterGroup, e *gin.Engine) {
	r.GET("/", workspace.Index)
	r.POST("/", workspace.Create)
	r.GET("/:id", workspace.Show)
	r.PUT("/:id", workspace.Update)
	r.DELETE("/:id", workspace.Delete)


	r.POST("/:id/traefik.yml", workspace.GenerateTraefikConfig)
	r.PUT("/:id/traefik.yml", workspace.UpdateTraefikConfig)

	e.GET(config.Config.Server.BaseUrl + "/workspaces/:id/traefik.yml", workspace.ShowTraefikConfig)
	e.GET(config.Config.Server.BaseUrl + "/workspaces/:id/version", workspace.ShowConfigVersion)
}
