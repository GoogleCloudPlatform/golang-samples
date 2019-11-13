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
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
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
