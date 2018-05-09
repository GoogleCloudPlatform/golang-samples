// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START memorystore_main_go]

// Command redis is a basic app that connects to a managed Redis instance.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	conn := redisPool.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("INCR", "visits"))
	if err != nil {
		http.Error(w, "Error incrementing visitor counter", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Visitor number: %d", counter)
}

func main() {
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	const maxConnections = 10
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}, maxConnections)

	http.HandleFunc("/", incrementHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// [END memorystore_main_go]
