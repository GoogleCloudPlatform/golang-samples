// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sharded_counter

import (
	"fmt"
	"math/rand"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

type simpleCounterShard struct {
	Count int
}

const (
	numShards = 20
	shardKind = "SimpleCounterShard"
)

// Count retrieves the value of the counter.
func Count(ctx context.Context) (int, error) {
	total := 0
	q := datastore.NewQuery(shardKind)
	for t := q.Run(ctx); ; {
		var s simpleCounterShard
		_, err := t.Next(&s)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return total, err
		}
		total += s.Count
	}
	return total, nil
}

// Increment increments the counter.
func Increment(ctx context.Context) error {
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		shardName := fmt.Sprintf("shard%d", rand.Intn(numShards))
		key := datastore.NewKey(ctx, shardKind, shardName, 0, nil)
		var s simpleCounterShard
		err := datastore.Get(ctx, key, &s)
		// A missing entity and a present entity will both work.
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		s.Count++
		_, err = datastore.Put(ctx, key, &s)
		return err
	}, nil)
}
