package certificate

import (
	"app/di"
	"app/jobs"
	"app/models"
	"app/schemas"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func Index(ctx *gin.Context) {
	// don't return password field
	var data []models.Certificate

	di.Container.DB.Select("id", "name", "domain", "status", "enable", "workspace_id", "expired_at", "created_at", "updated_at").Find(&data)

	schemas.MakeResponse(ctx, data, nil)
}

func Create(ctx *gin.Context) {
	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "Domain", "Enable"})
	data["WorkspaceID"] = ctx.Param("id")
	data["Domain"] = gstr.Trim(data["Domain"].(string))
	data["Status"] = models.CERTIFICATE_STATUS_INIT

	var obj models.Certificate
	err := gconv.Struct(data, &obj)
	if err != nil {
		schemas.MakeErrorResponse(ctx, err, 400)
		return
	}

	// 检查该域名是否存在
	var record models.Certificate
	di.Container.DB.Where("workspace_id = ?", data["WorkspaceID"]).Where("domain = ?", data["Domain"]).First(&record)
	if record.ID > 0 {
		schemas.MakeErrorResponse(ctx, "Domain already exists", 400)
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
	var obj models.Certificate
	results := di.Container.DB.First(&obj, ctx.Param("certificateId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Certificate Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	body, _ := ctx.Get("rawBody")
	data := gconv.Map(body)
	data = lo.PickByKeys(data, []string{"Name", "Enable"})
	payload := make(map[string]interface{})
	for k, v := range data {
		payload[gstr.CaseSnake(k)] = v
	}

	results = di.Container.DB.Model(&models.Certificate{}).Where("id = ?", ctx.Param("certificateId")).Updates(payload)
	if results.Error != nil {
		schemas.MakeErrorResponse(ctx, results.Error, 500)
		return
	}

	schemas.MakeResponse(ctx, true, nil)
}

func Delete(ctx *gin.Context) {
	var obj models.Certificate
	results := di.Container.DB.First(&obj, ctx.Param("certificateId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Certificate Not Found", 404)
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

func Renew(ctx *gin.Context) {
	var obj models.Certificate
	results := di.Container.DB.First(&obj, ctx.Param("certificateId"))
	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			schemas.MakeErrorResponse(ctx, "Certificate Not Found", 404)
		} else {
			schemas.MakeErrorResponse(ctx, results.Error, 500)
		}
		return
	}

	jobs.PushCertificateJob(obj.ID)

	schemas.MakeResponse(ctx, true, nil)
}