package models

import "gorm.io/datatypes"

const (
	WORKSPACE_CATEGORY_RANCHER = "rancher_v1"
	WORKSPACE_CATEGORY_PORTAINER = "portainer_swarm"
	WORKSPACE_CATEGORY_COMMON = "common"
	WORKSPACE_CATEGORY_CUSTOM = "custom"
)


type Workspace struct {
	BaseModel
	Name						string						`gorm:"type:varchar(100);"`
	Description			string						`gorm:"type:varchar(200);"`
	ManagerBaseUrl	string						`gorm:"type:varchar(500);"`
	Category				string						`gorm:"type:varchar(40);"`
	ApiBaseUrl			string						`gorm:"type:varchar(500);"`
	ApiKey					string						`gorm:"type:varchar(200);"`
	ApiSecret				string						`gorm:"type:varchar(200);"`
	Entrypoints			datatypes.JSON		`gorm:"type:json;"`		// 该实例可用的 entrypoints
	TraefikConfig		string						`gorm:"type:text;"`		// yaml 配置
	TraefikJsonConfig		string				`gorm:"type:text;"`		// json 配置
}