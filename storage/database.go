package storage

import (
	"database/sql"

	"github.com/wibl/webapp/model"
)

type dbStorage struct {
	*sql.DB
}

// NewDbStorage creates a database connection
func NewDbStorage(driverName, dataSource string) (Storage, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &dbStorage{db}, nil
}

//GetGroups gets groups from database
func (db *dbStorage) GetGroups() ([]*model.Group, error) {
	//TODO: implement
	return nil, nil
}

//CreateGroup creates group in database
func (db *dbStorage) CreateGroup(group *model.Group) error {
	res, err := db.Exec("INSERT INTO group VALUES(?)", group.Title)
	if err != nil {
		return err
	}
	group.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (db *dbStorage) CreateTemplate(template *model.Template) error {
	//TODO: implement
	return nil
}

func (db *dbStorage) GetTemplates(*model.Group) ([]*model.Template, error) {
	//TODO: implement
	return nil, nil
}

func (db *dbStorage) DeleteGroup(group *model.Group) error {
	//TODO: implement
	return nil
}
func (db *dbStorage) DeleteTemplate(template *model.Template) error {
	//TODO: implement
	return nil
}

func (db *dbStorage) UpdateGroup(group *model.Group) error {
	//TODO: implement
	return nil
}
func (db *dbStorage) UpdateTemplate(template *model.Template) error {
	//TODO: implement
	return nil
}
