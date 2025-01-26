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
	"time"

	"google.golang.org/appengine/datastore"
)

// [START gae_datastore_unindexed_properties]
// [START unindexed_properties]
type Person struct {
	Name string
	Age  int `datastore:",noindex"`
}

// [END unindexed_properties]
// [END gae_datastore_unindexed_properties]

// [START gae_datastore_exploding_index_example_3]
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
// [END gae_datastore_exploding_index_example_3]
