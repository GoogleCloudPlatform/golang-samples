// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package bookshelf

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DB          BookDatabase
	OAuthConfig *oauth2.Config

	PubsubTopicID string

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

	SessionStore sessions.Store

	PubsubClient *pubsub.Client

	// Force import of mgo library.
	_ mgo.Session
)

func init() {
	var err error

	// Read config.json
	configFile, err := ioutil.ReadFile("../config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Parse config.json
	var configJson map[string]interface{}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	// To use the in-memory test database, uncomment the next line.
	DB = newMemoryDB()

	// [START cloudsql]
	// When running locally, localhost:3306 is used, and the instance name is
  // ignored.
	if configJson["DATA_BACKEND"].(string) == "cloudsql" {
		DB, err = configureCloudSQL(cloudSQLConfig{
			Username: configJson["MYSQL_USER"].(string),
			Password: configJson["MYSQL_PASSWORD"].(string),
			// The connection name of the Cloud SQL v2 instance, i.e.,
			// "project:region:instance-id"
			// Cloud SQL v1 instances are not supported.
			Instance: configJson["INSTANCE_CONNECTION_NAME"].(string),
		})
	}
	// [END cloudsql]

	// [START mongo]
	// You can optionally update the credentials.
	if configJson["DATA_BACKEND"].(string) == "mongodb" {
		var cred *mgo.Credential
		DB, err = newMongoDB(configJson["MONGO_URL"].(string), cred)
	}
	// [END mongo]

	// [START datastore]
	if configJson["DATA_BACKEND"].(string) == "datastore" {
		// More options can be set, see the google package docs for details:
		// http://godoc.org/golang.org/x/oauth2/google
		DB, err = configureDatastoreDB(configJson["GCLOUD_PROJECT"].(string))
	}
	// [END datastore]

	if err != nil {
		log.Fatal(err)
	}

	// [START storage]
	StorageBucketName = configJson["CLOUD_BUCKET"].(string)
	StorageBucket, err = configureStorage(StorageBucketName)
	// [END storage]

	if err != nil {
		log.Fatal(err)
	}

	// [START auth]
	// To enable user sign-in, uncomment the following lines.
	// You will also need to update OAUTH2_CALLBACK in config.json when pushing to
	// production.
	//
	// OAuthConfig = configureOAuthClient(
	// 	configJson["OAUTH2_CLIENT_ID"].(string),
	// 	configJson["OAUTH2_CLIENT_SECRET"].(string),
	//	configJson["OAUTH2_CALLBACK"].(string)
	// )
	// [END auth]

	// [START sessions]
	// Configure storage method for session-wide information.
	cookieStore := sessions.NewCookieStore([]byte(configJson["COOKIE_STORE_SECRET"].(string)))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	SessionStore = cookieStore
	// [END sessions]

	// [START pubsub]
	PubsubTopicID = configJson["TOPIC_NAME"].(string)

	// To configure Pub/Sub, uncomment the following lines.
	//
	// PubsubClient, err = configurePubsub(
	// 	configJson["GCLOUD_PROJECT"].(string),
	// 	configJson["TOPIC_NAME"].(string)
	// )
	// [END pubsub]

	if err != nil {
		log.Fatal(err)
	}
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

func configurePubsub(projectID string, topicID string) (*pubsub.Client, error) {
	if _, ok := DB.(*memoryDB); ok {
		return nil, errors.New("Pub/Sub worker doesn't work with the in-memory DB " +
			"(worker does not share its memory as the main app). Configure another " +
			"database in bookshelf/config.go first (e.g. MySQL, Cloud Datastore, etc)")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Create the topic if it doesn't exist.
	if exists, err := client.Topic(topicID).Exists(ctx); err != nil {
		return nil, err
	} else if !exists {
		if _, err := client.CreateTopic(ctx, topicID); err != nil {
			return nil, err
		}
	}
	return client, nil
}

func configureOAuthClient(clientID, clientSecret string, redirectURL string) *oauth2.Config {
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

type cloudSQLConfig struct {
	Username, Password, Instance string
}

func configureCloudSQL(config cloudSQLConfig) (BookDatabase, error) {
	if os.Getenv("GAE_INSTANCE") != "" {
		// Running in production.
		return newMySQLDB(MySQLConfig{
			Username:   config.Username,
			Password:   config.Password,
			UnixSocket: "/cloudsql/" + config.Instance,
		})
	}

	// Running locally.
	return newMySQLDB(MySQLConfig{
		Username: config.Username,
		Password: config.Password,
		Host:     "localhost",
		Port:     3306,
	})
}
