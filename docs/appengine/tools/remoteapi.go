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

package main

import (
	"context"
	"log"
	"time"

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
