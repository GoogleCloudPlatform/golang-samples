// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package bookshelf

import (
	"errors"
	"log"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/sessions"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
	"google.golang.org/cloud/pubsub"
	"google.golang.org/cloud/storage"
)

var (
	DB          BookDatabase
	OAuthConfig *oauth2.Config

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

	SessionStore sessions.Store

	pubSubCtx context.Context

	// Force import of mgo library.
	_ mgo.Session
)

const PubSubTopic = "fill-book-details"

func init() {
	var err error

	// To use the in-memory test database, uncomment the next line.
	DB = newMemoryDB()

	// [START cloudsql]
	// To use MySQL, uncomment the following lines, and update the username,
	// password and host.
	//
	// DB, err = newMySQLDB(MySQLConfig{
	// 	Username: "",
	// 	Password: "",
	// 	Host:     "",
	// 	Port:     3306,
	// })
	// [END cloudsql]

	// [START mongo]
	// To use Mongo, uncomment the next lines and update the address string and
	// optionally, the credentials.
	//
	// var cred *mgo.Credential
	// DB, err = newMongoDB("localhost", cred)
	// [END mongo]

	// [START datastore]
	// To use Cloud Datastore, uncomment the following lines and update the
	// project ID.
	// More options can be set, see the google package docs for details:
	// http://godoc.org/golang.org/x/oauth2/google
	//
	// DB, err = configureDatastoreDB("<your-project-id>")
	// [END datastore]

	if err != nil {
		log.Fatal(err)
	}

	// [START storage]
	// To configure Cloud Storage, uncomment the following lines and update the
	// bucket name.
	//
	// StorageBucketName = "<your-storage-bucket>"
	// StorageBucket, err = configureStorage(StorageBucketName)
	// [END storage]

	if err != nil {
		log.Fatal(err)
	}

	// [START auth]
	// To enable user sign-in, uncomment the following lines and update the
	// Client ID and Client Secret.
	// You will also need to update OAUTH2_CALLBACK in app.yaml when pushing to
	// production.
	//
	// OAuthConfig = configureOAuthClient("clientid", "clientsecret")
	// [END auth]

	// [START sessions]
	// Configure storage method for session-wide information.
	// Update "something-very-secret" with a hard to guess string or byte sequence.
	cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	SessionStore = cookieStore
	// [END sessions]

	// [START pubsub]
	// To configure Pub/Sub, uncomment the following lines and update the project ID.
	//
	// pubSubCtx, err = cloudContext("<your-project-id>")
	// [END pubsub]

	if err != nil {
		log.Fatal(err)
	}
}

func PubSubEnabled() bool {
	return pubSubCtx != nil
}

// PubSubCtx returns the Pub/Sub context, or an error if Pub/Sub is not
// configured or misconfigured.
func PubSubCtx() (context.Context, error) {
	if pubSubCtx == nil {
		return nil, errors.New("You must configure Pub/Sub in bookshelf/config.go " +
			"before running the Pub/Sub worker.")
	}

	if _, ok := DB.(*memoryDB); ok {
		return nil, errors.New("Pub/Sub worker doesn't work with the in-memory DB " +
			"(worker does not share its memory as the main app). Configure another " +
			"database in bookshelf/config.go first (e.g. MySQL, Cloud Datastore, etc)")
	}

	return pubSubCtx, nil
}

func configureDatastoreDB(projectID string) (BookDatabase, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return newDatastoreDB(client)
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}

func cloudContext(projectID string) (context.Context, error) {
	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, storage.ScopeFullControl, pubsub.ScopePubSub)
	if err != nil {
		return nil, err
	}
	return cloud.WithContext(ctx, projectID, httpClient), nil
}

func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
	redirectURL := os.Getenv("OAUTH2_CALLBACK")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/oauth2callback"
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
