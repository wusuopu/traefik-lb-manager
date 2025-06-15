package acme

import (
	"app/di"
	"app/models"

	"github.com/gin-gonic/gin"
)

func TokenVerify (ctx *gin.Context) {
	var cert models.Certificate
	di.Container.DB.
		Where("acme_token = ?", ctx.Param("token")).
		Order("updated_at desc").
		First(&cert)
	ctx.Data(200, "text/plain", []byte(cert.AcmeKeyAuth))
}