// Copyright 2021 Google LLC
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

// Sample sigterm-handler shows how to gracefully shut down in response to a SIGTERM signal.
package main

// [START cloudrun_sigterm_handler]
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Create channel to listen for signals.
var signalChan chan (os.Signal) = make(chan os.Signal, 1)

func main() {
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: http.HandlerFunc(handler),
	}

	// SIGINT handles Ctrl+C locally.
	// SIGTERM handles Cloud Run termination signal.
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server.
	go func() {
		log.Printf("listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Receive output from signalChan.
	sig := <-signalChan
	log.Printf("%s signal caught", sig)

	// Timeout if waiting for connections to return idle.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add extra handling here to clean up resources, such as flushing logs and
	// closing any database or Redis connections.

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %+v", err)
	}
	log.Print("server exited")
}

// [END cloudrun_sigterm_handler]

func handler(w http.ResponseWriter, r *http.Request) {
	// The 'terminate' parameter is used by tests in sigterm_handler.e2e_test.go
	if r.URL.Query().Get("terminate") != "" {
		fmt.Fprint(w, "Goodbye World!\n")
		signalChan <- syscall.SIGTERM
		return
	}
	fmt.Fprint(w, "Hello World!\n")
}
