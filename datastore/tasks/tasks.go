// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START all]

// A simple command-line task list manager to demonstrate using the
// cloud.google.com/go//datastore package.
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
	"time"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"
)

func main() {
	projID := os.Getenv("DATASTORE_PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
	}
	// [START build_service]
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projID)
	// [END build_service]
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

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
			key, err := AddTask(ctx, client, args)
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
			if err := MarkDone(ctx, client, n); err != nil {
				log.Printf("Failed to mark task done: %v", err)
				break
			}
			fmt.Printf("Task %d marked done\n", n)

		case "list":
			tasks, err := ListTasks(ctx, client)
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
			if err := DeleteTask(ctx, client, n); err != nil {
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

// [START add_entity]
// Task is the model used to store tasks in the datastore.
type Task struct {
	Desc    string    `datastore:"description"`
	Created time.Time `datastore:"created"`
	Done    bool      `datastore:"done"`
	id      int64     // The integer ID used in the datastore.
}

// AddTask adds a task with the given description to the datastore,
// returning the key of the newly created entity.
func AddTask(ctx context.Context, client *datastore.Client, desc string) (*datastore.Key, error) {
	task := &Task{
		Desc:    desc,
		Created: time.Now(),
	}
	key := datastore.IncompleteKey("Task", nil)
	return client.Put(ctx, key, task)
}

// [END add_entity]

// [START update_entity]
// MarkDone marks the task done with the given ID.
func MarkDone(ctx context.Context, client *datastore.Client, taskID int64) error {
	// Create a key using the given integer ID.
	key := datastore.IDKey("Task", taskID, nil)

	// In a transaction load each task, set done to true and store.
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var task Task
		if err := tx.Get(key, &task); err != nil {
			return err
		}
		task.Done = true
		_, err := tx.Put(key, &task)
		return err
	})
	return err
}

// [END update_entity]

// [START retrieve_entities]
// ListTasks returns all the tasks in ascending order of creation time.
func ListTasks(ctx context.Context, client *datastore.Client) ([]*Task, error) {
	var tasks []*Task

	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery("Task").Order("created")
	keys, err := client.GetAll(ctx, query, &tasks)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Task from the corresponding key.
	for i, key := range keys {
		tasks[i].id = key.ID
	}

	return tasks, nil
}

// [END retrieve_entities]

// [START delete_entity]
// DeleteTask deletes the task with the given ID.
func DeleteTask(ctx context.Context, client *datastore.Client, taskID int64) error {
	return client.Delete(ctx, datastore.IDKey("Task", taskID, nil))
}

// [END delete_entity]

// [START format_results]
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

// [END format_results]

func usage() {
	fmt.Println(`Usage:

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

// [END all]
