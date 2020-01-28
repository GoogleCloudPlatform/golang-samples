package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// MarkdownRenderer defines an interface for rendering Markdown to HTML.
type MarkdownRenderer interface {
	Render([]byte) ([]byte, error)
}

// render converts Markdown text into HTML.
var render MarkdownRenderer

func init() {
	url := os.Getenv("EDITOR_UPSTREAM_RENDER_URL")
	if url == "" {
		log.Fatalf("no configuration for upstream render service: add EDITOR_UPSTREAM_RENDER_URL environment variable")
	}
	auth := os.Getenv("EDITOR_UPSTREAM_UNAUTHENTICATED") == ""
	if !auth {
		log.Println("editor: starting in unauthenticated upstream mode")
	}
	render = &RenderService{
		URL:           url,
		Authenticated: auth,
	}
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	out, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var d struct{ Data string }
	if err := json.Unmarshal(out, &d); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	rendered, err := render.Render([]byte(d.Data))
	if err != nil {
		log.Printf("MarkdownRenderer.Render: %v", err)

		msg := http.StatusText(http.StatusInternalServerError)
		if errors.Is(err, errNotOk) {
			msg = fmt.Sprintf("<h3>%s (%d)</h3>\n<p>The request to the upstream render service failed with the message:</p>\n<p>%s</p>", http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, rendered)
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Write(rendered)
}
