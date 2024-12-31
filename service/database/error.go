package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

var NoResult = errors.New("no rows in result set")
var ForeignKeyConstraint = errors.New("foreign key constraint violation")
var UniqueConstraint = errors.New("unique constraint violation")
var UnexpectedError = errors.New("unexpected error")

func HandleSqlError(err error) error {
	fmt.Println(err)
	var sqlErr sqlite3.Error
	if !errors.As(err, &sqlErr) {
		if errors.Is(err, sql.ErrNoRows) {
			return NoResult
		}
		return UnexpectedError
	}

	if errors.Is(sqlErr.Code, sqlite3.ErrConstraint) {
		return ForeignKeyConstraint
	} else if errors.Is(sqlErr.Code, sqlite3.ErrConstraintUnique) {
		return UniqueConstraint
	}

	return UnexpectedError
}
