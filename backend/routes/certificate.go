package routes

import (
	"app/controllers/certificate"
	"app/middlewares"

	"github.com/gin-gonic/gin"
)

func InitCertificate(r *gin.RouterGroup, e *gin.Engine) {
	const baseUrl = "/:id/certificates"

	r.GET(baseUrl + "/", middlewares.CheckWorkspace, certificate.Index)
	r.POST(baseUrl + "/", middlewares.CheckWorkspace, certificate.Create)
	r.PUT(baseUrl + "/:certificateId", middlewares.CheckWorkspace, certificate.Update)
	r.DELETE(baseUrl + "/:certificateId", middlewares.CheckWorkspace, certificate.Delete)
	r.PUT(baseUrl + "/:certificateId/renew", middlewares.CheckWorkspace, certificate.Renew)
}
