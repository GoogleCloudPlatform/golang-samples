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
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// [START intro_1]
import "google.golang.org/appengine/memcache"

// [END intro_1]

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// [START intro_2]
	// Create an Item
	item := &memcache.Item{
		Key:   "lyric",
		Value: []byte("Oh, give me a home"),
	}
	// Add the item to the memcache, if the key does not already exist
	if err := memcache.Add(ctx, item); err == memcache.ErrNotStored {
		log.Infof(ctx, "item with key %q already exists", item.Key)
	} else if err != nil {
		log.Errorf(ctx, "error adding item: %v", err)
	}

	// Change the Value of the item
	item.Value = []byte("Where the buffalo roam")
	// Set the item, unconditionally
	if err := memcache.Set(ctx, item); err != nil {
		log.Errorf(ctx, "error setting item: %v", err)
	}

	// Get the item from the memcache
	if item, err := memcache.Get(ctx, "lyric"); err == memcache.ErrCacheMiss {
		log.Infof(ctx, "item not in the cache")
	} else if err != nil {
		log.Errorf(ctx, "error getting item: %v", err)
	} else {
		log.Infof(ctx, "the lyric is %q", item.Value)
	}
	// [END intro_2]

}
