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

// [START gae_datastore_overview]
// [START overview]
import (
	"context"
	"time"

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
// [END gae_datastore_overview]
