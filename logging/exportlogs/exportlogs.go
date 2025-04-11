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

// Sample exportlogs lists, creates, updates, and deletes log sinks.
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		usage("Missing command.")
	}
	if len(os.Args) != 3 {
		usage("")
	}

	projID := os.Args[1]
	command := os.Args[2]

	switch command {
	case "list":
		log.Print("Listing log sinks.")
		sinks, err := listSinks(projID)
		if err != nil {
			log.Fatalf("Could not list log sinks: %v", err)
		}
		for _, sink := range sinks {
			fmt.Printf("Sink: %v\n", sink)
		}
	case "create":
		if _, err := createSink(projID); err != nil {
			log.Fatalf("Could not create sink: %v", err)
		}
	case "update":
		if _, err := updateSink(projID); err != nil {
			log.Fatalf("Could not update sink: %v", err)
		}
	case "delete":
		if err := deleteSink(projID); err != nil {
			log.Fatalf("Could not delete sink: %v", err)
		}
	default:
		usage("Unknown command.")
	}
}

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, "usage: exportlogs <project-id> [list|create|update|delete]")
	os.Exit(2)
}
