package devserver

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IsDevAppServer: %v", appengine.IsDevAppServer())
}
