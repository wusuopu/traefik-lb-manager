package models

import "gorm.io/datatypes"


type Middleware struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	Category		string						`gorm:"type:varchar(40);"`
	Options			datatypes.JSON		`gorm:"type:json;"`		// 保存该中间件的配置信息
	WorkspaceID int
}