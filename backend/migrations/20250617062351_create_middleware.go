package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250617062351-CreateMiddleware

func init() {
	goose.AddMigrationContext(up20250617062351, down20250617062351)
}
func createModel20250617062351 () interface{} {
	type Middleware struct {
		BaseModel
		Name				string						`gorm:"type:varchar(100);"`
		Category		string						`gorm:"type:varchar(40);"`
		Options			datatypes.JSON		`gorm:"type:json;"`		// 保存该中间件的配置信息
		WorkspaceID int
	}
	return &Middleware{}
}
func up20250617062351(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250617062351()
		createTable(migrator, model)
	})
}

func down20250617062351(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250617062351()
		dropTable(migrator, model)
	})
}
