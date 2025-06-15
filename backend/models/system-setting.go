package models

import (
	"gorm.io/datatypes"
)


type SystemSetting struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	Value				datatypes.JSON		`gorm:"type:json;"`
}
