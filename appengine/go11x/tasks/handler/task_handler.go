// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START cloud_tasks_appengine_quickstart]

// Sample task_handler is an App Engine app demonstrating Cloud Tasks handling.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// Allow confirmation the task handling service is running.
	http.HandleFunc("/", indexHandler)

	// Handle all tasks.
	http.HandleFunc("/task_handler", taskHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

// taskHandler processes task requests.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	t, ok := r.Header["X-Appengine-Taskname"]
	if !ok || len(t[0]) == 0 {
		// You may use the presence of the X-Appengine-Taskname header to validate
		// the request comes from Cloud Tasks.
		log.Println("Invalid Task: No X-Appengine-Taskname request header found")
		http.Error(w, "Bad Request - Invalid Task", http.StatusBadRequest)
		return
	}
	taskName := t[0]

	// Pull useful headers from Task request.
	q, ok := r.Header["X-Appengine-Queuename"]
	queueName := ""
	if ok {
		queueName = q[0]
	}

	// Extract the request body for further task details.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll: %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	// Log & output details of the task.
	output := fmt.Sprintf("Completed task: task queue(%s), task name(%s), payload(%s)",
		queueName,
		taskName,
		string(body),
	)
	log.Println(output)

	// Set a non-2xx status code to indicate a failure in task processing that should be retried.
	// For example, http.Error(w, "Internal Server Error: Task Processing", http.StatusInternalServerError)
	fmt.Fprintln(w, output)
}

// [END cloud_tasks_appengine_quickstart]
