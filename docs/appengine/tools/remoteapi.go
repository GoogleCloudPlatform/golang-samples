// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/remote_api"
)

func main() {
	const host = "<your-app-id>.appspot.com"

	ctx := context.Background()

	hc, err := google.DefaultClient(ctx,
		"https://www.googleapis.com/auth/appengine.apis",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		log.Fatal(err)
	}

	remoteCtx, err := remote_api.NewRemoteContext(host, hc)
	if err != nil {
		log.Fatal(err)
	}

	q := datastore.NewQuery("Greeting").
		Filter("Date >=", time.Now().AddDate(0, 0, -7))

	log.Println(q.Count(remoteCtx))
}
