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

package visitcount

import (
  "fmt"
  "net/http"
  "os"

  "github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

// Initialize the connection pool on instance startup
func init() {
  redisHost := os.Getenv("REDISHOST")
  redisPort := os.Getenv("REDISPORT")
  redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

  const maxConnections = 10
  redisPool = redis.NewPool(func() (redis.Conn, error) {
    return redis.Dial("tcp", redisAddr)
  }, maxConnections)
}

func VisitCount(w http.ResponseWriter, r *http.Request) {
  conn := redisPool.Get()
  defer conn.Close()

  counter, err := redis.Int(conn.Do("INCR", "visits"))
  if err != nil {
    http.Error(w, "Error incrementing visit count", http.StatusInternalServerError)
    return
  }
  fmt.Fprintf(w, "Visit count: %d", counter)
}

// [END functions_memorystore_redis]
