// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
<<<<<<< HEAD
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
=======
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
>>>>>>> appengine/tasks: fix gofmt on helloworld_test.go

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.route, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(indexHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.status {
			t.Errorf(
				"unexpected status: got (%v) want (%v)",
				status,
				test.status,
			)
		}

		expected := test.body
		if rr.Body.String() != expected {
			t.Errorf(
				"unexpected body: got (%v) want (%v)",
				rr.Body.String(),
				test.body,
			)
		}
	}
}
