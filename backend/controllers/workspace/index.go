package workspace

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
	// don't return password field
	var data []models.Workspace

	di.Container.DB.
		Select("id", "name", "description", "manager_base_url", "category", "api_base_url", "api_key", "api_secret", "Entrypoints", "created_at", "updated_at").
		Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}

func Create(ctx *gin.Context) {
	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "Description", "ManagerBaseUrl", "Category", "ApiBaseUrl", "ApiKey", "ApiSecret", "Entrypoints"})

	if value, exists := data["Entrypoints"]; exists {
		data["Entrypoints"], _ = json.Marshal(value)
	}

	var obj models.Workspace
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

func Show(ctx *gin.Context) {
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

	schemas.MakeResponse(ctx, obj, nil)
}

func Update(ctx *gin.Context) {
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

	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "Description", "ManagerBaseUrl", "Category", "ApiBaseUrl", "ApiKey", "ApiSecret", "Entrypoints"})
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[gstr.CaseSnake(k)] = v
	}
	if value, exists := payload["entrypoints"]; exists {
		payload["entrypoints"], _ = json.Marshal(value)
	}

	results = di.Container.DB.Model(&models.Workspace{}).Where("id = ?", ctx.Param("id")).Updates(payload)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, true, nil)
}

func Delete(ctx *gin.Context) {
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
	if ctx.Query("name") != obj.Name {
		schemas.MakeErrorResponse(ctx, "Forbidden", 403)
		return
	}

	di.Container.DB.Delete(&obj)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	// TODO 删除相关的其他资源

	schemas.MakeResponse(ctx, true, nil)
}

func GenerateTraefikConfig(ctx *gin.Context) {
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
	if obj.Category == models.WORKSPACE_CATEGORY_CUSTOM {
		schemas.MakeErrorResponse(ctx, "Not Found", 400)
		return
	}
	// TODO 生成配置
	schemas.MakeResponse(ctx, true, nil)
}

func UpdateTraefikConfig(ctx *gin.Context) {
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

	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	payload := lo.PickByKeys(data, []string{"traefik_config"})

	results = di.Container.DB.Model(&models.Workspace{}).Where("id = ?", ctx.Param("id")).Updates(payload)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, true, nil)
}

func ShowTraefikConfig(ctx *gin.Context) {
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
	if ctx.Query("name") != obj.Name {
		schemas.MakeErrorResponse(ctx, "Forbidden", 403)
		return
	}

	ctx.Data(200, "text/yaml", []byte(obj.TraefikConfig))
}

func ShowConfigVersion(ctx *gin.Context) {
	// TODO 计算当前配置的版本号，以 workspace 或者 certificate 的更新时间作为版本号
	schemas.MakeResponse(ctx, 100, nil)
}