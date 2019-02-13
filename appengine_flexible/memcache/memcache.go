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

// [START gae_flex_redislabs_memcache]

// Sample memcache demonstrates use of a memcached client from App Engine flexible environment.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bradfitz/gomemcache/memcache"

	"google.golang.org/appengine"
)

var memcacheClient *memcache.Client

func main() {
	host := os.Getenv("MEMCACHE_PORT_11211_TCP_ADDR")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MEMCACHE_PORT_11211_TCP_PORT")
	if port == "" {
		port = "11211"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	memcacheClient = memcache.New(addr)

	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var count uint64
	var err error

	for {
		count, err = memcacheClient.Increment("count", 1)
		if err == nil {
			break
		}
		if err != memcache.ErrCacheMiss {
			msg := fmt.Sprintf("Could not increment count: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		initial := &memcache.Item{
			Key:   "count",
			Value: []byte("0"),
		}
		err := memcacheClient.Add(initial)
		if err != nil && err != memcache.ErrNotStored {
			msg := fmt.Sprintf("Could not populate initial value: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		// Increment via the next iteration of the loop.
	}

	fmt.Fprintf(w, "Count: %d", count)
}

// [END gae_flex_redislabs_memcache]
