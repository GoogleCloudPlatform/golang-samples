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

package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	app := startApp()
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	app.indexHandler(writer, request)
	resp := writer.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Errorf("Expected StatusCode of 200. Got %d", resp.StatusCode)
	}
	if !strings.Contains(string(body[:]), "Tabs VS Spaces") {
		t.Errorf("Expected to see 'Tabs VS Spaces' in index response body")
	}
}

func TestCastVote(t *testing.T) {
	app := startApp()
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("team=SPACES")))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.indexHandler(writer, request)
	resp := writer.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Errorf("Expected StatusCode of 200. Got %d", resp.StatusCode)
	}
	if !strings.Contains(string(body[:]), "Vote successfully cast for SPACES") {
		t.Errorf("Expected to see 'Vote successfully cast for SPACES' in response body")
	}
}
