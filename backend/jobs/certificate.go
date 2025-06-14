package jobs

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
)


var certQueue *gqueue.Queue = nil

func startCertificateInterval(t *time.Ticker) {
	for {
		select {
		case now := <-t.C:
			fmt.Printf("[%s] start check certificate\n", now.String())
		}
	}
  
}

func startCertificateQueueCheck() {
	for {
		select {
		case item := <-certQueue.C:
			fmt.Printf("[%s] start renew certificate: %d %d\n", time.Now().String(), item, certQueue.Len())
			if item == 2 {
				continue
			}

			fmt.Println("wait 3 seconds....")
			// 控制 letsencrypt api 请求频率
			time.Sleep(time.Second * 3)
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