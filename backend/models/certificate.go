package models

import (
	"time"
)

const (
	CERTIFICATE_STATUS_INIT = "init"
	CERTIFICATE_STATUS_PENDING = "pending"
	CERTIFICATE_STATUS_COMPLETE = "complete"
)

type Certificate struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	Domain			string						`gorm:"type:varchar(100);"`
	Cert				string						`gorm:"type:text;"`
	Key					string						`gorm:"type:text;"`
	ExpiredAt 	time.Time
	Status			string						`gorm:"type:varchar(20);default:init;"`
	Enable			bool
	WorkspaceID int
}