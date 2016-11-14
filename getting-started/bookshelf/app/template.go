// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf"
)

// parseTemplate applies a given file to the body of the base template.
func parseTemplate(filename string) *appTemplate {
	tmpl := template.Must(template.ParseFiles("templates/base.html"))

	// Put the named file into a template called "body"
	path := filepath.Join("templates", filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not read template: %v", err))
	}
	template.Must(tmpl.New("body").Parse(string(b)))

	return &appTemplate{tmpl.Lookup("base.html")}
}

// appTemplate is a user login-aware wrapper for a html/template.
type appTemplate struct {
	t *template.Template
}

// Execute writes the template using the provided data, adding login and user
// information to the base template.
func (tmpl *appTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}) *appError {
	d := struct {
		Data        interface{}
		AuthEnabled bool
		Profile     *Profile
		LoginURL    string
		LogoutURL   string
	}{
		Data:        data,
		AuthEnabled: bookshelf.OAuthConfig != nil,
		LoginURL:    "/login?redirect=" + r.URL.RequestURI(),
		LogoutURL:   "/logout?redirect=" + r.URL.RequestURI(),
	}

	if d.AuthEnabled {
		// Ignore any errors.
		d.Profile = profileFromSession(r)
	}

	if err := tmpl.t.Execute(w, d); err != nil {
		return appErrorf(err, "could not write template: %v")
	}
	return nil
}
