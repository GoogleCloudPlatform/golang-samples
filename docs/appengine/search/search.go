// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

// [START creating_a_document]
import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/search"
)

type User struct {
	Name      string
	Comment   search.HTML
	Visits    float64
	LastVisit time.Time
	Birthday  time.Time
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	id := "PA6-5000"
	user := &User{
		Name:      "Joe Jackson",
		Comment:   "this is <em>marked up</em> text",
		Visits:    7,
		LastVisit: time.Now(),
		Birthday:  time.Date(1960, time.June, 19, 0, 0, 0, 0, nil),
	}
	// ...
	// [END creating_a_document]

	// [START putting_documents_in_an_index_1]
	// ...
	ctx := appengine.NewContext(r)
	index, err := search.Open("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = index.Put(ctx, id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")

	// [END putting_documents_in_an_index_1]

	// [START putting_documents_in_an_index_2]
	id, err = index.Put(ctx, "", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, id)
	// [END putting_documents_in_an_index_2]
}

// [START retrieving_documents_by_doc_ids]
func getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	index, err := search.Open("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := "PA6-5000"
	var user User
	if err := index.Get(ctx, id, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Retrieved document: ", user)
}

// [END retrieving_documents_by_doc_ids]

// [START deleting_documents_from_an_index]
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	index, err := search.Open("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := "PA6-5000"
	err = index.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Deleted document: ", id)
}

// [END deleting_documents_from_an_index]

type Doc struct{}

// [START search_example_1]
func searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	index, err := search.Open("myIndex")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for t := index.Search(ctx, "Product: piano AND Price < 5000", nil); ; {
		var doc Doc
		id, err := t.Next(&doc)
		if err == search.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "Search error: %v\n", err)
			break
		}
		fmt.Fprintf(w, "%s -> %#v\n", id, doc)
	}
}

// [END search_example_1]

func sample() {
	var ctx context.Context
	var index search.Index

	// [START queries_1]
	index.Search(ctx, "rose water", nil)
	// [END queries_1]

	// [START queries_2]
	index.Search(ctx, "1776-07-04", nil)
	// [END queries_2]

	// [START queries_3]
	// search for documents with pianos that cost less than $5000
	index.Search(ctx, "Product = piano AND Price < 5000", nil)
	// [END queries_3]

}
