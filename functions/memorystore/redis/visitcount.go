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

// [START functions_memorystore_redis]

// Package visitcount provides a Cloud Function that connects
// to a managed Redis instance.
package visitcount

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func init() {
	// Register the HTTP handler with the Functions Framework
	functions.HTTP("VisitCount", visitCount)
}

// initializeRedis initializes and returns a connection pool
func initializeRedis() (*redis.Pool, error) {
	redisHost := os.Getenv("REDISHOST")
	if redisHost == "" {
		return nil, errors.New("REDISHOST must be set")
	}
	redisPort := os.Getenv("REDISPORT")
	if redisPort == "" {
		return nil, errors.New("REDISPORT must be set")
	}
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	const maxConnections = 10
	return &redis.Pool{
		MaxIdle: maxConnections,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, fmt.Errorf("redis.Dial: %w", err)
			}
			return c, err
		},
	}, nil
}

// visitCount increments the visit count on the Redis instance
// and prints the current count in the HTTP response.
func visitCount(w http.ResponseWriter, r *http.Request) {
	// Initialize connection pool on first invocation
	if redisPool == nil {
		// Pre-declare err to avoid shadowing redisPool
		var err error
		redisPool, err = initializeRedis()
		if err != nil {
			log.Printf("initializeRedis: %v", err)
			http.Error(w, "Error initializing connection pool", http.StatusInternalServerError)
			return
		}
	}

	conn := redisPool.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("INCR", "visits"))
	if err != nil {
		log.Printf("redis.Int: %v", err)
		http.Error(w, "Error incrementing visit count", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Visit count: %d", counter)
}

// [END functions_memorystore_redis]
