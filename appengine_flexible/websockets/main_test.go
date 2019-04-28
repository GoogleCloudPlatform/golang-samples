package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSocketHandler(t *testing.T) {
	var err error

	server := httptest.NewServer(http.HandlerFunc(socketHandler))
	defer server.Close()
	dialer := websocket.Dialer{}

	conn, resp, err := dialer.Dial("ws://"+server.Listener.Addr().String()+"/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	got := resp.StatusCode
	want := http.StatusSwitchingProtocols
	if got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("echo test"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", strings.NewReader(""))

	rr := httptest.NewRecorder()
	healthCheckHandler(rr, req)

	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	want := "ok"
	if got := string(out); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
