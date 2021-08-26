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

//go:build go1.8
// +build go1.8

// Unit tests for example app

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandleCheckMessages(t *testing.T) {
	os.Setenv("MESSAGE_SERVICE", "mock")
	r := httptest.NewRequest("GET", "http://messages?user=Friend2", nil)
	w := httptest.NewRecorder()
	handleCheckMessages(w, r)
	resp := w.Result()
	expected := http.StatusOK
	result := resp.StatusCode
	if result != expected {
		t.Errorf("TestHandleCheckMessages: Expected: %d, got %d\n", expected,
			result)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	expect := "message(s)"
	if !strings.Contains(content, expect) {
		t.Errorf("TestHandleCheckMessages: Expect to contain: %s, got, %s\n",
			expect, content)
	}
}

func TestHandleDefault(t *testing.T) {
	os.Setenv("MESSAGE_SERVICE", "mock")
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	handleDefault(w, r)
	resp := w.Result()
	expected := http.StatusOK
	result := resp.StatusCode
	if result != expected {
		t.Errorf("TestHandleDefault: Expected: %d, got %d\n", expected, result)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	expect := "two options"
	if !strings.Contains(content, expect) {
		t.Errorf("TestHandleDefault: Expect to contain: %s, got, %s\n",
			expect, content)
	}
}

func TestHandleSend(t *testing.T) {
	os.Setenv("MESSAGE_SERVICE", "mock")
	r := httptest.NewRequest("GET",
		"http://send?user=Friend1&friend=Friend2&text=We+miss+you!", nil)
	w := httptest.NewRecorder()
	handleSend(w, r)
	resp := w.Result()
	expected := http.StatusOK
	result := resp.StatusCode
	if result != expected {
		t.Errorf("TestHandleSend: Expected: %d, got %d\n", expected,
			result)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	expect := "Message sent"
	if !strings.Contains(content, expect) {
		t.Errorf("TestHandleSend: Expect to contain: %s, got, %s\n",
			expect, content)
	}
}
