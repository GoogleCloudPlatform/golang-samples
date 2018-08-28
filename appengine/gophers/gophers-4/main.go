// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"

	// [START gae_go_env_data_imports]
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	// [END gae_go_env_data_imports]
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

// [START gae_go_env_post_struct]
type Post struct {
	Author  string
	Message string
	Posted  time.Time
}

// [END gae_go_env_post_struct]

type templateParams struct {
	Notice string

	Name string
	// [START gae_go_env_template_params_fields]
	Message string

	Posts []Post
	// [END gae_go_env_template_params_fields]

}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// [START gae_go_env_new_context]
	ctx := appengine.NewContext(r)
	// [END gae_go_env_new_context]
	params := templateParams{}

	// [START gae_go_env_new_query]
	q := datastore.NewQuery("Post").Order("-Posted").Limit(20)
	// [END gae_go_env_new_query]
	// [START gae_go_env_get_posts]
	if _, err := q.GetAll(ctx, &params.Posts); err != nil {
		log.Errorf(ctx, "Getting posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest posts. Refresh?"
		indexTemplate.Execute(w, params)
		return
	}
	// [END gae_go_env_get_posts]

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}

	// It's a POST request, so handle the form submission.
	// [START gae_go_env_new_post]
	post := Post{
		Author:  r.FormValue("name"),
		Message: r.FormValue("message"),
		Posted:  time.Now(),
	}
	// [END gae_go_env_new_post]
	if post.Author == "" {
		post.Author = "Anonymous Gopher"
	}
	params.Name = post.Author

	if post.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "No message provided"
		indexTemplate.Execute(w, params)
		return
	}
	// [START gae_go_env_new_key]
	key := datastore.NewIncompleteKey(ctx, "Post", nil)
	// [END gae_go_env_new_key]
	// [START gae_go_env_add_post]
	if _, err := datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't add new post. Try again?"
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// [END gae_go_env_add_post]

	// Prepend the post that was just added.
	// [START gae_go_env_prepend_post]
	params.Posts = append([]Post{post}, params.Posts...)
	// [END gae_go_env_prepend_post]

	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", post.Author)
	indexTemplate.Execute(w, params)
}
