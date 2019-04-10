package storage

import (
	"github.com/wibl/webapp/model"
)

// DataStorage is interface for working with data storage
type DataStorage interface {
	CreateGroup(group *model.Group) error
}
