package models

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

type Server struct {
	BaseModel
	Name				string						`gorm:"type:varchar(100);"`
	Host				datatypes.JSON		`gorm:"type:json;"`		// 可以设置多个域名，用json 数组存储
	Enable			bool
	WorkspaceID int
}

func (s *Server) GetHostRules () string {
	hosts, err := gjson.LoadJson(s.Host)
	if err != nil {
		return ""
	}
	
	return "(" + strings.Join(
		lo.Map(hosts.Array(), func(x interface{}, _ int) string {
			return fmt.Sprintf("Host: `%s`", x.(string))
		}),
		" || ",
	) + ")"
}