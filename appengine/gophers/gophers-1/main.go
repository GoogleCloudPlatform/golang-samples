package main

// [START import_statements]
import (
	"fmt"
	"net/http"

	"google.golang.org/appengine" // Required external App Engine library
)

// [END import_statements]
// [START main_func]
func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}

// [END main_func]
// [START indexHandler]
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// if statement redirects all invalid URLs to the root homepage.
	// Ex: if URL is http://[YOUR_PROJECT_ID].appspot.com/FOO, it will be
	// redirected to http://[YOUR_PROJECT_ID].appspot.com.
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Fprintln(w, "Hello, Gopher Network!")
}

// [END indexHandler]
