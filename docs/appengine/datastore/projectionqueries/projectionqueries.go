// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

var ctx context.Context

type EventLog struct {
	Title, ReadPath, DateWritten string
}

func example() {
	// [START using_1]
	q := datastore.NewQuery("People").Project("FirstName", "LastName")
	// [END using_1]
	_ = q
}

func example2() {
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
}

func example3() {
	// [START grouping]
	q := datastore.NewQuery("Person").
		Project("LastName", "Height").Distinct().
		Filter("Height >", 20).
		Order("-Height").Order("LastName")
	// [END grouping]
	_ = q

	type Foo struct {
		A []int
		B []string
	}

	// [START projections_and_multiple_valued_properties_1]
	entity := Foo{A: []int{1, 1, 2, 3}, B: []string{"x", "y", "x"}}
	// [END projections_and_multiple_valued_properties_1]
	_ = entity
}

func example4() {
	// [START projections_and_multiple_valued_properties_2]
	q := datastore.NewQuery("Foo").Project("A", "B").Filter("A <", 3)
	// [END projections_and_multiple_valued_properties_2]
	_ = q
}
