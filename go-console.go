package main

import (
	"log"
	"errors"
	"os"
	"regexp"
	"io"
	"io/ioutil"
	"net/http"
	"html/template"
)

var templates = template.Must(template.ParseFiles("view.html"))
var validPath = regexp.MustCompile("^/console/([a-zA-Z0-9]+.html)$")
var (
	Trace *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[1], nil // The title is the second subexpression.
}

func loadPage(title string) (*Page, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}


func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		Info.Println("Wrong page title")
		return
	}
	p, err := loadPage(title)
	if err != nil {
		Info.Println("Can't load "+title+", redirecting to index.html")
		http.Redirect(w, r, "/console/index.html", http.StatusFound)
	} else {
		Info.Println("Rendering index.html")
		renderTemplate(w, "view", p)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	p1 := &Page{Title: "index.html", Body: []byte("This is the first page.")}
	p1.save()
	Info.Println("Created index.html")
        p2 := &Page{Title: "index2.html", Body: []byte("This is the second page.")}
        p2.save()
	Info.Println("Created index2.html")
	http.HandleFunc("/console/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
