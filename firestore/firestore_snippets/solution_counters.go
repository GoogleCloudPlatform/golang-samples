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

package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// [START fs_counter_classes]
// counters/${ID}
type Counter struct {
	numShards int
}

// counters/${ID}/shards/${NUM}
type Shard struct {
	Count int
}

// [END fs_counter_classes]

// [START fs_create_counter]
func (c *Counter) сreateCounter(ctx context.Context, docRef *firestore.DocumentRef) error {
	// Initialize the counter document, then initialize each shard.
	var err error
	colRef := docRef.Collection("shards")

	// Initialize each shard with count=0
	for num := 0; num < c.numShards; num++ {
		shard := Shard{0}

		_, err = colRef.Doc(strconv.Itoa(num)).Set(ctx, shard)
		if err != nil {
			return fmt.Errorf("Set: %v", err)
		}
	}
	return nil
}

// [END fs_create_counter]

// [START fs_increment_counter]
func (c *Counter) incrementCounter(ctx context.Context, docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	docID := strconv.Itoa(rand.Intn(c.numShards))

	shardRef := docRef.Collection("shards").Doc(docID)
	return shardRef.Update(ctx, []firestore.Update{
		{Path: "Count", Value: firestore.Increment(1)},
	})
}

// [END fs_increment_counter]

// [START fs_get_count]
func (c *Counter) getCount(ctx context.Context, docRef *firestore.DocumentRef) (int64, error) {
	// Sum the count of each shard in the subcollection
	var total int64
	shards := docRef.Collection("shards").Documents(ctx)
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("Next: %v", err)
		}

		vTotal := doc.Data()["Count"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return 0, fmt.Errorf("firestore: invalid dataType %T, want int64", vTotal)
		}
		total += shardCount
	}
	return total, nil
}

// [END fs_get_count]
