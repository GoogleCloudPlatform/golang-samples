// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// This sample demonstrates AccessTokenCredentials:
// https://godoc.org/golang.org/x/oauth2/google#JWTAccessTokenSourceFromJSON
// https://developers.google.com/identity/protocols/OAuth2ServiceAccount#jwt-auth

// To use, create a service accountJSON file and allow it atleast
// Cloud Datastore Viewer IAM on the project.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	projectID = flag.String("project", "", "Project ID")
	keyfile   = flag.String("keyfile", "", "Service Account JSON keyfile")
)

func main() {
	flag.Parse()

	if *projectID == "" {
		fmt.Fprintln(os.Stderr, "missing -project flag")
		flag.Usage()
		os.Exit(2)
	}
	if *keyfile == "" {
		fmt.Fprintln(os.Stderr, "missing -keyfile flag")
		flag.Usage()
		os.Exit(2)
	}

	// [START jwtaccesstoken_sample]
	// audience values for other services can be found in the repo here similar to
	// FireStore
	// https://github.com/googleapis/googleapis/blob/master/google/firestore/firestore.yaml
	var aud string = "https://firestore.googleapis.com/google.firestore.v1beta1.Firestore"

	ctx := context.Background()
	keyBytes, err := ioutil.ReadFile(*keyfile)
	if err != nil {
		log.Fatalf("Unable to read service account key file  %v", err)
	}
	tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, aud)
	if err != nil {
		log.Fatalf("Error building JWT access token source: %v", err)
	}
	jwt, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Unable to generate JWT token: %v", err)
	}
	log.Println(jwt.AccessToken)

	client, err := firestore.NewClient(ctx, *projectID, option.WithTokenSource(tokenSource))
	if err != nil {
		log.Fatalf("Could not create FireStore Client: %v", err)
	}

	collections, err := client.Collections(ctx).GetAll()
	if err != nil {
		log.Fatalf("Unable to get Collections %v", err)
	}
	for _, col := range collections {
		log.Printf("Collection ID:  %v", col.ID)
	}
	// [END jwtaccesstoken_sample]
}
