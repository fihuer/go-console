package httpServer

import (
	"os"
	"io"
	"errors"
	"log"
	"net/http"
	"strings"
	"html/template"
	"io/ioutil"
	"regexp"
)

var templates = template.Must(template.ParseFiles("httpServer/templates/console.html"))
var validPath = regexp.MustCompile("^/console/([a-zA-Z0-9-_]+.html)$")
var libPath = regexp.MustCompile("^/bower_components/[-_a-zA-Z0-9/]*([a-zA-Z0-9-_]+.(html|js|css))$")
var elemPath = regexp.MustCompile("^/elements/[-_a-zA-Z0-9/]*([a-zA-Z0-9-_]+.(html|js|css))$")
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

func (p *Page) Save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func Start(cert, key string) error {
	//Init loggers
	initLoggers(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	//Http redirector
	go func() {
		http.ListenAndServe(":8081", http.HandlerFunc(redir))
	}()

	//Handlers
	http.HandleFunc("/console/",viewHandler)
	http.HandleFunc("/bower_components/", libHandler)
	http.HandleFunc("/elements/", elemHandler)

	err := http.ListenAndServeTLS(":8080", cert, key, nil)
	return err
}

func initLoggers(
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

func redir(w http.ResponseWriter, r *http.Request) {
	m := strings.Split(r.Host, ":")
	h := m[0]
	http.Redirect(w, r, "https://"+h+":8080"+r.RequestURI, http.StatusMovedPermanently)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		Info.Println("Wrong page title")
		http.Redirect(w, r, "/console/console.html", http.StatusFound)
		return
	}
	p, err := loadPage(title)
	if err != nil {
		Info.Println("Can't load "+title+", redirecting to console.html")
		http.Redirect(w, r, "/console/console.html", http.StatusFound)
	} else {
		Info.Println("Rendering console.html")
		renderTemplate(w, "console", p)
	}
}

func libHandler(w http.ResponseWriter, r *http.Request) {
	m := libPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		Warning.Println("Invalid Page Title : "+r.URL.Path)
		return
	}
	Info.Println("Loading Lib : "+r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}
func elemHandler(w http.ResponseWriter, r *http.Request) {
	m := elemPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		Warning.Println("Invalid Page Title : "+r.URL.Path)
		return
	}
	Info.Println("Loading Lib : "+r.URL.Path[1:])
	http.ServeFile(w, r, "httpServer/"+r.URL.Path[1:])
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		//http.NotFound(w, r)
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	w.Header().Set("Content-type", "text/html")
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
