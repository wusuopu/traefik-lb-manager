package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

const dir = "/etc/traefik/ssl"

type SSLCertResponseItem struct {
	ID 				uint
	Name			string
	Domain		string
	Cert			string
	Key				string
	ExpiredAt time.Time
	Enable		bool
}
type SSLCertResponse struct {
	Data			[]SSLCertResponseItem
}

type SSLManager struct {
}
var sslManager = SSLManager{}

func InitSSLDir() error {
  // 清空 /etc/traefik/ssl 目录
	os.RemoveAll(dir)
	return os.MkdirAll(dir, 0777)
}

func (s *SSLManager) downloadCerts(url string) error {
	// 下载证书文件
	body, err := request(url, "GET")
	if err != nil {
		return err
	}


	// 保存证书
	var resp SSLCertResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	os.MkdirAll(dir, 0777)
	for _, item := range resp.Data {
		if item.Cert == "" || item.Key == "" {
			continue
		}

		name := item.Name
		if name == "" {
			name = fmt.Sprintf("%s__%d", item.Domain, item.ID)
		}
		filename := path.Join(dir, name)
		err = os.WriteFile(filename + ".crt", []byte(item.Cert), 0666)
		if err != nil {
			Logger().Printf("[ERROR] write file %s.crt failed: %s", filename, err.Error())
		}
		err = os.WriteFile(filename + ".key", []byte(item.Key), 0666)
		if err != nil {
			Logger().Printf("[ERROR] write file %s.key failed: %s", filename, err.Error())
		}
	}

	return nil
}