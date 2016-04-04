// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
