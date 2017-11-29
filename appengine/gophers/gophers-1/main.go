package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	indexTemplate.Execute(w, nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	body := struct{ Message string }{}

	if r.FormValue("message") == "" {
		body.Message = "No message provided"
		writeJSON(w, http.StatusBadRequest, body)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous Gopher"
	}

	// TODO: save the message into a database.

	body.Message = fmt.Sprintf("Thank you for your submission, %s!", name)

	writeJSON(w, http.StatusOK, body)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		http.Error(w, `{"Error":"Could not marshal payload."}`, http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}
