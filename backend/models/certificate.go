package models

import (
	"time"
)

const (
	CERTIFICATE_STATUS_INIT = "init"
	CERTIFICATE_STATUS_PENDING = "pending"		// 开始调用 api 获取证书
	CERTIFICATE_STATUS_COMPLETE = "complete"
	CERTIFICATE_STATUS_FAILED = "failed"
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
	AcmeToken		string						`gorm:"type:varchar(100);"`
	AcmeKeyAuth	string						`gorm:"type:varchar(100);"`
}