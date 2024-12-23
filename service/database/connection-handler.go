package database

import "database/sql"

func SaveClose(rows *sql.Rows) error {
	if err := rows.Close(); err != nil {
		return err
	}
	return nil
}
