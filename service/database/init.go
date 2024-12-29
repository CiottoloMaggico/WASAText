package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"os"
	"sort"
)

func RunMigrations(db *sqlx.DB) error {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		migration, err := os.ReadFile("migrations/" + file.Name())
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}

		if _, err := tx.Exec(string(migration)); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sqlx.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// TODO: add run migrations from cmd line parameters like django
	//if err := RunMigrations(db); err != nil {
	//	return nil, err
	//}

	return &appdbimpl{
		db,
	}, nil
}
