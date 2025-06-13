package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250613030724-CreateServer

func init() {
	goose.AddMigrationContext(up20250613030724, down20250613030724)
}
func createModel20250613030724 () interface{} {
	type Server struct {
		BaseModel
		Name				string						`gorm:"type:varchar(100);"`
		Host				datatypes.JSON		`gorm:"type:json;"`
		EnableSSL		bool
		Enable			bool
		WorkspaceID int
	}
	return &Server{}
}
func up20250613030724(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250613030724()
		createTable(migrator, model)
	})
}

func down20250613030724(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250613030724()
		dropTable(migrator, model)
	})
}
