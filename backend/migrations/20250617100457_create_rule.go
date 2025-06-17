package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250617100457-CreateRule

func init() {
	goose.AddMigrationContext(up20250617100457, down20250617100457)
}
func createModel20250617100457 () interface{} {
	type Rule struct {
		BaseModel
		Options			datatypes.JSON		`gorm:"type:json;"`		// 保存配置信息
		Enable			bool
		WorkspaceID int
		ServerID		int
	}
	return &Rule{}
}
func up20250617100457(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250617100457()
		createTable(migrator, model)
	})
}

func down20250617100457(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250617100457()
		dropTable(migrator, model)
	})
}
