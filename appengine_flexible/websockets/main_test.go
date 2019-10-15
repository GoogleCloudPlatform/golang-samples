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

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSocketHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(socketHandler))
	defer server.Close()
	dialer := websocket.Dialer{}

	conn, resp, err := dialer.Dial("ws://"+server.Listener.Addr().String()+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	want := http.StatusSwitchingProtocols
	if got := resp.StatusCode; got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}

	message := []byte("echo test")
	if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
		t.Fatal(err)
	}

	_, got, err := conn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, message) {
		t.Errorf("got %q, want %q", got, message)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(""))

	rr := httptest.NewRecorder()
	healthCheckHandler(rr, req)

	if got, want := rr.Body.String(), "ok"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
