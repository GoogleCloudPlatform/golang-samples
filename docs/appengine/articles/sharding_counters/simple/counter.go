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

package sharded_counter

import (
	"context"
	"fmt"
	"math/rand"

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
