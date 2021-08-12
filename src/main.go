package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(
	template.ParseFiles(
		TEMPLATES_DIR_PATH+"/edit.html",
		TEMPLATES_DIR_PATH+"/view.html",
		TEMPLATES_DIR_PATH+"/index.html",
		TEMPLATES_DIR_PATH+"/error.html",
	),
)
var validPath = regexp.MustCompile(`^/(edit|save|view)/(\S+)$`)

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
