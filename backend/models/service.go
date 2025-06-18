package models

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"gorm.io/datatypes"
)

// https://doc.traefik.io/traefik/routing/services/#servers-load-balancer
type ServiceLoadBalancerServer struct {
  Url						string
	PreservePath	bool
	Weight				uint
}

type Service struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	LBServers		datatypes.JSON		`gorm:"type:json;"`		// 保存 loadBalancer.servers 数组数据
	WorkspaceID int
}

func (s *Service) GetLBRuleMap() map[string]interface{} {
	options, err := gjson.LoadJson(s.LBServers)
	if err != nil {
		return nil
	}

	lbServers := make([]map[string]interface{}, 0)
	for _, v := range options.Array() {
		if v == nil {
			continue
		}
		item := v.(map[string]interface{})
		url, ok := item["url"]
		if !ok {
			continue
		}

		obj := make(map[string]interface{})
		obj["url"] = url
		if v, ok := item["preservePath"]; ok && v.(bool) {
			obj["preservePath"] = v.(bool)
		}
		if v, ok := item["weight"]; ok && v.(float64) > 0 {
			obj["weight"] = gconv.Uint(v)
		}
		lbServers = append(lbServers, obj)

	}

	if len(lbServers) == 0 {
		return nil
	}

	ret := make(map[string]interface{})
	ret["servers"] = lbServers
	ret["passHostHeader"] = true

	return ret
}