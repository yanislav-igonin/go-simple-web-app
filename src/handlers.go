package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

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
