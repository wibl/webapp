package storage

import (
	"database/sql"
	"log"

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
	db.Exec("CREATE TABLE IF NOT EXISTS `group` (id INTEGER PRIMARY KEY AUTOINCREMENT, title VARCHAR(255));")
	return &dbStorage{db}, nil
}

//GetGroups gets groups from database
func (db *dbStorage) GetAllGroups() ([]*model.Group, error) {
	rows, err := db.Query("SELECT * FROM `group`")
	if err != nil {
		return nil, err
	}

	groups := []*model.Group{}
	for rows.Next() {
		var group model.Group
		err = rows.Scan(&group.ID, &group.Title)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}

func (db *dbStorage) GetGroup(id int64) (*model.Group, error) {
	//TODO: implement
	return nil, nil
}

//CreateGroup creates group in database
func (db *dbStorage) CreateGroup(group *model.Group) error {
	stmt, err := db.Prepare("INSERT INTO `group`(title) VALUES(?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(group.Title)
	if err != nil {
		log.Fatal(err)
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

func (db *dbStorage) GetAllTemplates(*model.Group) ([]*model.Template, error) {
	//TODO: implement
	return nil, nil
}

func (db *dbStorage) GetTemplate(id int64) (*model.Template, error) {
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

// func (db *dbStorage) UpdateGroup(group *model.Group) error {
// 	//TODO: implement
// 	return nil
// }
// func (db *dbStorage) UpdateTemplate(template *model.Template) error {
// 	//TODO: implement
// 	return nil
// }
