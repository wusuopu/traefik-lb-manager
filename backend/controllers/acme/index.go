package acme

import "github.com/gin-gonic/gin"

func TokenVerify (ctx *gin.Context) {
	ctx.Data(200, "text/plain", []byte("ok"))
}