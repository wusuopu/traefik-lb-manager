package models

import (
	"github.com/gogf/gf/v2/util/gconv"
	"gorm.io/datatypes"
)


type Middleware struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	Description	string						`gorm:"type:text;"`
	Category		string						`gorm:"type:varchar(40);"`
	Options			datatypes.JSON		`gorm:"type:json;"`		// 保存该中间件的配置信息
	WorkspaceID int
}

func (m *Middleware) GetRuleMap() map[string]interface{} {
	if m.Options == nil {
		return nil
	}
	ret := gconv.Map(string(m.Options))
	return ret
}