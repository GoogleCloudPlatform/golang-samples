// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START using_transactions]
package counter

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func init() {
	http.HandleFunc("/", handler)
}

type Counter struct {
	Count int
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	key := datastore.NewKey(ctx, "Counter", "mycounter", 0, nil)
	count := new(Counter)
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		// Note: this function's argument ctx shadows the variable ctx
		//       from the surrounding function.
		err := datastore.Get(ctx, key, count)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		count.Count++
		_, err = datastore.Put(ctx, key, count)
		return err
	}, nil)
	if err != nil {
		log.Errorf(ctx, "Transaction failed: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	fmt.Fprintf(w, "Current count: %d", count.Count)
}

// [END using_transactions]

// [START uses_for_transactions_1]
func increment(ctx context.Context, key *datastore.Key) error {
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		count := new(Counter)
		if err := datastore.Get(ctx, key, count); err != nil {
			return err
		}
		count.Count++
		_, err := datastore.Put(ctx, key, count)
		return err
	}, nil)
}

// [END uses_for_transactions_1]

// [START uses_for_transactions_2]
type Account struct {
	Address string
	Phone   string
}

func GetOrUpdate(ctx context.Context, id, addr, phone string) error {
	key := datastore.NewKey(ctx, "Account", id, 0, nil)
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		acct := new(Account)
		err := datastore.Get(ctx, key, acct)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		acct.Address = addr
		acct.Phone = phone
		_, err = datastore.Put(ctx, key, acct)
		return err
	}, nil)
}

// [END uses_for_transactions_2]

func example() {
	var ctx context.Context
	// [START transactional_task_enqueuing]
	datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		t := &taskqueue.Task{Path: "/path/to/worker"}
		if _, err := taskqueue.Add(ctx, t, ""); err != nil {
			return err
		}
		// ...
		return nil
	}, nil)
	// [END transactional_task_enqueuing]
}
