package database

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"strings"
)

var ErrNoResult = errors.New("no rows in result set")
var ErrForeignKey = errors.New("foreign key constraint violation")
var ErrUnique = errors.New("unique constraint violation")
var ErrCheck = errors.New("check constraint violation")
var ErrTrigger = errors.New("trigger constraint violation")
var ErrUnexpected = errors.New("unexpected error")

type DBError struct {
	ErrType error
	Detail  string
}

func (e DBError) Error() string {
	return e.Detail
}

func NewDBError(errType error, msg string) error {
	return DBError{ErrType: errType, Detail: msg}
}

func HandleDBError(err error) error {
	var sqlErr sqlite3.Error
	if errors.As(err, &sqlErr) {
		errExtended := sqlErr.ExtendedCode
		if errors.Is(errExtended, sqlite3.ErrConstraintForeignKey) {
			return NewDBError(ErrForeignKey, err.Error())
		} else if errors.Is(errExtended, sqlite3.ErrConstraintUnique) {
			return NewDBError(ErrUnique, err.Error())
		} else if errors.Is(errExtended, sqlite3.ErrConstraintTrigger) {
			errMsg := strings.TrimPrefix(err.Error(), "TRIGGER: ")
			return NewDBError(ErrTrigger, errMsg)
		} else if errors.Is(errExtended, sqlite3.ErrConstraintCheck) {
			return NewDBError(ErrCheck, err.Error())
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return NewDBError(ErrNoResult, "no rows in the result set")
	}
	return NewDBError(ErrUnexpected, "unexpected error")
}
