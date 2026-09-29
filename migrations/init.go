package migrations

import (
	"context"
	"database/sql"
)

func Run(ctx context.Context, db *sql.DB) error {
	if err := createUsersTable(ctx, db); err != nil {
		return err
	}

	if err := createUserPhonesTable(ctx, db); err != nil {
		return err
	}

	return nil
}

const createUsersTableQuery = `
	CREATE TABLE IF NOT EXISTS users 
	(
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		login VARCHAR NOT NULL UNIQUE,
		password VARCHAR NOT NULL,
		username VARCHAR NOT NULL UNIQUE,
		age VARCHAR NOT NULL DEFAULT 0,
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	)
`

func createUsersTable(ctx context.Context, db *sql.DB) error {
	stmt, err := db.Prepare(createUsersTableQuery)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

const createUserPhonesTableQuery = `
	CREATE TABLE IF NOT EXISTS user_phones 
	(
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id VARCHAR NOT NULL,
		phone VARCHAR NOT NULL UNIQUE,
		description VARCHAR NOT NULL,
		is_fax BOOLEAN DEFAULT false,
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE INDEX phone_user_id ON user_phones (user_id);
`

func createUserPhonesTable(ctx context.Context, db *sql.DB) error {
	stmt, err := db.Prepare(createUserPhonesTableQuery)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}
