// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
