// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Greeting struct{}

var ctx context.Context

func guestbookKey(ctx context.Context) *datastore.Key {
	return nil
}

func example() {
	// [START master_slave_data_definition_code]
	g := Greeting{ /* ... */ }
	key := datastore.NewIncompleteKey(ctx, "Greeting", nil)
	// [END master_slave_data_definition_code]

	// [START master_slave_query_code]
	q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
	// [END master_slave_query_code]

	_ = g
	_ = key
	_ = q
}

func example2() {
	// [START high_replication_data_definition_code]
	g := Greeting{ /* ... */ }
	key := datastore.NewIncompleteKey(ctx, "Greeting", guestbookKey(ctx))
	// [END high_replication_data_definition_code]

	// [START high_replication_query_code]
	q := datastore.NewQuery("Greeting").Ancestor(guestbookKey(ctx)).Order("-Date").Limit(10)
	// [END high_replication_query_code]

	_ = g
	_ = key
	_ = q
}
