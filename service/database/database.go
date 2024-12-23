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
	"errors"
	"fmt"
)

// TODO: simplify database package using sqlx to populate structs and then translators in controllers

type AppDatabase interface {
	QueryRowGroup(query string, params ...any) (*Group, error)
	QueryRowImage(query string, params ...any) (*Image, error)
	QueryMessageWithAuthorAndAttachment(query string, params ...any) (MessageWithAuthorAndAttachmentList, error)
	QueryRowMessageWithAuthorAndAttachment(query string, params ...any) (*MessageWithAuthorAndAttachment, error)
	QueryMessageInfo(query string, params ...any) (MessageInfoList, error)
	QueryRowMessageInfo(query string, params ...any) (*MessageInfo, error)
	QueryRowUserConversation(query string, params ...any) (*UserConversation, error)
	QueryUserConversation(query string, params ...any) ([]UserConversation, error)
	QueryRowUser(query string, params ...any) (*User, error)
	QueryUserWithImage(query string, params ...any) (UserWithImageList, error)
	QueryRowUserWithImage(query string, params ...any) (*UserWithImage, error)
	QueryRowMessage(query string, params ...any) (*Message, error)
	QueryRowChat(query string, params ...any) (*Chat, error)
	Exec(query string, params ...any) (sql.Result, error)

	Ping() error
}

type BaseDatabase interface {
	AppDatabase
	BeginTx() (TransactionDatabase, error)
	Ping() error
}

type TransactionDatabase interface {
	AppDatabase
	Commit() error
	Rollback() error
}

type databaseInterface interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type transactionInterface interface {
	databaseInterface
	Commit() error
	Rollback() error
}

type appdbimpl struct {
	c databaseInterface
}

type apptransactionimpl struct {
	appdbimpl
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}
	for key, value := range tables {
		err := tx.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?;`, key).Scan(&tableName)
		if errors.Is(err, sql.ErrNoRows) {
			completeQuery := fmt.Sprintf("%s\n%s", value, initializers[key])
			if _, err = tx.Exec(completeQuery); err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					return nil, fmt.Errorf("could not rollback transaction: %w", rollbackErr)
				}
				return nil, fmt.Errorf("error creating database structure, rolled back to initial state: %w", err)
			}
		}
	}

	for _, trigger := range triggers {
		if _, err := tx.Exec(trigger); err != nil {
			return nil, fmt.Errorf("could not create trigger: %w", err)
		}
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", commitErr)
	}

	return &appdbimpl{
		c: databaseInterface(db),
	}, nil
}

func (db *appdbimpl) Ping() error {
	appDb, ok := db.c.(*sql.DB)
	if !ok {
		return errors.New("database is not pingable")
	}
	return appDb.Ping()
}

func (db *appdbimpl) Exec(query string, params ...any) (sql.Result, error) {
	return db.c.Exec(query, params...)
}

func (db *appdbimpl) BeginTx() (TransactionDatabase, error) {
	transactionDb, ok := db.c.(*sql.DB)
	if !ok {
		return nil, errors.New("impossible to start transactions with this database")
	}

	tx, err := transactionDb.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transaction: %w", err)
	}

	return &apptransactionimpl{
		appdbimpl{tx},
	}, nil
}

func (db *apptransactionimpl) Commit() error {
	tx, ok := db.c.(*sql.Tx)
	if !ok {
		return errors.New("impossible to perform transaction commit with this database")
	}

	return tx.Commit()
}

func (db *apptransactionimpl) Rollback() error {
	tx, ok := db.c.(*sql.Tx)
	if !ok {
		return errors.New("impossible to perform transaction rollback with this database")
	}
	return tx.Rollback()
}
