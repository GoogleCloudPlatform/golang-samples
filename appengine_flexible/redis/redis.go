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

// [START gae_flex_golang_redis]

// Sample redis demonstrates use of a redis client from App Engine flexible environment.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"

	"google.golang.org/appengine"
)

var redisPool *redis.Pool

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisAddr)
			if redisPassword == "" {
				return conn, err
			}
			if err != nil {
				return nil, err
			}
			if _, err := conn.Do("AUTH", redisPassword); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		// TODO: Tune other settings, like IdleTimeout, MaxActive, MaxIdle, TestOnBorrow.
	}

	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	redisConn := redisPool.Get()
	defer redisConn.Close()

	count, err := redisConn.Do("INCR", "count")
	if err != nil {
		msg := fmt.Sprintf("Could not increment count: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Count: %d", count)
}

// [END gae_flex_golang_redis]
