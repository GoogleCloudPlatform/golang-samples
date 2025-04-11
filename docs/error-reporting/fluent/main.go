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

// TODO(cbro): find a better sample - panic, defer, recover should not be present.

// Sample fluent demonstrates integration of fluent and Cloud Error reporting.
package main

import (
	"log"
	"net/http"
	"os"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
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
