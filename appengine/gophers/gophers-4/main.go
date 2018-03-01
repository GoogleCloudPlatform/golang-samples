// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	// [START imports]
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	// [END imports]
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

// [START post_struct]
type Post struct {
	Author  string
	Message string
	Posted  time.Time
}

// [END post_struct]

type templateParams struct {
	Notice string

	Name string
	// [START added_templateParams_fields]
	Message string

	Posts []Post
	// [END added_templateParams_fields]

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
	// [START new_context]
	ctx := appengine.NewContext(r)
	// [END new_context]
	params := templateParams{}

	// [START new_query]
	q := datastore.NewQuery("Post").Order("-Posted").Limit(20)
	// [END new_query]
	// [START get_posts]
	if _, err := q.GetAll(ctx, &params.Posts); err != nil {
		log.Errorf(ctx, "Getting posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest posts. Refresh?"
		indexTemplate.Execute(w, params)
		return
	}
	// [END get_posts]

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}

	// It's a POST request, so handle the form submission.
	// [START new_post]
	post := Post{
		Author:  r.FormValue("name"),
		Message: r.FormValue("message"),
		Posted:  time.Now(),
	}
	// [END new_post]
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
	// [START new_key]
	key := datastore.NewIncompleteKey(ctx, "Post", nil)
	// [END new_key]
	// [START add_post]
	if _, err := datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't add new post. Try again?"
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// [END add_post]

	// Prepend the post that was just added.
	// [START prepend_post]
	params.Posts = append([]Post{post}, params.Posts...)
	// [END prepend_post]

	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", post.Author)
	indexTemplate.Execute(w, params)
}
