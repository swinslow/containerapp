// Package models defines the database and model framework.
package models

import (
	"database/sql"
	"time"

	// postgres driver
	_ "github.com/lib/pq"
)

// Datastore defines the interface to be implemented by models,
// using either a backing database (production) or mocks (test).
type Datastore interface {
	// Users
	GetAllUsers() ([]*User, error)
	GetUserByID(id uint32) (*User, error)
	GetUserByEmail(email string) (*User, error)
	AddUser(uint32, string, string, bool) error
	// VisitedPaths
	GetAllVisitedPaths() ([]*VisitedPath, error)
	GetAllVisitedPathsForUserID(uint32) ([]*VisitedPath, error)
	AddVisitedPath(string, time.Time, uint32) error
}

// DB holds the actual database/sql object as well as its related
// database statements.
type DB struct {
	sqldb *sql.DB
}

// NewDB opens and returns an initialized DB object.
func NewDB(srcName string) (*DB, error) {
	sqldb, err := sql.Open("postgres", srcName)
	if err != nil {
		return nil, err
	}
	if err = sqldb.Ping(); err != nil {
		return nil, err
	}

	db := &DB{sqldb: sqldb}
	return db, nil
}

// InitDBTables confirms that the DB tables are good to go.
func (db *DB) InitDBTables() error {
	var err error
	err = db.CreateTableVisitedPath()
	if err != nil {
		return err
	}

	err = db.CreateTableUsers()
	if err != nil {
		return err
	}

	return nil
}

// CloseDB closes the DB object when the program is exiting.
func (db *DB) CloseDB() {
	if db != nil && db.sqldb != nil {
		db.sqldb.Close()
	}
}
