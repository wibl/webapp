package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/wibl/webapp/mq"
)

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

type appContext struct {
	sender mq.Sender
}

func sendTestMessage(context *appContext) error {
	return context.sender.SendMessage("/queue/test-1", "TEST")
}

func main() {
	stompSender, err := mq.New("localhost:61613")
	if err != nil {
		panic(err)
	}

	context := &appContext{
		sender: stompSender,
	}

	defer stompSender.Disconnect()

	/*http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))*/

	sendTestMessage(context)
}
