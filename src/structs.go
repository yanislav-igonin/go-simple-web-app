package main

import (
	"io/ioutil"
	"path/filepath"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	path := filepath.Join(PAGES_DIR_PATH, filename)
	return ioutil.WriteFile(path, p.Body, 0600)
}

type ErrorData struct {
	Message string
}
