// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// TODO(cbro): find a better sample - panic, defer, recover should not be present.

// Sample fluent demonstrates integration of fluent and Cloud Error reporting.
package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/fluent/fluent-logger-golang/fluent"
)

var logger *fluent.Fluent

func main() {
	var err error
	logger, err = fluent.New(fluent.Config{
		FluentHost: "localhost",
		FluentPort: 24224,
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/demo", demoHandler)
	http.ListenAndServe(":8080", nil)
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
	if err := logger.Post("myapp.errors", payload); err != nil {
		log.Print(err)
	}
}

// Handler for the incoming requests.
func demoHandler(w http.ResponseWriter, r *http.Request) {
	// How to handle a panic.
	defer func() {
		if e := recover(); e != nil {
			stack := make([]byte, 1<<16)
			stackSize := runtime.Stack(stack, true)
			report(string(stack[:stackSize]), r)
		}
	}()

	// Panic is triggered.
	x := 0
	log.Println(100500 / x)
}
