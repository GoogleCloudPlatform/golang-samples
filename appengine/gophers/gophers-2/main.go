package main

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	indexTemplate.Execute(w, nil)
}
