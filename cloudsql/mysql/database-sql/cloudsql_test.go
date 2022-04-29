// Copyright 2020 Google LLC
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

package cloudsql

import (
	"bytes"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// dbConfig holds database connection information derived from the environment.
type dbConfig struct {
	// user is the user name
	user string
	// pass is the user password
	pass string
	// name is the database name
	name string
	// port is the port used with the TCP host
	port string
	// host is a TCP address
	host string
	// unixPath is the unix socket file path
	unixPath string
	// instConnName is the Cloud SQL instance connection name
	instConnName string
}

type connType int

const (
	useTCP connType = iota
	useUnix
	useConnector
)

// dbConfigFromEnv reads database configuration from the environment, selecting
// the provided connType in the configuration.
func dbConfigFromEnv(t *testing.T, ct connType) dbConfig {
	testEnv := func(name string) string {
		n := os.Getenv(name)
		if n == "" {
			t.Fatalf("failed to get env var %q", name)
		}
		return n
	}
	d := dbConfig{
		user:         testEnv("MYSQL_USER"),
		pass:         testEnv("MYSQL_PASSWORD"),
		name:         testEnv("MYSQL_DATABASE"),
		port:         testEnv("MYSQL_PORT"),
		host:         testEnv("MYSQL_HOST"),
		unixPath:     testEnv("MYSQL_UNIX_SOCKET"),
		instConnName: testEnv("MYSQL_INSTANCE"),
	}
	// Zero out all but requested conn type
	switch ct {
	case useTCP:
		// use host
		d.unixPath = ""
		d.instConnName = ""
	case useUnix:
		d.host = ""
		// use unix path
		d.instConnName = ""
	default: // connector
		d.host = ""
		d.unixPath = ""
		// use instance connection name
	}
	return d
}

func setupTestEnv(c dbConfig) func() {
	oldDBUser := os.Getenv("DB_USER")
	oldDBPass := os.Getenv("DB_PASS")
	oldDBName := os.Getenv("DB_NAME")
	oldDBPort := os.Getenv("DB_PORT")
	oldHost := os.Getenv("INSTANCE_HOST")
	oldUnix := os.Getenv("INSTANCE_UNIX_SOCKET")
	oldInstance := os.Getenv("INSTANCE_CONNECTION_NAME")

	os.Setenv("DB_USER", c.user)
	os.Setenv("DB_PASS", c.pass)
	os.Setenv("DB_NAME", c.name)
	os.Setenv("DB_PORT", c.port)
	os.Setenv("INSTANCE_HOST", c.host)
	os.Setenv("INSTANCE_UNIX_SOCKET", c.unixPath)
	os.Setenv("INSTANCE_CONNECTION_NAME", c.instConnName)

	return func() {
		os.Setenv("DB_USER", oldDBUser)
		os.Setenv("DB_PASS", oldDBPass)
		os.Setenv("DB_NAME", oldDBName)
		os.Setenv("DB_PORT", oldDBPort)
		os.Setenv("INSTANCE_HOST", oldHost)
		os.Setenv("INSTANCE_UNIX", oldUnix)
		os.Setenv("INSTANCE_CONNECTION_NAME", oldInstance)
	}
}

func testGetVotes(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	Votes(rr, req)
	resp := rr.Result()
	body := rr.Body.String()

	wantStatus := 200
	if gotStatus := resp.StatusCode; wantStatus != gotStatus {
		t.Errorf("want = %v, got = %v", wantStatus, gotStatus)
	}

	want := "Tabs VS Spaces"
	if !strings.Contains(body, want) {
		t.Errorf("failed to find %q in resp = %v", want, body)
	}
}

func testCastVote(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("team=SPACES")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	Votes(rr, req)
	resp := rr.Result()
	body := rr.Body.String()

	wantStatus := 200
	if gotStatus := resp.StatusCode; wantStatus != gotStatus {
		t.Errorf("want = %v, got = %v", wantStatus, gotStatus)
	}

	want := "Vote successfully cast for SPACES"
	if !strings.Contains(body, want) {
		t.Errorf("failed to find %q in resp = %v", want, body)
	}
}

func TestGetVotes(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip()
	}

	testCases := []struct {
		desc string
		ct   connType
	}{
		{desc: "connecting with a TCP Socket", ct: useTCP},
		{desc: "connecting with a Unix Socket", ct: useUnix},
		{desc: "connecting with a connector", ct: useConnector},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			conf := dbConfigFromEnv(t, tc.ct)
			cleanup := setupTestEnv(conf)
			defer cleanup()

			// initialize database connection based on environment
			db = mustConnect()

			testGetVotes(t)
		})
	}
}

func TestCastVote(t *testing.T) {
	if os.Getenv("GOLANG_SAMPLES_E2E_TEST") == "" {
		t.Skip()
	}

	testCases := []struct {
		desc string
		ct   connType
	}{
		{desc: "connecting with a TCP Socket", ct: useTCP},
		{desc: "connecting with a Unix Socket", ct: useUnix},
		{desc: "connecting with a connector", ct: useConnector},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			conf := dbConfigFromEnv(t, tc.ct)
			cleanup := setupTestEnv(conf)
			defer cleanup()

			// initialize database connection based on environment
			db = mustConnect()

			testCastVote(t)
		})
	}
}
