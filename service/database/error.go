package database

import (
	"database/sql"
	"errors"
	sqlite3 "github.com/mattn/go-sqlite3"
)

var NoResult = errors.New("no rows in result set")
var ForeignKeyConstraint = errors.New("foreign key constraint violation")
var UniqueConstraint = errors.New("unique constraint violation")
var CheckConstraint = errors.New("check constraint violation")
var TriggerConstraint = errors.New("trigger constraint violation")
var UnexpectedError = errors.New("unexpected error")

func DBError(err error) error {
	var sqlErr sqlite3.Error
	if errors.As(err, &sqlErr) {
		sqlExtendedCode := sqlErr.ExtendedCode
		switch sqlExtendedCode {
		case sqlite3.ErrConstraintForeignKey:
			return ForeignKeyConstraint
		case sqlite3.ErrConstraintUnique:
			return UniqueConstraint
		case sqlite3.ErrConstraintCheck:
			return CheckConstraint
		case sqlite3.ErrConstraintTrigger:
			return TriggerConstraint
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return NoResult
	}
	return UnexpectedError
}
