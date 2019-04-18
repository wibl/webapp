package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wibl/webapp/api"
	"github.com/wibl/webapp/model"
	"github.com/wibl/webapp/mq"
	"github.com/wibl/webapp/storage"

	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func main() {
	logger := log.New(os.Stdout, "Log Message ", log.Ldate|log.Ltime|log.Lshortfile)

	stompSender, err := mq.New("localhost:61613")
	if err != nil {
		logger.Fatal(err)
	}
	defer stompSender.Disconnect()

	context := &appContext{
		sender: stompSender,
	}

	sendTestMessage(context)

	stor := initializeStorage()

	group1 := &model.Group{Title: "test_group1"}
	stor.CreateGroup(group1)
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

type appContext struct {
	sender mq.Sender
}

func sendTestMessage(context *appContext) error {
	return context.sender.SendMessage("/queue/test-1", "TEST")
}

func initializeStorage() storage.Storage {
	//stor, _ := storage.NewDb("", "")
	stor, _ := storage.NewMemStorage()
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

// func initRPC() {
// 	s := rpc.NewServer()
// 	s.RegisterCodec(json.NewCodec(), "application/json")
// 	s.RegisterService(new(HelloService), "")
// 	http.Handle("/rpc", s)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}
