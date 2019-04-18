package storage

import (
	"github.com/wibl/webapp/model"
)

// Storage is interface for working with data storage
type Storage interface {
	GetAllGroups() ([]*model.Group, error)
	GetAllTemplates(*model.Group) ([]*model.Template, error)

	GetGroup(id int64) (*model.Group, error)
	GetTemplate(id int64) (*model.Template, error)

	CreateGroup(group *model.Group) error
	CreateTemplate(template *model.Template) error

	DeleteGroup(group *model.Group) error
	DeleteTemplate(template *model.Template) error

	// UpdateGroup(group *model.Group) error
	// UpdateTemplate(template *model.Template) error
}
