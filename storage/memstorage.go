package storage

import (
	"fmt"
	"sync"

	"github.com/wibl/webapp/model"
)

type memStorage struct {
	lastGroupID    int64
	groups         []*model.Group
	lastTemplateID int64
	templates      []*model.Template
	sync.Mutex
}

// NewMemStorage creates a MemStorage instance
func NewMemStorage() (Storage, error) {
	stor := &memStorage{groups: make([]*model.Group, 0)}
	return stor, nil
}

func (s *memStorage) CreateGroup(group *model.Group) error {
	s.Lock()
	s.lastGroupID++
	group.ID = s.lastGroupID
	s.groups = append(s.groups, group)
	s.Unlock()
	return nil
}

func (s *memStorage) GetAllGroups() ([]*model.Group, error) {
	return s.groups, nil
}

func (s *memStorage) GetGroup(id int64) (*model.Group, error) {
	for _, gr := range s.groups {
		if gr.ID == id {
			return gr, nil
		}
	}
	return nil, fmt.Errorf("%s %d", "Group not found with ID ", id)
}

func (s *memStorage) CreateTemplate(template *model.Template) error {
	s.Lock()
	s.lastTemplateID++
	template.ID = s.lastTemplateID
	s.templates = append(s.templates, template)
	s.Unlock()
	return nil
}

func (s *memStorage) GetAllTemplates(group *model.Group) ([]*model.Template, error) {
	templates := make([]*model.Template, 0)
	for _, template := range s.templates {
		if template.GroupID == group.ID {
			templates = append(templates, template)
		}
	}
	return templates, nil
}

func (s *memStorage) GetTemplate(id int64) (*model.Template, error) {
	for _, template := range s.templates {
		if template.ID == id {
			return template, nil
		}
	}
	return nil, fmt.Errorf("%s %d", "Template not found with ID ", id)
}

func (s *memStorage) DeleteGroup(deletedGroup *model.Group) error {
	s.Lock()

	templates := make([]*model.Template, 0)
	for _, template := range s.templates {
		if template.GroupID != deletedGroup.ID {
			templates = append(templates, template)
		}
	}

	groups := make([]*model.Group, 0)
	for _, group := range s.groups {
		if group.ID != deletedGroup.ID {
			groups = append(groups, group)
		}
	}

	s.groups = groups
	s.templates = templates

	s.Unlock()
	return nil
}

func (s *memStorage) DeleteTemplate(deletedTemplate *model.Template) error {
	s.Lock()

	templates := make([]*model.Template, 0)
	for _, template := range s.templates {
		if template.ID != deletedTemplate.ID {
			templates = append(templates, template)
		}
	}

	s.templates = templates

	s.Unlock()
	return nil
}

// func (s *memStorage) UpdateGroup(updatedGroup *model.Group) error {
// 	s.Lock()

// 	for idx, group := range s.groups {
// 		if group.ID == updatedGroup.ID {
// 			s.groups[idx] = updatedGroup
// 			break
// 		}
// 	}

// 	s.Unlock()
// 	return nil
// }

// func (s *memStorage) UpdateTemplate(updatedTemplate *model.Template) error {
// 	s.Lock()

// 	for idx, template := range s.templates {
// 		if template.ID == updatedTemplate.ID {
// 			s.templates[idx] = updatedTemplate
// 		}
// 	}

// 	s.Unlock()
// 	return nil
// }
