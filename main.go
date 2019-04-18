package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wibl/webapp/api"
	"github.com/wibl/webapp/model"
	"github.com/wibl/webapp/storage"

	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := log.New(os.Stdout, "Log Message ", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println("")

	stor := initializeStorage()

	group1 := &model.Group{Title: "test_group1"}
	err := stor.CreateGroup(group1)
	if err != nil {
		log.Fatal(err)
	}
	stor.CreateTemplate(&model.Template{GroupID: group1.ID, Title: "template1_in_group1"})
	stor.CreateTemplate(&model.Template{GroupID: group1.ID, Title: "template2_in_group1"})

	group2 := &model.Group{Title: "test_group2"}
	stor.CreateGroup(group2)

	group2template1 := &model.Template{GroupID: group2.ID, Title: "template1_in_group2"}
	stor.CreateTemplate(group2template1)
	stor.CreateTemplate(&model.Template{GroupID: group2.ID, Title: "template2_in_group2"})

	printAllGroups(stor)

	groups, _ := stor.GetAllGroups()
	for _, group := range groups {
		group.Title = group.Title + "_new"
		//stor.UpdateGroup(group)
		templates, _ := stor.GetAllTemplates(group)
		for _, template := range templates {
			template.Body = "Body of " + template.Title
			//stor.UpdateTemplate(template)
		}
	}

	printAllGroups(stor)

	stor.DeleteGroup(group1)

	stor.DeleteTemplate(group2template1)

	printAllGroups(stor)

	initRPC(stor)
}

func initializeStorage() storage.Storage {
	stor, err := storage.NewDbStorage("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		panic(err)
	}
	return stor
}

func printAllGroups(stor storage.Storage) {
	fmt.Println("-----------------------------------")
	groups, _ := stor.GetAllGroups()
	for _, group := range groups {
		fmt.Printf("%+v\n", group)
		templates, _ := stor.GetAllTemplates(group)
		for _, template := range templates {
			fmt.Printf("%+v\n", template)
		}
	}
}

func initRPC(stor storage.Storage) {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(&api.GroupService{Stor: stor}, "GS")
	http.Handle("/api", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
