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

package app

// [START gae_writing_logs]
import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func logHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	post := &Post{Body: "sample post"}
	key := datastore.NewIncompleteKey(ctx, "Posts", nil)
	if _, err := datastore.Put(ctx, key, post); err != nil {
		log.Errorf(ctx, "could not put into datastore: %v", err)
		http.Error(w, "An error occurred. Try again.", http.StatusInternalServerError)
		return
	}
	log.Debugf(ctx, "Datastore put successful")

	w.Write([]byte("ok!"))
}

// [END gae_writing_logs]

type Post struct {
	Body string
}

func init() {
	http.HandleFunc("/log", logHandler)
}
