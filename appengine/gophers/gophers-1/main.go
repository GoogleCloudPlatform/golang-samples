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
		Response(body).Status(http.StatusBadRequest).WriteJSON(w)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous Gopher"
	}

	// TODO: save the message into a database.

	body.Message = fmt.Sprintf("Thank you for your submission, %s!", name)

	Response(body).WriteJSON(w)
}

func Response(payload interface{}) response {
	return response{payload: payload}
}

type response struct {
	status  int
	payload interface{}
}

func (r response) Status(status int) response {
	r.status = status
	return r
}

func (r response) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(r.payload, "", "\t")
	if err != nil {
		http.Error(w, `{"Error":"Could not marshal payload."}`, http.StatusInternalServerError)
		return err
	}
	if r.status != 0 {
		w.WriteHeader(r.status)
	}
	_, err = w.Write(b)
	return err
}
