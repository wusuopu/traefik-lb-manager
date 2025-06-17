package middleware

import (
	"app/di"
	"app/models"
	"app/schemas"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"gorm.io/gorm"
)


func Index(ctx *gin.Context) {
	var data []models.Middleware

	di.Container.DB.
		Where("workspace_id = ?", ctx.Param("id")).
		Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}

func Create(ctx *gin.Context) {
	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "Category", "Options"})
	data["WorkspaceID"] = ctx.Param("id")

	if value, exists := data["Options"]; exists {
		data["Options"], _ = json.Marshal(value)
	}

	var obj models.Middleware
	err := gconv.Struct(data, &obj)
	if err != nil {
		schemas.MakeErrorResponse(ctx, err, 400)
		return
	}

	results := di.Container.DB.Create(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}

func Update(ctx *gin.Context) {
	var obj models.Middleware
	results := di.Container.DB.First(&obj, ctx.Param("middlewareId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Middleware Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "Options"})

	if value, exists := data["Options"]; exists {
		data["Options"], _ = json.Marshal(value)
	}
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[gstr.CaseSnake(k)] = v
	}

	results = di.Container.DB.Model(&models.Middleware{}).Where("id = ?", ctx.Param("middlewareId")).Updates(payload)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, true, nil)
}

func Delete(ctx *gin.Context) {
	var obj models.Middleware
	results := di.Container.DB.First(&obj, ctx.Param("middlewareId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Middleware Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	di.Container.DB.Delete(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, true, nil)
}
