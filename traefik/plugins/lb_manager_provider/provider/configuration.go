package provider

import "fmt"


type ConfigurationMarshaler struct {
	BaseEndpoint   string		// http://<host>/workspaces/:id/
	QueryParams    string

	lastConfig		string
	version				string
}

func (c *ConfigurationMarshaler) MarshalJSON() ([]byte, error) {
	if c.lastConfig == "" {
		return nil, nil
	}

  return []byte(c.lastConfig), nil
}

func (c *ConfigurationMarshaler) LoadConfiguration() {
	// 检查是否需要更新配置
	needUpdate, err := c.checkVersion()
	if err != nil {
		Logger().Printf("[ERROR] checkVersion error %s", err)
		return
	}
	if !needUpdate {
		return
	}

	// 下载证书
	if err = sslManager.downloadCerts(fmt.Sprintf("%s/certs.json?%s", c.BaseEndpoint, c.QueryParams)); err != nil {
		Logger().Printf("[ERROR] download ssl certs error %s", err)
		return
	}

	// 下载配置
	if err = c.downloadConfig(); err != nil {
		Logger().Printf("[ERROR] download config error %s", err)
		return
	}
}

func (c *ConfigurationMarshaler) checkVersion() (bool, error) {
	ret, err := request(fmt.Sprintf("%s/version?%s", c.BaseEndpoint, c.QueryParams), "GET")
	if err != nil {
		return false, err
	}

	if string(ret) != c.version {
		c.version = string(ret)
		return true, nil
	}

  return false, nil
}

func (c *ConfigurationMarshaler) downloadConfig() error {
	ret, err := request(fmt.Sprintf("%s/traefik.json?%s", c.BaseEndpoint, c.QueryParams), "GET")
	if err != nil {
		return err
	}
	c.lastConfig = string(ret)
  return nil
}