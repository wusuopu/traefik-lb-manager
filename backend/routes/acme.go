package routes

import (
	"app/config"
	"app/controllers/acme"

	"github.com/gin-gonic/gin"
)

func InitAcme(e *gin.Engine) {
	e.GET(config.Config.Server.SSLChallengeBaseUrl + "/.well-known/acme-challenge/:token", acme.TokenVerify)
}