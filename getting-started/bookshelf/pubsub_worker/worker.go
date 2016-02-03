// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/api/books/v1"
	"google.golang.org/cloud/pubsub"

	"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf"
)

const subName = "book-worker-sub"

var (
	countMu sync.Mutex
	count   int

	booksClient *books.Service
	pubSubCtx   context.Context
)

func main() {
	var err error

	pubSubCtx, err = bookshelf.PubSubCtx()
	if err != nil {
		log.Fatal(err)
	}

	booksClient, err = books.New(http.DefaultClient)
	if err != nil {
		log.Fatalf("could not access Google Books API: %v", err)
	}

	// ignore returned errors, which will be "already exists". If they're fatal
	// errors, then following calls (e.g. in the subscribe function) will also fail.
	pubsub.CreateTopic(pubSubCtx, bookshelf.PubSubTopic)
	pubsub.CreateSub(pubSubCtx, subName, bookshelf.PubSubTopic, 0, "")

	// Start worker goroutine.
	go subscribe()

	// [START http]
	// Publish a count of processed requests to the server homepage.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		countMu.Lock()
		defer countMu.Unlock()
		fmt.Fprintf(w, "This worker has processed %d books.", count)
	})

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
	// [END http]
}

func subscribe() {
	for {
		// Pull up to 10 messages (maybe fewer) from the subscription.
		// Blocks for an indeterminate amount of time.
		msgs, err := pubsub.PullWait(pubSubCtx, subName, 10)
		if err != nil {
			log.Fatalf("could not pull: %v", err)
		}

		for _, m := range msgs {
			msg := m

			var id int64
			if err := json.Unmarshal(msg.Data, &id); err != nil {
				log.Printf("could not decode message data: %#v", msg)
				go pubsub.Ack(pubSubCtx, subName, msg.AckID)
				continue
			}

			log.Printf("[ID %d] Processing.", id)
			go func() {
				if err := update(id); err != nil {
					log.Printf("[ID %d] could not update: %v", id, err)
					return
				}

				countMu.Lock()
				count++
				countMu.Unlock()

				pubsub.Ack(pubSubCtx, subName, msg.AckID)
				log.Printf("[ID %d] ACK", id)
			}()
		}
	}
}

// update retrieves the book with the given ID, finds metata from the Books
// server and updates the database with the book's details.
func update(bookID int64) error {
	book, err := bookshelf.DB.GetBook(bookID)
	if err != nil {
		return err
	}

	vols, err := booksClient.Volumes.List(book.Title).Do()
	if err != nil {
		return err
	}

	if len(vols.Items) == 0 {
		return nil
	}

	info := vols.Items[0].VolumeInfo
	book.Title = info.Title
	book.Author = strings.Join(info.Authors, ", ")
	book.PublishedDate = info.PublishedDate
	if book.Description == "" {
		book.Description = info.Description
	}
	if book.ImageURL == "" && info.ImageLinks != nil {
		url := info.ImageLinks.Thumbnail
		// Replace http with https to prevent Content Security errors on the page.
		book.ImageURL = strings.Replace(url, "http://", "https://", 1)
	}

	return bookshelf.DB.UpdateBook(book)
}
