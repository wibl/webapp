package api

import (
	"fmt"
	"net/http"

	"github.com/wibl/webapp/model"
	"github.com/wibl/webapp/storage"
)

// TemplateService contains methods for working with templates
type TemplateService struct {
	Stor storage.Storage
}

// TemplateArgs contains possible arguments
type TemplateArgs struct {
	ID      int64
	GroupID int64
	Title   string
	Queue   string
	Body    string
}

// TemplateReply contains info for replay
type TemplateReply struct {
	Message   string
	Templates []*model.Template
}

// UpdateTemplateReply contains info for replay
type UpdateTemplateReply struct {
	Message  string
	Template *model.Template
}

// CreateTemplate allows to create a template
func (gs *TemplateService) CreateTemplate(rq *http.Request, args *TemplateArgs, reply *UpdateTemplateReply) error {
	newTmp := &model.Template{GroupID: args.GroupID, Title: args.Title, Queue: args.Queue, Body: args.Body}
	err := gs.Stor.CreateTemplate(newTmp)
	if err != nil {
		return err
	}
	reply.Message = "Created Template " + args.Title
	reply.Template = newTmp
	return nil
}

// GetAllTemplates return all templates
func (gs *TemplateService) GetAllTemplates(rq *http.Request, args *TemplateArgs, reply *TemplateReply) error {
	group, err := gs.Stor.GetGroup(args.GroupID)
	if err != nil {
		return err
	}
	tmps, err := gs.Stor.GetAllTemplates(group)
	if err != nil {
		return err
	}
	reply.Message = "Templates successfully received"
	reply.Templates = tmps
	return nil
}

// UpdateTemplate updates the template title and queue
func (gs *TemplateService) UpdateTemplate(rq *http.Request, args *TemplateArgs, reply *UpdateTemplateReply) error {
	tmp, err := gs.Stor.GetTemplate(args.ID)
	if err != nil {
		return err
	}
	tmp.Title = args.Title
	tmp.Queue = args.Queue
	tmp.Body = args.Body
	err = gs.Stor.UpdateTemplate(tmp)
	if err != nil {
		return err
	}
	reply.Message = "Template was successfully update"
	reply.Template = tmp
	return nil
}

// DeleteTemplate deletes the template by id
func (gs *TemplateService) DeleteTemplate(rq *http.Request, args *TemplateArgs, reply *TemplateReply) error {
	tmp, err := gs.Stor.GetTemplate(args.ID)
	if err != nil {
		return err
	}
	gs.Stor.DeleteTemplate(tmp)
	reply.Message = "Template with ID " + fmt.Sprint(args.ID) + " was successfully deleted"
	return nil
}
