package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const PAGES_DIR_PATH = "../pages"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("../pages/"+filename, p.Body, 0600)
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

func createDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func main() {
	err := createDir(PAGES_DIR_PATH)
	if err != nil {
		panic(err)
	}

	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	err = p1.save()
	if err != nil {
		panic(err)
	}

	p2, err := loadPage("TestPage")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(p2.Body))
}
