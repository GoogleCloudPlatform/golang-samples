// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/fluent/fluent-logger-golang/fluent"
)

var logger *fluent.Fluent

func init() {
	// Initializes a Fluentd logger.
	var err error
	logger, err = fluent.New(
		fluent.Config{FluentPort: 24224, FluentHost: "localhost"})
	if err != nil {
		panic(err)
	}

	// Registers a handler.
	http.HandleFunc("/demo", demoHandler)
}

func report(stackTrace string, r *http.Request) {
	payload := map[string]interface{}{
		"serviceContext": map[string]interface{}{
			"service": "myapp",
		},
		"message": stackTrace,
		"context": map[string]interface{}{
			"httpRequest": map[string]interface{}{
				"method":    r.Method,
				"url":       r.URL.String(),
				"userAgent": r.UserAgent(),
				"referrer":  r.Referer(),
				"remoteIp":  r.RemoteAddr,
			},
		},
	}
	err := logger.Post("myapp.errors", payload)
	if err != nil {
		log.Fatal(err)
	}
}

// Handler for the incomming requests.
func demoHandler(w http.ResponseWriter, r *http.Request) {
	// How to handle an error.
	defer func() {
		if e := recover(); e != nil {
			stack := make([]byte, 1<<16)
			stackSize := runtime.Stack(stack, true)
			report(string(stack[:stackSize]), r)
		}
	}()

	// Error is generated here.
	x := 0
	log.Println(100500 / x)
}

// Http server starts serving requests.
func main() {
	http.ListenAndServe(":8080", nil)
}
