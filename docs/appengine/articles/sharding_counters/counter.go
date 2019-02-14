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
	"google.golang.org/appengine/memcache"
)

type counterConfig struct {
	Shards int
}

type shard struct {
	Name  string
	Count int
}

const (
	defaultShards = 20
	configKind    = "GeneralCounterShardConfig"
	shardKind     = "GeneralCounterShard"
)

func memcacheKey(name string) string {
	return shardKind + ":" + name
}

// Count retrieves the value of the named counter.
func Count(ctx context.Context, name string) (int, error) {
	total := 0
	mkey := memcacheKey(name)
	if _, err := memcache.JSON.Get(ctx, mkey, &total); err == nil {
		return total, nil
	}
	q := datastore.NewQuery(shardKind).Filter("Name =", name)
	for t := q.Run(ctx); ; {
		var s shard
		_, err := t.Next(&s)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return total, err
		}
		total += s.Count
	}
	memcache.JSON.Set(ctx, &memcache.Item{
		Key:        mkey,
		Object:     &total,
		Expiration: 60,
	})
	return total, nil
}

// Increment increments the named counter.
func Increment(ctx context.Context, name string) error {
	// Get counter config.
	var cfg counterConfig
	ckey := datastore.NewKey(ctx, configKind, name, 0, nil)
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		err := datastore.Get(ctx, ckey, &cfg)
		if err == datastore.ErrNoSuchEntity {
			cfg.Shards = defaultShards
			_, err = datastore.Put(ctx, ckey, &cfg)
		}
		return err
	}, nil)
	if err != nil {
		return err
	}
	var s shard
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		shardName := fmt.Sprintf("%s-shard%d", name, rand.Intn(cfg.Shards))
		key := datastore.NewKey(ctx, shardKind, shardName, 0, nil)
		err := datastore.Get(ctx, key, &s)
		// A missing entity and a present entity will both work.
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		s.Name = name
		s.Count++
		_, err = datastore.Put(ctx, key, &s)
		return err
	}, nil)
	if err != nil {
		return err
	}
	memcache.IncrementExisting(ctx, memcacheKey(name), 1)
	return nil
}

// IncreaseShards increases the number of shards for the named counter to n.
// It will never decrease the number of shards.
func IncreaseShards(ctx context.Context, name string, n int) error {
	ckey := datastore.NewKey(ctx, configKind, name, 0, nil)
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		var cfg counterConfig
		mod := false
		err := datastore.Get(ctx, ckey, &cfg)
		if err == datastore.ErrNoSuchEntity {
			cfg.Shards = defaultShards
			mod = true
		} else if err != nil {
			return err
		}
		if cfg.Shards < n {
			cfg.Shards = n
			mod = true
		}
		if mod {
			_, err = datastore.Put(ctx, ckey, &cfg)
		}
		return err
	}, nil)
}
