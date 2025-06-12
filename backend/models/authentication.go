package models


type Authentication struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	Username		string						`gorm:"type:varchar(100);"`
	Password		string						`gorm:"type:varchar(100);"`
	HashedPw		string						`gorm:"type:varchar(100);"`
}