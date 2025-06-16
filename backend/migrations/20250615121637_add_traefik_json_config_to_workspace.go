package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250615121637-AddTraefikJsonConfigToWorkspace

func init() {
	goose.AddMigrationContext(up20250615121637, down20250615121637)
}
func createModel20250615121637 () interface{} {
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
	return &Workspace{}
}
func up20250615121637(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250615121637()
		addTableColumn(migrator, model, "TraefikJsonConfig")
	})
}

func down20250615121637(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250615121637()
		dropTableColumn(migrator, model, "TraefikJsonConfig")
	})
}
