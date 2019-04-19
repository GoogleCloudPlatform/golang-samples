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

package sample

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// [START example]
	w.Header().Set("Content-Type", "text/plain")
	c := appengine.NewContext(r)

	who := "nobody"
	item, err := memcache.Get(c, "who")
	if err == nil {
		who = string(item.Value)
	} else if err != memcache.ErrCacheMiss {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Previously incremented by %s\n", who)
	memcache.Set(c, &memcache.Item{
		Key:   "who",
		Value: []byte("Go"),
	})

	count, _ := memcache.Increment(c, "count", 1, 0)
	fmt.Fprintf(w, "Count incremented by Go = %d\n", count)
	// [END example]
}

func init() {
	http.HandleFunc("/", handle)
}
