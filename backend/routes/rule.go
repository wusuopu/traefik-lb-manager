package routes

import (
	"app/controllers/rule"
	"app/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRule(r *gin.RouterGroup, e *gin.Engine) {
	const baseUrl = "/:id/rules"

	r.GET(baseUrl + "/", middlewares.CheckWorkspace, rule.Index)
	r.POST(baseUrl + "/", middlewares.CheckWorkspace, rule.Create)
	r.PUT(baseUrl + "/:ruleId", middlewares.CheckWorkspace, rule.Update)
	r.DELETE(baseUrl + "/:ruleId", middlewares.CheckWorkspace, rule.Delete)
}
