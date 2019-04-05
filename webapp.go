package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
	"github.com/go-stomp/stomp"
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

func sendMess(destination, contentType, message string) error {
	netConn, err := net.DialTimeout("tcp", "stomp.server.com:61613", 10*time.Second)
	if err != nil {
		return err
	}
	defer netConn.Close()

	stompConn, err := stomp.Connect(netConn)
	if err != nil {
		return err
	}
	defer stompConn.Disconnect()
	
	stompConn.Send(destination, contentType, []byte(message))
	return nil
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	sendMess("/queue/test-1", "text/plain", "TEST")
}
