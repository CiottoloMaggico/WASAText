package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"os"
	"sort"
)

func RunMigrations(db *sqlx.DB) {
	files, err := os.ReadDir("./service/database/migrations")
	if err != nil {
		panic("Unable to read db migrations directory")
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		migration, err := os.ReadFile("./service/database/migrations/" + file.Name())
		if err != nil {
			panic(err)
		}

		db.MustExec(string(migration))
	}
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sqlx.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	db.MustExec("PRAGMA foreign_keys = ON;")
	RunMigrations(db)

	return &appdbimpl{
		db,
	}, nil
}
