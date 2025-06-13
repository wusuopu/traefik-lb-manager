package models

import "gorm.io/datatypes"

type Server struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	Host				datatypes.JSON		`gorm:"type:json;"`		// 可以设置多个域名，用json 数组存储
	EnableSSL		bool
	Enable			bool
	WorkspaceID int
}