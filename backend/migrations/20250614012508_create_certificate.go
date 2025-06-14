package migrations

import (
	"app/utils"
	"context"
	"database/sql"
	"time"

	"github.com/pressly/goose/v3"
)

// https://gorm.io/docs/migration.html
// 20250614012508-CreateCertificate

func init() {
	goose.AddMigrationContext(up20250614012508, down20250614012508)
}
func createModel20250614012508 () interface{} {
	type Certificate struct {
		BaseModel
		Name				string						`gorm:"type:varchar(100);"`
		Domain			string						`gorm:"type:varchar(100);"`
		Cert				string						`gorm:"type:text;"`
		Key					string						`gorm:"type:text;"`
		ExpiredAt 	time.Time
		Status			string						`gorm:"type:varchar(20);default:init;"`
		Enable			bool
		WorkspaceID int
		AcmeToken		string						`gorm:"type:varchar(100);"`
		AcmeKeyAuth	string						`gorm:"type:varchar(100);"`
	}
	return &Certificate{}
}
func up20250614012508(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250614012508()
		createTable(migrator, model)
	})
}

func down20250614012508(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250614012508()
		dropTable(migrator, model)
	})
}
