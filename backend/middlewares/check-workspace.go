package middlewares

import (
	"app/di"
	"app/models"
	"app/schemas"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CheckWorkspace(ctx *gin.Context) {
	var workspace models.Workspace
	results := di.Container.DB.First(&workspace, ctx.Param("id"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Workspace Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	ctx.Next()
}