package api

import (
	"fmt"
	"net/http"

	"github.com/wibl/webapp/model"
	"github.com/wibl/webapp/storage"
)

// GroupService contains methods for working with groups
type GroupService struct {
	Stor storage.Storage
}

// GroupArgs contains possible arguments
type GroupArgs struct {
	ID    int64
	Title string
}

// GroupReply contains info for replay
type GroupReply struct {
	Message string
	Groups  []*model.Group
}

// CreateGroup allows to create a group
func (gs *GroupService) CreateGroup(rq *http.Request, args *GroupArgs, reply *GroupReply) error {
	newGroup := &model.Group{Title: args.Title}
	err := gs.Stor.CreateGroup(newGroup)
	if err != nil {
		return err
	}
	reply.Message = "Created Group " + args.Title
	return nil
}

// GetAllGroups return all groups
func (gs *GroupService) GetAllGroups(rq *http.Request, args *GroupArgs, reply *GroupReply) error {
	groups, err := gs.Stor.GetAllGroups()
	if err != nil {
		return err
	}
	reply.Message = "Groups successfully received"
	reply.Groups = groups
	return nil
}

// DeleteGroup deletes the group by id
func (gs *GroupService) DeleteGroup(rq *http.Request, args *GroupArgs, reply *GroupReply) error {
	group, err := gs.Stor.GetGroup(args.ID)
	if err != nil {
		return err
	}
	gs.Stor.DeleteGroup(group)
	reply.Message = "Group with ID " + fmt.Sprint(args.ID) + " was successfully deleted"
	return nil
}

// UpdateGroup updates the group title
func (gs *GroupService) UpdateGroup(rq *http.Request, args *GroupArgs, reply *GroupReply) error {
	group, err := gs.Stor.GetGroup(args.ID)
	if err != nil {
		return err
	}
	group.Title = args.Title
	reply.Message = "Group was successfully update"
	return nil
}
