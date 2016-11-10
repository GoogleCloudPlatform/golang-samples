// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample redis demonstrates use of a redis client from App Engine flexible environment.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"

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
