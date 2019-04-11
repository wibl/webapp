package storage

import (
	"github.com/wibl/webapp/model"
)

// Storage is interface for working with data storage
type Storage interface {
	GetGroups() ([]*model.Group, error)
	GetTemplates(*model.Group) ([]*model.Template, error)

	CreateGroup(group *model.Group) error
	CreateTemplate(template *model.Template) error

	DeleteGroup(group *model.Group) error
	DeleteTemplate(template *model.Template) error

	UpdateGroup(group *model.Group) error
	UpdateTemplate(template *model.Template) error
}
