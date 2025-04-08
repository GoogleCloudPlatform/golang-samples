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

// Sample simplelog writes some entries, lists them, then deletes the log.
package main

import (
	"fmt"
	"log"
	"os"
	"time"
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
	case "write":
		log.Print("Writing log entries.")
		structuredWrite(projID)

	case "read":
		log.Print("Fetching and printing log entries.")
		entries, err := getEntries(projID)
		if err != nil {
			log.Fatalf("Could not get entries: %v", err)
		}
		log.Printf("Found %d entries.", len(entries))
		for _, entry := range entries {
			fmt.Printf("Entry: %6s @%s: %v\n",
				entry.Severity,
				entry.Timestamp.Format(time.RFC3339),
				entry.Payload)
		}

	case "delete":
		log.Print("Deleting log.")
		if err := deleteLog(projID); err != nil {
			log.Fatalf("Could not delete log: %v", err)
		}

	default:
		usage("Unknown command.")
	}
}

func usage(msg string) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, "usage: simplelog <project-id> [write|read|delete]")
	os.Exit(2)
}
