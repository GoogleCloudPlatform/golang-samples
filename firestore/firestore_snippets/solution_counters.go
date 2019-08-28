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
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// [START fs_counter_classes]
// counters/${ID}
type Counter struct {
	numShards int
}

// counters/${ID}/shards/${NUM}
type Shard map[string]int

// [END fs_counter_classes]

// [START fs_create_counter]
func (d *Counter) сreateCounter(ctx context.Context, docRef *firestore.DocumentRef) []error {
	// Initialize the counter document, then initialize each shard.
	errsList := make([]error, 0, d.numShards)
	colRef := docRef.Collection("shards")

	// Initialize each shard with count=0
	for num := 0; num < d.numShards; num++ {
		shardRef := colRef.Doc(strconv.Itoa(num))
		shard := Shard{"count": 0}

		_, err := shardRef.Set(ctx, shard)
		errsList = append(errsList, err)
	}
	return errsList
}

// [END fs_create_counter]

// [START fs_increment_counter]
func (d *Counter) incrementCounter(ctx context.Context, docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	docID := strconv.Itoa(rng.Intn(d.numShards))

	shardRef := docRef.Collection("shards").Doc(docID)
	return shardRef.Update(
		ctx, []firestore.Update{
			{Path: "count", Value: firestore.Increment(1)},
		})
}

// [END fs_increment_counter]

// [START fs_get_count]
func (d *Counter) getCount(ctx context.Context, docRef *firestore.DocumentRef) (total int64, err error) {
	// Sum the count of each shard in the subcollection
	shards := docRef.Collection("shards").Documents(ctx)
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return total, err
		}

		vTotal := doc.Data()["count"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return -1, fmt.Errorf("firestore: invalid dataType %T, want int64", vTotal)
		}
		total += shardCount
	}
	return total, nil
}

// [END fs_get_count]
