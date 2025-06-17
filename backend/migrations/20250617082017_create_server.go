package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250617082017-CreateServer

func init() {
	goose.AddMigrationContext(up20250617082017, down20250617082017)
}
func createModel20250617082017 () interface{} {
	type Server struct {
		BaseModel
		Name				string						`gorm:"type:varchar(100);"`
		Host				datatypes.JSON		`gorm:"type:json;"`		// 可以设置多个域名，用json 数组存储
		Enable			bool
		WorkspaceID int
	}
	return &Server{}
}
func up20250617082017(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250617082017()
		createTable(migrator, model)
	})
}

func down20250617082017(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250617082017()
		dropTable(migrator, model)
	})
}
