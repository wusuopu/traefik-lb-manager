package models

import "gorm.io/datatypes"

type Rule struct {
	BaseModel
	Options			datatypes.JSON		`gorm:"type:json;"`		// 保存配置信息
	Enable			bool
	WorkspaceID int
	ServerID		int
}