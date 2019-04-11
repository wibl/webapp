package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wibl/webapp/model"
	"github.com/wibl/webapp/mq"
	"github.com/wibl/webapp/storage"
)

func initializeStorage() storage.Storage {
	//stor, _ := storage.NewDb("", "")
	stor, _ := storage.NewMemStorage()
	return stor
}

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

	groups, _ := stor.GetGroups()
	for _, group := range groups {
		group.Title = group.Title + "_new"
		stor.UpdateGroup(group)
		templates, _ := stor.GetTemplates(group)
		for _, template := range templates {
			template.Body = "Body of " + template.Title
			stor.UpdateTemplate(template)
		}
	}

	printAllGroups(stor)

	stor.DeleteGroup(group1)

	stor.DeleteTemplate(group2template1)

	printAllGroups(stor)

}

func printAllGroups(stor storage.Storage) {
	fmt.Println("-----------------------------------")
	groups, _ := stor.GetGroups()
	for _, group := range groups {
		fmt.Printf("%+v\n", group)
		templates, _ := stor.GetTemplates(group)
		for _, template := range templates {
			fmt.Printf("%+v\n", template)
		}
	}
}

type appContext struct {
	sender mq.Sender
}

func sendTestMessage(context *appContext) error {
	return context.sender.SendMessage("/queue/test-1", "TEST")
}

/* Deprecated
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Page as a struct with two fields representing the title and body.
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(wr http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	if name == "" {
		name = "MainPage"
	}
	fmt.Fprintf(wr, "Hi, there, it's %s!", name)
}

func viewHandler(wr http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(wr, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(wr, "view", p)
}

func editHandler(wr http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(wr, "edit", p)
}

func saveHandler(wr http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(wr, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(wr http.ResponseWriter, tmpl string, page *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(wr, page)
}
*/
