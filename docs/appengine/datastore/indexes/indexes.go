// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// [START unindexed_properties]
type Person struct {
	Name string
	Age  int `datastore:",noindex"`
}

// [END unindexed_properties]

// [START exploding_index_example_3]
type Widget struct {
	X    []int
	Y    []string
	Date time.Time
}

func f(ctx context.Context) {
	e2 := &Widget{
		X:    []int{1, 2, 3, 4},
		Y:    []string{"red", "green", "blue"},
		Date: time.Now(),
	}

	k := datastore.NewIncompleteKey(ctx, "Widget", nil)
	if _, err := datastore.Put(ctx, k, e2); err != nil {
		// Handle error.
	}
}

// [END exploding_index_example_3]
