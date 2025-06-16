package jobs

import (
	"app/di"
	"app/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/container/gqueue"
)


var certQueue *gqueue.Queue = nil

func startCertificateInterval(t *time.Ticker) {
	for {
		select {
		case now := <-t.C:
			di.Container.Logger.Info("start check certificate")

			var certificates []models.Certificate
			endTime := now.Add(time.Hour * 24 * 5)		// 续签还剩5天的证书
			results := di.Container.DB.
				Where("created_at < ?", now.Add(time.Hour * 2)).		// 只处理两小时前创建的数据
				Where("enable = ?", 1).
				Where("expired_at < ? OR expired_at is NULL", endTime).
				Limit(50).
				Find(&certificates)

			if results.Error != nil {
				continue
			}

			for _, v := range certificates {
				if certQueue.Len() < 600 {
					certQueue.Push(v.ID)
				}
			}
		}
	}
  
}

func startCertificateQueueCheck() {
	for {
		select {
		case item := <-certQueue.C:
			di.Container.Logger.Info(fmt.Sprintf("start renew certificate: %d %d", item, certQueue.Len()))
			var cert models.Certificate
			di.Container.DB.Where("id = ?", item).Find(&cert)
			if cert.ID == 0 {
				continue
			}
			if !cert.Enable {
				continue
			}
			if cert.Status == models.CERTIFICATE_STATUS_COMPLETE && cert.ExpiredAt.After(time.Now().Add(time.Hour * 24 * 5)) {
				// 证书还有5天有效，不需要续签
				continue
			}

			if gin.Mode() == gin.ReleaseMode {
				// 生产环境才需要签发证书
				di.Service.CertificateService.Obtain(&cert)
			} else {
				di.Container.Logger.Debug("mock obtain certcrypto")
			}

			// 控制 letsencrypt api 请求频率
			time.Sleep(time.Second * 10)
		}
	}
}

func PushCertificateJob(item interface{}) bool {
  if certQueue == nil {
		return false
	}
	certQueue.Push(item)
	return true
}