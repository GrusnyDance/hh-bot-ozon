package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTable, downCreateTable)
}

func upCreateTable(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE IF NOT EXISTS vacancies (
    				id bigserial CONSTRAINT vacancies_pk PRIMARY KEY,
    				created_at TIMESTAMP,
    				name VARCHAR,
    				link VARCHAR
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downCreateTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE IF EXISTS vacancies;`)
	if err != nil {
		return err
	}
	return nil
}
