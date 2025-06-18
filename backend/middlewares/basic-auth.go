package middlewares

import (
	"app/config"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/net/gclient"
)

var basicAuth gin.HandlerFunc

func BasicAuthMiddleware(c *gin.Context) {
	if config.Config.Server.Username != "" && config.Config.Server.Password != "" {
		if basicAuth == nil {
			basicAuth = gin.BasicAuth(gin.Accounts{
				config.Config.Server.Username: config.Config.Server.Password,
			})
		}
		basicAuth(c)
		return
	}

	if config.Config.Server.RancherV1AuthUrl != "" {
		// 使用 rancher v1 的 cookie 进行认，该程序需要与 rancher 服务在同一域名下
		token, err := c.Cookie("token")
		if token == "" || err != nil {
			c.AbortWithStatus(401)
			return
		}

		ctx := context.Background()
		client := gclient.New().Timeout(time.Second * 5).Cookie(map[string]string{
			"token": token,
		})
		resp, err := client.Get(ctx, fmt.Sprintf("%s/v2-beta/projects?all=true&limit=1&sort=name", config.Config.Server.RancherV1AuthUrl))
		statusCode := 500
		if resp != nil {
			statusCode = resp.StatusCode
		}
		resp.Close()

		if err != nil || statusCode > 300 {
			c.AbortWithStatus(401)
			return
		}
	}
	if config.Config.Server.PortainerAuthUrl != "" {
		// 使用 portainer 的 cookie 进行认，该程序需要与 portainer 服务在同一域名下
		token, err := c.Cookie("portainer_api_key")
		if token == "" || err != nil {
			c.AbortWithStatus(401)
			return
		}

		ctx := context.Background()
		client := gclient.New().Timeout(time.Second * 5).Cookie(map[string]string{
			"portainer_api_key": token,
		})
		resp, err := client.Get(ctx, fmt.Sprintf("%s/api/users/me", config.Config.Server.PortainerAuthUrl))
		statusCode := 500
		if resp != nil {
			statusCode = resp.StatusCode
		}
		resp.Close()

		if err != nil || statusCode > 300 {
			c.AbortWithStatus(401)
			return
		}
	}

	c.Next()
}