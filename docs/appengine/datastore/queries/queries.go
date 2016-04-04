// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

var maxHeight int
var minBirthYear, maxBirthYear int

// [START interface]
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

// [END interface]

func example() {
	var lastSeenKey *datastore.Key

	// [START key_filter_example]
	q := datastore.NewQuery("Person").Filter("__key__ >", lastSeenKey)
	// [END key_filter_example]
	_ = q
}

func example2() {
	// [START property_filter_example]
	q := datastore.NewQuery("Person").Filter("Height <=", maxHeight)
	// [END property_filter_example]
	_ = q
}

func example3() {
	var ancestorKey *datastore.Key

	// [START ancestor_filter_example]
	q := datastore.NewQuery("Person").Ancestor(ancestorKey)
	// [END ancestor_filter_example]
	_ = q
}

func example4() {
	// [START sort_order_example]
	// Order alphabetically by last name:
	q := datastore.NewQuery("Person").Order("LastName")

	// Order by height, tallest to shortest:
	q = datastore.NewQuery("Person").Order("-Height")
	// [END sort_order_example]
	_ = q
}

func example5() {
	// [START multiple_sort_orders_example]
	q := datastore.NewQuery("Person").Order("LastName").Order("-Height")
	// [END multiple_sort_orders_example]
	_ = q
}

func example6() {
	type Photo struct {
		URL string
	}
	var ctx context.Context

	// [START ancestor_query_example]
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
	// [END ancestor_query_example]
	_ = err
	_ = photos
}

func example7() {
	// [START keys_only_example]
	q := datastore.NewQuery("Person").KeysOnly()
	// [END keys_only_example]
	_ = q
}

func example8() {
	// [START inequality_filters_one_property_valid_example_1]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Filter("BirthYear <=", maxBirthYear)
	// [END inequality_filters_one_property_valid_example_1]
	_ = q
}

func example9() {
	// [START inequality_filters_one_property_invalid_example]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Filter("Height <=", maxHeight) // ERROR
	// [END inequality_filters_one_property_invalid_example]
	_ = q
}

func example10() {
	var targetLastName, targetCity string

	// [START inequality_filters_one_property_valid_example_2]
	q := datastore.NewQuery("Person").
		Filter("LastName =", targetLastName).
		Filter("City =", targetCity).
		Filter("BirthYear >=", minBirthYear).
		Filter("BirthYear <=", maxBirthYear)
	// [END inequality_filters_one_property_valid_example_2]
	_ = q
}

func example11() {
	// [START inequality_filters_sort_orders_valid_example]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Order("BirthYear").
		Order("LastName")
	// [END inequality_filters_sort_orders_valid_example]
	_ = q
}

func example12() {
	// [START inequality_filters_sort_orders_invalid_example_1]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Order("LastName") // ERROR
	// [END inequality_filters_sort_orders_invalid_example_1]
	_ = q
}

func example13() {
	// [START inequality_filters_sort_orders_invalid_example_2]
	q := datastore.NewQuery("Person").
		Filter("BirthYear >=", minBirthYear).
		Order("LastName").
		Order("BirthYear") // ERROR
	// [END inequality_filters_sort_orders_invalid_example_2]
	_ = q
}

func example14() {
	// [START surprising_behavior_example_1]
	q := datastore.NewQuery("Widget").
		Filter("x >", 1).
		Filter("x <", 2)
	// [END surprising_behavior_example_1]
	_ = q
}

func example15() {
	// [START surprising_behavior_example_2]
	q := datastore.NewQuery("Widget").
		Filter("x =", 1).
		Filter("x =", 2)
	// [END surprising_behavior_example_2]
	_ = q
}

func doSomething(k *datastore.Key, p Person) {}

func example16() {
	var ctx context.Context
	// [START retrieval_example]
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
	// [END retrieval_example]
}

func example17() {
	var ctx context.Context

	// [START all_entities_retrieval_example]
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
	// [END all_entities_retrieval_example]
}

func example18() {
	var ctx context.Context

	// [START query_limit_example]
	q := datastore.NewQuery("Person").Order("-Height").Limit(5)
	var people []Person
	_, err := q.GetAll(ctx, &people)
	// check err

	for _, p := range people {
		log.Infof(ctx, "%s %s, %d inches tall", p.FirstName, p.LastName, p.Height)
	}
	// [END query_limit_example]
	_ = err
}

func example19() {
	// [START query_offset_example]
	q := datastore.NewQuery("Person").Order("-Height").Limit(5).Offset(5)
	// [END query_offset_example]
	_ = q
}

func example20() {
	var ctx context.Context

	// [START cursors]
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
	// [END cursors]
}

func example21() {
	var lastSeenKey *datastore.Key

	// [START kindless_query_example]
	q := datastore.NewQuery("").Filter("__key__ >", lastSeenKey)
	// [END kindless_query_example]
	_ = q
}

func example22() {
	var ancestorKey, lastSeenKey *datastore.Key

	// [START kindless_ancestor_key_query_example]
	q := datastore.NewQuery("").Ancestor(ancestorKey).Filter("__key__ >", lastSeenKey)
	// [END kindless_ancestor_key_query_example]
	_ = q
}

func example23() {
	type Photo struct {
		URL string
	}
	type Video Photo
	var ctx context.Context
	doSomething := func(x interface{}) {}

	// [START kindless_ancestor_query_example]
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
	// [END kindless_ancestor_query_example]
	_ = err
}
