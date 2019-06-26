// Copyright 2019 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHelloPubSubErrors(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "no_payload"},
		{name: "not_base64"},
	}
	for _, test := range tests {
		var payload *strings.Reader
		if test.name == "no_payload" {
			payload = strings.NewReader("");
		} else {
			not_encoded := "Gopher"
			jsonStr := fmt.Sprintf(`{"message":{"data":"%s","id":"test-123"}}`, not_encoded)
			payload = strings.NewReader(jsonStr)
		}
		req := httptest.NewRequest("GET", "/", payload)
		rr := httptest.NewRecorder()

		HelloPubSub(rr, req)

		if rr.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("HelloPubSub(%q) should get BadRequest response", test.name)
		}
	}
}

func TestHelloPubSub(t *testing.T) {
	tests := []struct {
		data string
		want string
	}{
		{want: "Hello World!\n"},
		{data: "Go", want: "Hello Go!\n"},
	}
	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		originalFlags := log.Flags()
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		payload := strings.NewReader("{}")
		if test.data != "" {
			encoded := base64.StdEncoding.EncodeToString([]byte(test.data))
			jsonStr := fmt.Sprintf(`{"message":{"data":"%s","id":"test-123"}}`, encoded)
			payload = strings.NewReader(jsonStr)
		}
		req := httptest.NewRequest("GET", "/", payload)
		rr := httptest.NewRecorder()

		HelloPubSub(rr, req)

		w.Close()
		log.SetOutput(os.Stderr)
		log.SetFlags(originalFlags)

		if rr.Result().StatusCode == http.StatusBadRequest {
			t.Errorf("HelloPubSub received invalid input (%q)", test.data)
		}

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("HelloPubSub(%q) = %q, want %q", test.data, got, test.want)
		}
	}
}
