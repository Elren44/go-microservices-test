package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddUsersTable, downAddUsersTable)
}

func upAddUsersTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	_, err := tx.Exec(`
			CREATE TABLE users (
				id BIGSERIAL PRIMARY KEY,
				email TEXT UNIQUE NOT NULL,
				password_hash TEXT NOT NULL
			);
		`)
	return err
}

func downAddUsersTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.

	_, err := tx.Exec(`
			DROP TABLE users;
		`)
	return err
}
