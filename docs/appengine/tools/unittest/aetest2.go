// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package newsletter

// [START utility_example_2]
import (
	"errors"
	"testing"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestMyFunction(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	req1, err := inst.NewRequest("GET", "/gophers", nil)
	if err != nil {
		t.Fatalf("Failed to create req1: %v", err)
	}
	c1 := appengine.NewContext(req1)

	req2, err := inst.NewRequest("GET", "/herons", nil)
	if err != nil {
		t.Fatalf("Failed to create req2: %v", err)
	}
	c2 := appengine.NewContext(req2)

	// Run code and tests with *http.Request req1 and req2,
	// and context.Context c1 and c2.
	// [START_EXCLUDE]
	check(t, c1)
	check(t, c2)
	// [END_EXCLUDE]
}

// [END utility_example_2]

// [START datastore_example_1]
func TestWithdrawLowBal(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	key := datastore.NewKey(ctx, "BankAccount", "", 1, nil)
	if _, err := datastore.Put(ctx, key, &BankAccount{100}); err != nil {
		t.Fatal(err)
	}

	err = withdraw(ctx, "myid", 128, 0)
	if err == nil || err.Error() != "insufficient funds" {
		t.Errorf("Error: %v; want insufficient funds error", err)
	}

	b := BankAccount{}
	if err := datastore.Get(ctx, key, &b); err != nil {
		t.Fatal(err)
	}
	if bal, want := b.Balance, 100; bal != want {
		t.Errorf("Balance %d, want %d", bal, want)
	}
}

// [END datastore_example_1]

type BankAccount struct {
	Balance int
}

func withdraw(ctx context.Context, foo string, bar, baz int) error {
	return errors.New("insufficient funds")
}
