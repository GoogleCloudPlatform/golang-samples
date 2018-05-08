// Copyright 2018 Google Inc.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START main_go]
// A basic app that connects to a managed Redis instance.
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"os"
)

var redisPool *redis.Pool

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	conn := redisPool.Get()
	defer conn.Close()
	counter, err := redis.Int(conn.Do("INCR", "visits"))
	if err != nil {
		http.Error(w, "Error incrementing visitor counter", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Visitor number : %d", counter)
}

func main() {
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)
	maxConnections := 10
	redisPool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddress)
	}, maxConnections)
	defer redisPool.Close()
	http.HandleFunc("/", incrementHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
// [END main_go]
