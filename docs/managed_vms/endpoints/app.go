package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"google.golang.org/appengine"
)

func main() {
	r := mux.NewRouter()

	r.Path("/echo").Methods("POST").
		HandlerFunc(echoHandler)

	r.Path("/auth/info/googlejwt").Methods("GET").
		HandlerFunc(authInfoHandler)
	r.Path("/auth/info/googleidtoken").Methods("GET").
		HandlerFunc(authInfoHandler)

	http.Handle("/", r)
	appengine.Main()
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	var msg interface{}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			errorf(w, http.StatusBadRequest, "Body was not valid JSON: %v", err)
			return
		}
		errorf(w, http.StatusInternalServerError, "Could not get body: %v", err)
		return
	}

	b, err := json.Marshal(msg)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Could not marshal JSON: %v", err)
		return
	}
	w.Write(b)
}

func authInfoHandler(w http.ResponseWriter, r *http.Request) {
	encodedInfo := r.Header.Get("X-Endpoint-API-UserInfo")
	if encodedInfo == "" {
		w.Write([]byte(`{"id": "anonymous"}`))
		return
	}

	b, err := base64.StdEncoding.DecodeString(encodedInfo)
	if err != nil {
		errorf(w, http.StatusInternalServerError, "Could not decode auth info: %v", err)
		return
	}
	w.Write(b)
}

// errorf writes a swagger-compliant error response.
func errorf(w http.ResponseWriter, code int, format string, a ...interface{}) {
	var error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	error.Code = code
	error.Message = fmt.Sprintf(format, a...)

	b, err := json.Marshal(error)
	if err != nil {
		http.Error(w, `{"code": 500, "message": "Could not format JSON for original message."}`, 500)
		return
	}

	http.Error(w, string(b), code)
}
