package service

import (
	"app/di"
	"app/models"
	"app/schemas"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"gorm.io/gorm"
)


func Index(ctx *gin.Context) {
	var data []models.Service

	di.Container.DB.
		Where("workspace_id = ?", ctx.Param("id")).
		Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}

// 获取 rancher / portainer 的服务
func ExternalIndex(ctx *gin.Context) {
	var obj models.Workspace
	results := di.Container.DB.First(&obj, ctx.Param("id"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	if obj.Category == models.WORKSPACE_CATEGORY_RANCHER {
		data, err := di.Service.ServiceService.FetchRancherServices(&obj)
		if err != nil {
			schemas.MakeErrorResponse(ctx, err, 500)
			return
		}
		schemas.MakeResponse(ctx, data, nil)
		return
	} else if obj.Category == models.WORKSPACE_CATEGORY_PORTAINER {
		data, err := di.Service.ServiceService.FetchPortainerServices(&obj)
		if err != nil {
			schemas.MakeErrorResponse(ctx, err, 500)
			return
		}
		schemas.MakeResponse(ctx, data, nil)
		return
	}

	schemas.MakeErrorResponse(ctx, fmt.Errorf("workspace category is not support"), 400)
}

func Create(ctx *gin.Context) {
	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "LBServers"})
	data["WorkspaceID"] = ctx.Param("id")

	if value, exists := data["LBServers"]; exists {
		data["LBServers"], _ = json.Marshal(value)
	}

	var obj models.Service
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
	var obj models.Service
	results := di.Container.DB.First(&obj, ctx.Param("serviceId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Service Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "LBServers"})

	if value, exists := data["LBServers"]; exists {
		data["LBServers"], _ = json.Marshal(value)
	}
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[gstr.CaseSnake(k)] = v
	}

	results = di.Container.DB.Model(&models.Service{}).Where("id = ?", ctx.Param("serviceId")).Updates(payload)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, true, nil)
}

func Delete(ctx *gin.Context) {
	var obj models.Service
	results := di.Container.DB.First(&obj, ctx.Param("serviceId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Service Not Found", 404)
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
