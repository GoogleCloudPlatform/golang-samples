package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// parsedTemplate is the global parsed HTML template.
var parsedTemplate *template.Template
var markdownDefault string

func init() {
	var err error
	parsedTemplate, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("template.ParseFiles: %v", err)
	}

	out, err := ioutil.ReadFile("templates/markdown.md")
	if err != nil {
		log.Fatalf("ioutil.ReadFile: %v", err)
	}
	markdownDefault = string(out)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if err := parsedTemplate.Execute(w, map[string]string{"Default": markdownDefault}); err != nil {
		log.Printf("template.Execute: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
