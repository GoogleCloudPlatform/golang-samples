package main

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

type templateParams struct {
	Notice string
	Name   string
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	params := templateParams{}

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}

	// It's a POST request, so handle the form submission.

	name := r.FormValue("name")
	params.Name = name // Preserve the name field.
	if name == "" {
		name = "Anonymous Gopher"
	}

	if r.FormValue("message") == "" {
		w.WriteHeader(http.StatusBadRequest)

		params.Notice = "No message provided"
		indexTemplate.Execute(w, params)
		return
	}

	// TODO: save the message into a database.

	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", name)
	indexTemplate.Execute(w, params)
}
