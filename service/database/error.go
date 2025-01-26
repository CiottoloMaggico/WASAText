package database

import (
	"database/sql"
	"errors"
	"fmt"
	sqlite3 "github.com/mattn/go-sqlite3"
	"strings"
)

var ErrNoResult = errors.New("no rows in result set")
var ErrForeignKey = errors.New("foreign key constraint violation")
var ErrUnique = errors.New("unique constraint violation")
var ErrCheck = errors.New("check constraint violation")
var ErrTrigger = errors.New("trigger constraint violation")
var ErrUnexpected = errors.New("unexpected error")

type ErrDB struct {
	ErrType error
	Detail  string
}

func (e ErrDB) Error() string {
	return e.Detail
}

func NewErrDB(errType error, msg string) error {
	return ErrDB{ErrType: errType, Detail: msg}
}

func DBError(err error) error {
	var sqlErr sqlite3.Error
	if errors.As(err, &sqlErr) {
		errExtended := sqlErr.ExtendedCode
		if errors.Is(errExtended, sqlite3.ErrConstraintForeignKey) {
			return NewErrDB(ErrForeignKey, err.Error())
		} else if errors.Is(errExtended, sqlite3.ErrConstraintUnique) {
			fmt.Println("unireq")
			return NewErrDB(ErrUnique, err.Error())
		} else if errors.Is(errExtended, sqlite3.ErrConstraintTrigger) {
			errMsg := strings.TrimPrefix(err.Error(), "TRIGGER: ")
			return NewErrDB(ErrTrigger, errMsg)
		} else if errors.Is(errExtended, sqlite3.ErrConstraintCheck) {
			return NewErrDB(ErrCheck, err.Error())
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return NewErrDB(ErrNoResult, "no rows in the result set")
	}
	return NewErrDB(ErrUnexpected, "unexpected error")
}
