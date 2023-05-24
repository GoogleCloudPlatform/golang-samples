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

// A simple command-line task list manager to demonstrate using the
// cloud.google.com/go/datastore package.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

func main() {
	projID := os.Getenv("DATASTORE_PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
	}
	client, err := createClient(projID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}
	defer client.Close()

	// Print welcome message.
	fmt.Println("Cloud Datastore Task List")
	fmt.Println()
	usage()

	// Read commands from stdin.
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		cmd, args, n := parseCmd(scanner.Text())
		switch cmd {
		case "new":
			if args == "" {
				log.Printf("Missing description in %q command", cmd)
				usage()
				break
			}
			key, err := AddTask(projID, args)
			if err != nil {
				log.Printf("Failed to create task: %v", err)
				break
			}
			fmt.Printf("Created new task with ID %d\n", key.ID)

		case "done":
			if n == 0 {
				log.Printf("Missing numerical task ID in %q command", cmd)
				usage()
				break
			}
			if err := MarkDone(projID, n); err != nil {
				log.Printf("Failed to mark task done: %v", err)
				break
			}
			fmt.Printf("Task %d marked done\n", n)

		case "list":
			tasks, err := ListTasks(projID)
			if err != nil {
				log.Printf("Failed to fetch task list: %v", err)
				break
			}
			PrintTasks(os.Stdout, tasks)

		case "delete":
			if n == 0 {
				log.Printf("Missing numerical task ID in %q command", cmd)
				usage()
				break
			}
			if err := DeleteTask(projID, n); err != nil {
				log.Printf("Failed to delete task: %v", err)
				break
			}
			fmt.Printf("Task %d deleted\n", n)

		default:
			log.Printf("Unknown command %q", cmd)
			usage()
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed reading stdin: %v", err)
	}
}

// PrintTasks prints the tasks to the given writer.
func PrintTasks(w io.Writer, tasks []*Task) {
	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 8, 8, 1, ' ', 0) // Min cell size of 8.
	fmt.Fprintf(tw, "ID\tDescription\tStatus\n")
	for _, t := range tasks {
		if t.Done {
			fmt.Fprintf(tw, "%d\t%s\tdone\n", t.id, t.Desc)
		} else {
			fmt.Fprintf(tw, "%d\t%s\tcreated %v\n", t.id, t.Desc, t.Created)
		}
	}
	tw.Flush()
}

func usage() {
	fmt.Print(`Usage:

  new <description>  Adds a task with a description <description>
  done <task-id>     Marks a task as done
  list               Lists all tasks by creation time
  delete <task-id>   Deletes a task
`)
}

// parseCmd splits a line into a command and optional extra args.
// n will be set if the extra args can be parsed as an int64.
func parseCmd(line string) (cmd, args string, n int64) {
	if f := strings.Fields(line); len(f) > 0 {
		cmd = f[0]
		args = strings.Join(f[1:], " ")
	}
	if i, err := strconv.ParseInt(args, 10, 64); err == nil {
		n = i
	}
	return cmd, args, n
}
