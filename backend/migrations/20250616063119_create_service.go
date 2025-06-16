package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/datatypes"
)

// https://gorm.io/docs/migration.html
// 20250616063119-CreateService

func init() {
	goose.AddMigrationContext(up20250616063119, down20250616063119)
}
func createModel20250616063119 () interface{} {
	type Service struct {
		BaseModel
		Name				string						`gorm:"type:varchar(100);"`
		LBServers		datatypes.JSON		`gorm:"type:json;"`		// 保存 loadBalancer.servers 数组数据
		WorkspaceID int
	}
	return &Service{}
}
func up20250616063119(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20250616063119()
		createTable(migrator, model)
	})
}

func down20250616063119(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20250616063119()
		dropTable(migrator, model)
	})
}
