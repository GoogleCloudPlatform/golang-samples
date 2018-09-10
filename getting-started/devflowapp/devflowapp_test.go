// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
