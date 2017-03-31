package devserver

import (
    "appengine"
    "fmt"
    "net/http"
)

func init() {
    http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `Development server? %t`, appengine.IsDevAppServer())
}
