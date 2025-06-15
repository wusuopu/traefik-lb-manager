package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250614061333-CreateSystemSetting

func init() {
	goose.AddMigrationContext(up20250614061333, down20250614061333)
}
func createModel20250614061333 () interface{} {
	type SystemSetting struct{
		BaseModel
		Name				string						`gorm:"type:varchar(100);"`
		Value				datatypes.JSON		`gorm:"type:json;"`
	}
	return &SystemSetting{}
}
func up20250614061333(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250614061333()
		createTable(migrator, model)
	})
}

func down20250614061333(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250614061333()
		dropTable(migrator, model)
	})
}
