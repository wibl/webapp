package storage

import (
	"database/sql"
	// For memory
	_ "github.com/go-sql-driver/mysql"
)
// DataStorage implements methods of working with data storage
type DataStorage interface {
	CreateGroup(title string) ()
}

// DB contains a specific implementation
type DB struct {
	*sql.DB
}

// NewDb creates a database connection
func NewDb (dataSource string) (*DB, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}