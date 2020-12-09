package security

import (
	"bytes"
	"log"
	"net/http"
	"os"
)

// MakeGetRequestFixture is an HTTP Cloud Function to facilitate testing security snippets.
func MakeGetRequestFixture(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("TARGET_URL") == "" {
		w.Write([]byte("Success!"))
		return
	}

	var b bytes.Buffer
	if err := makeGetRequest(&b, os.Getenv("TARGET_URL")); err != nil {
		log.Printf("makeGetRequest: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(b.Bytes())
}
