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
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

var maxHeight int
var minBirthYear, maxBirthYear int

// [START gae_go_datastore_interface]
type Person struct {
	FirstName string
	LastName  string
	City      string
	BirthYear int
	Height    int
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// The Query type and its methods are used to construct a query.
	q := datastore.NewQuery("Person").
		Filter("LastName =", "Smith").
		Filter("Height <=", maxHeight).
		Order("-Height")

	// To retrieve the results,
	// you must execute the Query using its GetAll or Run methods.
	var people []Person
	if _, err := q.GetAll(ctx, &people); err != nil {
		// Handle error.
	}
	// ...
}

// [END gae_go_datastore_interface]

func example() {
	var lastSeenKey *datastore.Key

	// [START gae_go_datastore_key_filter]
	q := datastore.NewQuery("Person").Filter("__key__ >", lastSeenKey)
	// [END gae_go_datastore_key_filter]
	_ = q
}

func example2() {
	// [START gae_go_datastore_property_filter]
	q := datastore.NewQuery("Person").Filter("Height <=", maxHeight)
	// [END gae_go_datastore_property_filter]
	_ = q
}

func example3() {
	var ancestorKey *datastore.Key

	// [START gae_go_datastore_ancestor_filter]
	q := datastore.NewQuery("Person").Ancestor(ancestorKey)
	// [END gae_go_datastore_ancestor_filter]
	_ = q
}

func example4() {
	// [START gae_go_datastore_sort_order]
	// Order alphabetically by last name:
	q := datastore.NewQuery("Person").Order("LastName")

	// Order by height, tallest to shortest:
	q = datastore.NewQuery("Person").Order("-Height")
	// [END gae_go_datastore_sort_order]
	_ = q
}

func example5() {
	// [START gae_go_datastore_multiple_sort_orders]
	q := datastore.NewQuery("Person").Order("LastName").Order("-Height")
	// [END gae_go_datastore_multiple_sort_orders]
	_ = q
}

func example6() {
	type Photo struct {
		URL string
	}
	var ctx context.Context

	// [START gae_go_datastore_ancestor_query]
	// Create two Photo entities in the datastore with a Person as their ancestor.
	tomKey := datastore.NewKey(ctx, "Person", "Tom", 0, nil)

	wPhoto := Photo{URL: "http://example.com/some/path/to/wedding_photo.jpg"}
	wKey := datastore.NewKey(ctx, "Photo", "", 0, tomKey)
	_, err := datastore.Put(ctx, wKey, wPhoto)
	// check err

	bPhoto := Photo{URL: "http://example.com/some/path/to/baby_photo.jpg"}
	bKey := datastore.NewKey(ctx, "Photo", "", 0, tomKey)
	_, err = datastore.Put(ctx, bKey, bPhoto)
	// check err

	// Now fetch all Photos that have tomKey as an ancestor.
	// This will populate the photos slice with wPhoto and bPhoto.
	q := datastore.NewQuery("Photo").Ancestor(tomKey)
	var photos []Photo
	_, err = q.GetAll(ctx, &photos)
	// check err
	// do something with photos
	// [END gae_go_datastore_ancestor_query]
	_ = err
	_ = photos
}

func example7() {
	// [START gae_go_datastore_keys_only]
	q := datastore.NewQuery("Person").KeysOnly()
	// [END gae_go_datastore_keys_only]
	_ = q
}

func example8() {
	// [START gae_go_datastore_inequality_filters_one_property_valid_1]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Filter("BirthYear <=", maxBirthYear)
	// [END gae_go_datastore_inequality_filters_one_property_valid_1]
	_ = q
}

func example9() {
	// [START gae_go_datastore_inequality_filters_one_property_invalid]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Filter("Height <=", maxHeight) // ERROR
	// [END gae_go_datastore_inequality_filters_one_property_invalid]
	_ = q
}

func example10() {
	var targetLastName, targetCity string

	// [START gae_go_datastore_inequality_filters_one_property_valid_2]
	q := datastore.NewQuery("Person").
		Filter("LastName =", targetLastName).
		Filter("City =", targetCity).
		Filter("BirthYear >=", minBirthYear).
		Filter("BirthYear <=", maxBirthYear)
	// [END gae_go_datastore_inequality_filters_one_property_valid_2]
	_ = q
}

func example11() {
	// [START gae_go_datastore_inequality_filters_sort_orders_valid]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Order("BirthYear").
		Order("LastName")
	// [END gae_go_datastore_inequality_filters_sort_orders_valid]
	_ = q
}

func example12() {
	// [START gae_go_datastore_inequality_filters_sort_orders_invalid_1]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Order("LastName") // ERROR
	// [END gae_go_datastore_inequality_filters_sort_orders_invalid_1]
	_ = q
}

func example13() {
	// [START gae_go_datastore_inequality_filters_sort_orders_invalid_2]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Order("LastName").
		Order("BirthYear") // ERROR
	// [END gae_go_datastore_inequality_filters_sort_orders_invalid_2]
	_ = q
}

func example14() {
	// [START gae_go_datastore_surprising_behavior_1]
	q := datastore.NewQuery("Widget").
		Filter("x >", 1).
		Filter("x <", 2)
	// [END gae_go_datastore_surprising_behavior_1]
	_ = q
}

func example15() {
	// [START gae_go_datastore_surprising_behavior_2]
	q := datastore.NewQuery("Widget").
		Filter("x =", 1).
		Filter("x =", 2)
	// [END gae_go_datastore_surprising_behavior_2]
	_ = q
}

func doSomething(k *datastore.Key, p Person) {}

func example16() {
	var ctx context.Context
	// [START gae_go_datastore_retrieval]
	q := datastore.NewQuery("Person")
	t := q.Run(ctx)
	for {
		var p Person
		k, err := t.Next(&p)
		if err == datastore.Done {
			break // No further entities match the query.
		}
		if err != nil {
			log.Errorf(ctx, "fetching next Person: %v", err)
			break
		}
		// Do something with Person p and Key k
		doSomething(k, p)
	}
	// [END gae_go_datastore_retrieval]
}

func example17() {
	var ctx context.Context

	// [START gae_go_datastore_all_entities_retrieval]
	q := datastore.NewQuery("Person")
	var people []Person
	keys, err := q.GetAll(ctx, &people)
	if err != nil {
		log.Errorf(ctx, "fetching people: %v", err)
		return
	}
	for i, p := range people {
		k := keys[i]
		// Do something with Person p and Key k
		doSomething(k, p)
	}
	// [END gae_go_datastore_all_entities_retrieval]
}

func example18() {
	var ctx context.Context

	// [START gae_go_datastore_query_limit]
	q := datastore.NewQuery("Person").Order("-Height").Limit(5)
	var people []Person
	_, err := q.GetAll(ctx, &people)
	// check err

	for _, p := range people {
		log.Infof(ctx, "%s %s, %d inches tall", p.FirstName, p.LastName, p.Height)
	}
	// [END gae_go_datastore_query_limit]
	_ = err
}

func example19() {
	q := datastore.NewQuery("Person").Order("-Height").Limit(5).Offset(5)
	_ = q
}

func example20() {
	var ctx context.Context

	// [START gae_go_datastore_cursors]
	// Create a query for all Person entities.
	q := datastore.NewQuery("Person")

	// If the application stored a cursor during a previous request, use it.
	item, err := memcache.Get(ctx, "person_cursor")
	if err == nil {
		cursor, err := datastore.DecodeCursor(string(item.Value))
		if err == nil {
			q = q.Start(cursor)
		}
	}

	// Iterate over the results.
	t := q.Run(ctx)
	for {
		var p Person
		_, err := t.Next(&p)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "fetching next Person: %v", err)
			break
		}
		// Do something with the Person p
	}

	// Get updated cursor and store it for next time.
	if cursor, err := t.Cursor(); err == nil {
		memcache.Set(ctx, &memcache.Item{
			Key:   "person_cursor",
			Value: []byte(cursor.String()),
		})
	}
	// [END gae_go_datastore_cursors]
}

func example21() {
	var lastSeenKey *datastore.Key

	// [START gae_go_datastore_kindless_query]
	q := datastore.NewQuery("").Filter("__key__ >", lastSeenKey)
	// [END gae_go_datastore_kindless_query]
	_ = q
}

func example22() {
	var ancestorKey, lastSeenKey *datastore.Key

	// [START gae_go_datastore_kindless_ancestor_key_query]
	q := datastore.NewQuery("").Ancestor(ancestorKey).Filter("__key__ >", lastSeenKey)
	// [END gae_go_datastore_kindless_ancestor_key_query]
	_ = q
}

func example23() {
	type Photo struct {
		URL string
	}
	type Video Photo
	var ctx context.Context
	doSomething := func(x interface{}) {}

	// [START gae_go_datastore_kindless_ancestor_query]
	tomKey := datastore.NewKey(ctx, "Person", "Tom", 0, nil)

	weddingPhoto := &Photo{URL: "http://example.com/some/path/to/wedding_photo.jpg"}
	_, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Photo", tomKey), weddingPhoto)

	weddingVideo := &Video{URL: "http://example.com/some/path/to/wedding_video.avi"}
	_, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Video", tomKey), weddingVideo)

	// The following query returns both weddingPhoto and weddingVideo,
	// even though they are of different entity kinds.
	q := datastore.NewQuery("").Ancestor(tomKey)
	t := q.Run(ctx)
	for {
		var x interface{}
		_, err := t.Next(&x)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "fetching next Photo/Video: %v", err)
			break
		}
		// Do something (e.g. switch on types)
		doSomething(x)
	}
	// [END gae_go_datastore_kindless_ancestor_query]
	_ = err
}
