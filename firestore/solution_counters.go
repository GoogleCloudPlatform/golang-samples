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

package firestore

// [START firestore_solution_sharded_counter_custom_type]
import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Counter is a collection of documents (shards)
// to realize counter with high frequency.
type Counter struct {
	numShards int
}

// Shard is a single counter, which is used in a group
// of other shards within Counter.
type Shard struct {
	Count int
}

// [END firestore_solution_sharded_counter_custom_type]

// [START firestore_solution_sharded_counter_create]

// initCounter creates a given number of shards as
// subcollection of specified document.
func (c *Counter) initCounter(ctx context.Context, docRef *firestore.DocumentRef) error {
	colRef := docRef.Collection("shards")

	// Initialize each shard with count=0
	for num := 0; num < c.numShards; num++ {
		shard := Shard{0}

		if _, err := colRef.Doc(strconv.Itoa(num)).Set(ctx, shard); err != nil {
			return fmt.Errorf("Set: %w", err)
		}
	}
	return nil
}

// [END firestore_solution_sharded_counter_create]

// [START firestore_solution_sharded_counter_increment]

// incrementCounter increments a randomly picked shard.
func (c *Counter) incrementCounter(ctx context.Context, docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	docID := strconv.Itoa(rand.Intn(c.numShards))

	shardRef := docRef.Collection("shards").Doc(docID)
	return shardRef.Update(ctx, []firestore.Update{
		{Path: "Count", Value: firestore.Increment(1)},
	})
}

// [END firestore_solution_sharded_counter_increment]

// [START firestore_solution_sharded_counter_get]

// getCount returns a total count across all shards.
func (c *Counter) getCount(ctx context.Context, docRef *firestore.DocumentRef) (int64, error) {
	var total int64
	shards := docRef.Collection("shards").Documents(ctx)
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("Next: %w", err)
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

// [END firestore_solution_sharded_counter_get]
