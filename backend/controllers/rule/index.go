package rule

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
	var data []models.Rule

	di.Container.DB.
		Where("workspace_id = ?", ctx.Param("id")).
		Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}

func Create(ctx *gin.Context) {
	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Options", "Enable", "ServerID"})
	data["WorkspaceID"] = ctx.Param("id")

	serverId, exists := data["ServerID"]
	if !exists {
		schemas.MakeErrorResponse(ctx, "ServerID is required", 400)
		return
	}
	var s models.Server
	results := di.Container.DB.First(&s, serverId)
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Server Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	if value, exists := data["Options"]; exists {
		data["Options"], _ = json.Marshal(value)
	}

	var obj models.Rule
	err := gconv.Struct(data, &obj)
	if err != nil {
		schemas.MakeErrorResponse(ctx, err, 400)
		return
	}

	results = di.Container.DB.Create(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, obj, nil)
}

func Update(ctx *gin.Context) {
	var obj models.Rule
	results := di.Container.DB.First(&obj, ctx.Param("ruleId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Rule Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Options", "Enable"})

	if value, exists := data["Options"]; exists {
		data["Options"], _ = json.Marshal(value)
	}
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[gstr.CaseSnake(k)] = v
	}

	results = di.Container.DB.Model(&models.Rule{}).Where("id = ?", ctx.Param("ruleId")).Updates(payload)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, true, nil)
}

func Delete(ctx *gin.Context) {
	var obj models.Rule
	results := di.Container.DB.First(&obj, ctx.Param("ruleId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Rule Not Found", 404)
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
