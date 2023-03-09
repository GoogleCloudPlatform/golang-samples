// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// MarkdownRenderer defines an interface for rendering Markdown to HTML.
type MarkdownRenderer interface {
	Render([]byte) ([]byte, error)
}

// Service manages centralized resources of the service
type Service struct {
	Renderer MarkdownRenderer

	parsedTemplate  *template.Template
	markdownDefault string
}

// NewServiceFromEnv creates a new Service instance from environment variables.
func NewServiceFromEnv() (*Service, error) {
	url := os.Getenv("EDITOR_UPSTREAM_RENDER_URL")
	if url == "" {
		return nil, errors.New("no configuration for upstream render service: add EDITOR_UPSTREAM_RENDER_URL environment variable")
	}

	// The use case of this service is the UI driven by these files.
	// Loading them as part of the server startup process keeps failures easily
	// discoverable and minimizes latency for the first request.
	parsedTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("template.ParseFiles: %w", err)
	}

	out, err := ioutil.ReadFile("templates/markdown.md")
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %w", err)
	}
	markdownDefault := string(out)

	return &Service{
		Renderer: &RenderService{
			URL: url,
		},
		parsedTemplate:  parsedTemplate,
		markdownDefault: markdownDefault,
	}, nil
}

// RegisterHandlers registers all HTTP handler routes to a new ServeMux.
func (s *Service) RegisterHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.editorHandler)
	mux.HandleFunc("/render", s.renderHandler)

	return mux
}

func (s *Service) editorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	p := map[string]string{
		"Default": s.markdownDefault,
	}

	if err := s.parsedTemplate.Execute(w, p); err != nil {
		log.Printf("template.Execute: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// renderHandler expects a JSON body payload with a 'data' property holding plain text for rendering.
func (s *Service) renderHandler(w http.ResponseWriter, r *http.Request) {
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

	rendered, err := s.Renderer.Render([]byte(d.Data))
	if err != nil {
		log.Printf("MarkdownRenderer.Render: %v", err)
		msg := http.StatusText(http.StatusInternalServerError)
		if strings.Contains(err.Error(), "http.Client.Do") {
			msg = fmt.Sprintf("<h3>%s (%d)</h3>\n<p>The request to the upstream render service failed with the message:</p>\n<p>%s</p>", http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, rendered)
		}

		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Write(rendered)
}
