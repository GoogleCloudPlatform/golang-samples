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
	"google.golang.org/appengine/log"
)

var ctx context.Context

type EventLog struct {
	Title, ReadPath, DateWritten string
}

func example() {
	// [START gae_datastore_using_1]
	// [START using_1]
	q := datastore.NewQuery("People").Project("FirstName", "LastName")
	// [END using_1]
	// [END gae_datastore_using_1]
	_ = q
}

func example2() {
	// [START gae_datastore_using_2]
	// [START using_2]
	q := datastore.NewQuery("EventLog").
		Project("Title", "ReadPath", "DateWritten").
		Order("DateWritten")
	t := q.Run(ctx)
	for {
		var l EventLog
		_, err := t.Next(&l)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "Running query: %v", err)
			break
		}
		log.Infof(ctx, "Log record: %v, %v, %v", l.Title, l.ReadPath, l.DateWritten)
	}
	// [END using_2]
	// [END gae_datastore_using_2]
}

func example3() {
	// [START gae_datastore_grouping]
	// [START grouping]
	q := datastore.NewQuery("Person").
		Project("LastName", "Height").Distinct().
		Filter("Height >", 20).
		Order("-Height").Order("LastName")
	// [END grouping]
	// [END gae_datastore_grouping]
	_ = q

	type Foo struct {
		A []int
		B []string
	}

	// [START gae_datastore_projections_and_multiple_valued_properties_1]
	// [START projections_and_multiple_valued_properties_1]
	entity := Foo{A: []int{1, 1, 2, 3}, B: []string{"x", "y", "x"}}
	// [END projections_and_multiple_valued_properties_1]
	// [END gae_datastore_projections_and_multiple_valued_properties_1]
	_ = entity
}

func example4() {
	// [START gae_datastore_projections_and_multiple_valued_properties_2]
	// [START projections_and_multiple_valued_properties_2]
	q := datastore.NewQuery("Foo").Project("A", "B").Filter("A <", 3)
	// [END projections_and_multiple_valued_properties_2]
	// [END gae_datastore_projections_and_multiple_valued_properties_2]
	_ = q
}
