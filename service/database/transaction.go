package database

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type AppTransaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	QueryStructRow(dest interface{}, query string, args ...interface{}) error
	QueryStruct(dest interface{}, query string, args ...interface{}) error
	Commit() error
	Rollback() error
}
type apptransactionimpl struct {
	*sqlx.Tx
}

func (db *apptransactionimpl) QueryStructRow(dest interface{}, query string, args ...interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		panic("dest must be a pointer")
	}
	destType = destType.Elem()
	if destType.Kind() != reflect.Struct {
		panic("dest must be a pointer to a struct")
	}

	if err := db.QueryRowx(query, args...).StructScan(dest); err != nil {
		return HandleDBError(err)
	}

	return nil
}

func (db *apptransactionimpl) QueryStruct(dest interface{}, query string, args ...interface{}) error {
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
		return HandleDBError(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	destType = destType.Elem()
	destValue := reflect.ValueOf(dest).Elem()
	for rows.Next() {
		newRow := reflect.New(destType)

		if err := rows.StructScan(newRow.Interface()); err != nil {
			return HandleDBError(err)
		}

		destValue.Set(reflect.Append(destValue, newRow.Elem()))
	}

	if err := rows.Err(); err != nil {
		return HandleDBError(err)
	}

	return nil
}
