package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250621005343-AddDescriptionToMiddleware

func init() {
	goose.AddMigrationContext(up20250621005343, down20250621005343)
}
func createModel20250621005343 () interface{} {
	type Middleware struct {
		BaseModel
		Name				string						`gorm:"type:varchar(100);"`
		Description	string						`gorm:"type:text;"`
		Category		string						`gorm:"type:varchar(40);"`
		Options			datatypes.JSON		`gorm:"type:json;"`		// 保存该中间件的配置信息
		WorkspaceID int
	}
	return &Middleware{}
}
func up20250621005343(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250621005343()
		addTableColumn(migrator, model, "Description")
	})
}

func down20250621005343(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250621005343()
		dropTableColumn(migrator, model, "Description")
	})
}
