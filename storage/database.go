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
	db.Exec("CREATE TABLE IF NOT EXISTS `group` (id INTEGER PRIMARY KEY AUTOINCREMENT, title VARCHAR(255));")
	db.Exec("CREATE TABLE IF NOT EXISTS `template` (id INTEGER PRIMARY KEY AUTOINCREMENT, groupid INTEGER, title VARCHAR(255), queue VARCHAR(255), body VARCHAR(255));")
	return &dbStorage{db}, nil
}

//CreateGroup creates group in database
func (db *dbStorage) CreateGroup(group *model.Group) error {
	stmt, err := db.Prepare("INSERT INTO `group`(title) VALUES(?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(group.Title)
	if err != nil {
		return err
	}
	group.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

//GetAllGroups gets groups from database
func (db *dbStorage) GetAllGroups() ([]*model.Group, error) {
	rows, err := db.Query("SELECT * FROM `group`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []*model.Group{}
	for rows.Next() {
		var group model.Group
		if err = rows.Scan(&group.ID, &group.Title); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}

// GetGroup returns the group by id
func (db *dbStorage) GetGroup(id int64) (*model.Group, error) {
	group := model.Group{}
	err := db.QueryRow("SELECT * FROM `group` WHERE id = $1", id).Scan(&group.ID, &group.Title)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (db *dbStorage) UpdateGroup(group *model.Group) error {
	_, err := db.Exec("UPDATE `group` SET title = $1 WHERE id = $2", group.Title, group.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbStorage) DeleteGroup(group *model.Group) error {
	_, err := db.Exec("DELETE FROM `group` WHERE id = $1", group.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbStorage) CreateTemplate(template *model.Template) error {
	res, err := db.Exec("INSERT INTO `template`(groupid, title, queue, body) VALUES($1, $2, $3, $4)",
		template.GroupID, template.Title, template.Queue, template.Body)
	if err != nil {
		return err
	}
	template.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (db *dbStorage) GetAllTemplates(gr *model.Group) ([]*model.Template, error) {
	rows, err := db.Query("SELECT * FROM `template` WHERE groupid = $1", gr.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := []*model.Template{}
	for rows.Next() {
		var tmp model.Template
		err = rows.Scan(&tmp.ID, &tmp.GroupID, &tmp.Title, &tmp.Queue, &tmp.Body)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (db *dbStorage) GetTemplate(id int64) (*model.Template, error) {
	tmp := model.Template{}
	err := db.QueryRow("SELECT * FROM `template` WHERE id = $1", id).Scan(&tmp.ID, &tmp.GroupID, &tmp.Title, &tmp.Queue, &tmp.Body)
	if err != nil {
		return nil, err
	}
	return &tmp, nil
}

func (db *dbStorage) UpdateTemplate(template *model.Template) error {
	_, err := db.Exec(
		"UPDATE `template` SET title = $1, queue = $2, body = $3 WHERE id = $4",
		template.Title, template.Queue, template.Body, template.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbStorage) DeleteTemplate(template *model.Template) error {
	_, err := db.Exec("DELETE FROM `template` WHERE id = $1", template.ID)
	if err != nil {
		return err
	}
	return nil
}
