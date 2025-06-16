package models

import "gorm.io/datatypes"

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