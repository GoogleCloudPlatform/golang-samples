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

var err error
var ctx context.Context
var k1, k2, k3 *datastore.Key
var e1, e2, e3 interface{}

type T struct{}
type Address struct{}

func example() {
	// [START batch]
	// A batch put.
	_, err = datastore.PutMulti(ctx, []*datastore.Key{k1, k2, k3}, []interface{}{e1, e2, e3})

	// A batch get.
	var entities = make([]*T, 3)
	err = datastore.GetMulti(ctx, []*datastore.Key{k1, k2, k3}, entities)

	// A batch delete.
	err = datastore.DeleteMulti(ctx, []*datastore.Key{k1, k2, k3})
	// [END batch]
	_ = err
}

func example2() {
	// [START delete]
	key := datastore.NewKey(ctx, "Employee", "asalieri", 0, nil)
	err = datastore.Delete(ctx, key)
	// [END delete]
}

func example3() {
	// [START get_key]
	employeeKey := datastore.NewKey(ctx, "Employee", "asalieri", 0, nil)
	addressKey := datastore.NewKey(ctx, "Address", "", 1, employeeKey)
	var addr Address
	err = datastore.Get(ctx, addressKey, &addr)
	// [END get_key]
}

func example4() {
	// [START key_id]
	// Create a key such as Employee:8261.
	key := datastore.NewKey(ctx, "Employee", "", 0, nil)
	// This is equivalent:
	key = datastore.NewIncompleteKey(ctx, "Employee", nil)
	// [END key_id]
	_ = key
}

func example5() {
	// [START key_name]
	// Create a key with a key name "asalieri".
	key := datastore.NewKey(
		ctx,        // context.Context
		"Employee", // Kind
		"asalieri", // String ID; empty means no string ID
		0,          // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		nil,        // Parent Key; nil means no parent
	)
	// [END key_name]
	_ = key
}

func example6() {
	// [START parent]
	// Create Employee entity
	employee := &Employee{ /* ... */ }
	employeeKey, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Employee", nil), employee)

	// Use Employee as Address entity's parent
	// and save Address entity to datastore
	address := &Address{ /* ... */ }
	addressKey := datastore.NewIncompleteKey(ctx, "Address", employeeKey)
	_, err = datastore.Put(ctx, addressKey, address)
	// [END parent]
	_ = err
}

func example7() {
	// [START put_with_keyname]
	employee := &Employee{
		FirstName: "Antonio",
		LastName:  "Salieri",
		HireDate:  time.Now(),
	}
	employee.AttendedHRTraining = true
	key := datastore.NewKey(ctx, "Employee", "asalieri", 0, nil)
	_, err = datastore.Put(ctx, key, employee)
	// [END put_with_keyname]
}

func example8() {
	// [START put_without_keyname]
	employee := &Employee{
		FirstName: "Antonio",
		LastName:  "Salieri",
		HireDate:  time.Now(),
	}
	employee.AttendedHRTraining = true
	key := datastore.NewIncompleteKey(ctx, "Employee", nil)
	_, err = datastore.Put(ctx, key, employee)
	// [END put_without_keyname]
}
