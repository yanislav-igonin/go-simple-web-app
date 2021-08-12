package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const PAGES_DIR_PATH = "../pages"
const TEMPLATES_DIR_PATH = "../templates"

var templates = template.Must(
	template.ParseFiles(
		TEMPLATES_DIR_PATH+"/edit.html",
		TEMPLATES_DIR_PATH+"/view.html",
		TEMPLATES_DIR_PATH+"/index.html",
		TEMPLATES_DIR_PATH+"/error.html",
	),
)
var validPath = regexp.MustCompile(`^/(edit|save|view)/(\S+)$`)

type Page struct {
	Title string
	Body  []byte
}

type ErrorData struct {
	Message string
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	path := filepath.Join(PAGES_DIR_PATH, filename)
	return ioutil.WriteFile(path, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	path := filepath.Join(PAGES_DIR_PATH, filename)
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func createDir() error {
	return os.MkdirAll(PAGES_DIR_PATH, os.ModePerm)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(PAGES_DIR_PATH)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var filanames []string
	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		filanames = append(filanames, name)
	}

	renderTemplate(w, "index", filanames)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	err := createDir()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
