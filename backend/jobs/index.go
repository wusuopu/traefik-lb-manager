package jobs

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/container/gqueue"
)

var ticker *time.Ticker = nil
func Start() {
	if ticker == nil {
		// 开发环境 5秒 检查一次
		duration := time.Second * 5
		if gin.Mode() == gin.ReleaseMode {
			// 生产环境 2小时 检查一次
			duration = time.Hour * 2
		}
		ticker = time.NewTicker(duration)
		go startCertificateInterval(ticker)
	}
	if certQueue == nil {
		certQueue = gqueue.New()
		go startCertificateQueueCheck()
	}
}

func Stop () {
	if ticker != nil {
		ticker.Stop()
		ticker = nil
	}
	if certQueue != nil {
		certQueue.Close()
		certQueue = nil
	}
}