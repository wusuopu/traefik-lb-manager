package routes

import (
	"app/config"
	"app/controllers/acme"

	"github.com/gin-gonic/gin"
)

func InitAcme(e *gin.Engine) {
	e.GET(config.Config.SSL.ChallengeBaseUrl + "/.well-known/acme-challenge/:token", acme.TokenVerify)
}