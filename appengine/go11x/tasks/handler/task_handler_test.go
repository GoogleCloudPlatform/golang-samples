// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	tests := []struct {
		route  string
		status int
		body   string
	}{
		{
			route:  "/",
			status: http.StatusOK,
			body:   "Hello, World!",
		},
		{
			route:  "/404",
			status: http.StatusNotFound,
			body:   "404 page not found\n",
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.route, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(indexHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.status {
			t.Errorf("unexpected status: got (%v) want (%v)", status, test.status)
		}

		if rr.Body.String() != test.body {
			t.Errorf("unexpected body: got (%v) want (%v)", rr.Body.String(), test.body)
		}
	}
}

func TestTaskHandler(t *testing.T) {
	tests := []struct {
		name     string
		taskName string
		message  string
		body     string
		status   int
	}{
		{
			name:     "Invalid Task",
			taskName: "",
			message:  "",
			body:     "Bad Request - Invalid Task\n",
			status:   http.StatusBadRequest,
		},
		{
			name:     "Valid Task, No Message",
			taskName: "1234",
			message:  "",
			status:   http.StatusOK,
		},
		{
			name:     "Valid Task, Text Message",
			taskName: "1234",
			message:  "task details",
			status:   http.StatusOK,
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/test_handler", strings.NewReader(test.message))
		if err != nil {
			t.Fatal(err)
		}
		req.Header["X-Appengine-Taskname"] = []string{test.taskName}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(taskHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.status {
			t.Errorf("%s: unexpected status: got (%v) want (%v)", test.name, status, test.status)
		}

		// Allow test cases to override the body message.
		want := test.body
		if test.body == "" {
			want = fmt.Sprintf("Completed task: task queue(%s), task name(%s), payload(%s)\n", "", test.taskName, test.message)
		}

		// HTTP Body might have embedded NUL characters.
		if got := rr.Body.String(); got != want {
			t.Errorf("%s: unexpected body:\n\tgot (%s)\n\twant (%s)", test.name, got, want)
		}
	}
}
