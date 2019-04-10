package storage

import (
	"database/sql"

	"github.com/wibl/webapp/model"
)

// DataStorage is interface for working with data storage
type DataStorage interface {
	CreateGroup(group model.Group) (model.Group, error)
}

// DB is implementation of DataStorage interface for connecting to database
type DB struct {
	*sql.DB
}

// NewDb creates a database connection
func NewDb(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

//CreateGroup creates group in database
func (db *DB) CreateGroup(group model.Group) (model.Group, error) {
	//TODO: implement
	return group, nil
}
