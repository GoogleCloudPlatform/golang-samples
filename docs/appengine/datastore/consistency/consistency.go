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
	"context"

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
