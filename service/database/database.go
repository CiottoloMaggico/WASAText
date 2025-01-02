/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type AppDatabase interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryStructRow(dest interface{}, query string, args ...interface{}) error
	QueryStruct(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	StartTx() (AppTransaction, error)

	Ping() error
}

type appdbimpl struct {
	*sqlx.DB
}

func (db *appdbimpl) StartTx() (AppTransaction, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, DBError(err)
	}
	return &apptransactionimpl{tx}, nil
}

func (db *appdbimpl) QueryStructRow(dest interface{}, query string, args ...interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		panic("dest must be a pointer")
	}
	destType = destType.Elem()
	if destType.Kind() != reflect.Struct {
		panic("dest must be a pointer to a struct")
	}

	if err := db.QueryRowx(query, args...).StructScan(dest); err != nil {
		return DBError(err)
	}

	return nil
}

func (db *appdbimpl) QueryStruct(dest interface{}, query string, args ...interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		panic("dest must be a pointer to an array or slice")
	}
	destType = destType.Elem()
	if destType.Kind() != reflect.Slice && destType.Kind() != reflect.Array {
		panic("dest must be a pointer to a slice or array")
	}

	rows, err := db.Queryx(query, args...)
	if err != nil {
		return DBError(err)
	}
	defer func(rows *sqlx.Rows) error {
		if err := rows.Close(); err != nil {
			return DBError(err)
		}
		return nil
	}(rows)

	destType = destType.Elem()
	destValue := reflect.ValueOf(dest).Elem()
	for rows.Next() {
		newRow := reflect.New(destType)

		if err := rows.StructScan(newRow.Interface()); err != nil {
			return DBError(err)
		}

		destValue.Set(reflect.Append(destValue, newRow.Elem()))
	}

	if err := rows.Err(); err != nil {
		return DBError(err)
	}

	return nil
}
