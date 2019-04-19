// Copyright 2019 Google LLC
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

// Sample pubsub_worker demonstrates the use of the Cloud Pub/Sub API to communicate between two modules.
// See https://cloud.google.com/go/getting-started/using-pub-sub
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	books "google.golang.org/api/books/v1"

	"cloud.google.com/go/pubsub"

	"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf"
)

const subName = "book-worker-sub"

var (
	countMu sync.Mutex
	count   int

	booksClient  *books.Service
	subscription *pubsub.Subscription
)

func main() {
	ctx := context.Background()

	if bookshelf.PubsubClient == nil {
		log.Fatal("You must configure the Pub/Sub client in config.go before running pubsub_worker.")
	}

	var err error
	booksClient, err = books.New(http.DefaultClient)
	if err != nil {
		log.Fatalf("could not access Google Books API: %v", err)
	}

	// [START pubsub_create_topic]
	// Create pubsub topic if it does not yet exist.
	topic := bookshelf.PubsubClient.Topic(bookshelf.PubsubTopicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Error checking for topic: %v", err)
	}
	if !exists {
		if _, err := bookshelf.PubsubClient.CreateTopic(ctx, bookshelf.PubsubTopicID); err != nil {
			log.Fatalf("Failed to create topic: %v", err)
		}
	}

	// Create topic subscription if it does not yet exist.
	subscription = bookshelf.PubsubClient.Subscription(subName)
	exists, err = subscription.Exists(ctx)
	if err != nil {
		log.Fatalf("Error checking for subscription: %v", err)
	}
	if !exists {
		if _, err = bookshelf.PubsubClient.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{Topic: topic}); err != nil {
			log.Fatalf("Failed to create subscription: %v", err)
		}
	}
	// [END pubsub_create_topic]

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
	ctx := context.Background()
	err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		var id int64
		if err := json.Unmarshal(msg.Data, &id); err != nil {
			log.Printf("could not decode message data: %#v", msg)
			msg.Ack()
			return
		}

		log.Printf("[ID %d] Processing.", id)
		if err := update(id); err != nil {
			log.Printf("[ID %d] could not update: %v", id, err)
			msg.Nack()
			return
		}

		countMu.Lock()
		count++
		countMu.Unlock()

		msg.Ack()
		log.Printf("[ID %d] ACK", id)
	})
	if err != nil {
		log.Fatal(err)
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
