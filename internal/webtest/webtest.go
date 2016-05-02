// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package webtest provides helpers for testing web applications.
package webtest

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

// W holds the configuration for a web test.
type W struct {
	t    *testing.T
	host string

	Client *http.Client
}

// New creates a web test for a given a host tring (e.g. "localhost:8080")
func New(t *testing.T, host string) *W {
	return &W{
		t:      t,
		host:   host,
		Client: http.DefaultClient,
	}
}

// WaitForNet waits for the host to come live.
// After a 30s timeout, it will call t.Fatal
func (w *W) WaitForNet() {
	const retryDelay = 100 * time.Millisecond
	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", w.host)
		if err != nil {
			time.Sleep(retryDelay)
			continue
		}
		conn.Close()
		return
	}

	w.t.Fatalf("Timed out wating for net %s", w.host)
}

// GetBody performs a GET request to a given path.
func (w *W) GetBody(path string) (body string, resp *http.Response, err error) {
	resp, err = w.Get(path)
	if err != nil {
		return "", resp, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp, err
	}
	return string(b), resp, err
}

// Get performs a GET request to a given path.
func (w *W) Get(path string) (*http.Response, error) {
	return w.Client.Get("http://" + w.host + path)
}

// Post performs a POST request to a given path.
func (w *W) Post(path, bodyType string, body io.Reader) (*http.Response, error) {
	return w.Client.Post("http://"+w.host+path, bodyType, body)
}

// PostForm performs a POST request to a given path.
func (w *W) PostForm(path string, v url.Values) (*http.Response, error) {
	return w.Client.PostForm("http://"+w.host+path, v)
}

// NewRequest constructs a http.Request for the web tests's host.
func (w *W) NewRequest(method, path string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, "http://"+w.host+path, body)
	if err != nil {
		w.t.Fatal(err)
	}
	return r
}
