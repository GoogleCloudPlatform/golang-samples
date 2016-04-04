// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

// [START overview]
import (
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

type Employee struct {
	FirstName          string
	LastName           string
	HireDate           time.Time
	AttendedHRTraining bool
}

func f(ctx context.Context) {
	// ...
	employee := &Employee{
		FirstName: "Antonio",
		LastName:  "Salieri",
		HireDate:  time.Now(),
	}
	employee.AttendedHRTraining = true

	key := datastore.NewIncompleteKey(ctx, "Employee", nil)
	if _, err := datastore.Put(ctx, key, employee); err != nil {
		// Handle err
	}
	// ...
}

// [END overview]
