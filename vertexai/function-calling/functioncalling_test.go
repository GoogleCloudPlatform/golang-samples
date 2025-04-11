// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package functioncalling

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_functionCalling(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	location := "us-central1"
	modelName := "gemini-1.5-flash-002"

	err := functionCalling(&buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("functionCalling failed: %v", err.Error())
	}

	funcOut := buf.String()

	expOut := `The model suggests to call the function "getCurrentWeather" with args: map[location:Boston]`
	if !strings.Contains(funcOut, expOut) {
		t.Errorf("expected output to contain text %q, got: %q", expOut, funcOut)
	}

	expOut = "weather in Boston"
	if !strings.Contains(funcOut, expOut) {
		t.Errorf("expected output to contain text %q, got: %q", expOut, funcOut)
	}
}

func Test_functionCallsChat(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	err := functionCallsChat(&buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("Test_functionCallsChat: %v", err.Error())
	}
}

func Test_parallelFunctionCalling(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	location := "us-central1"
	modelName := "gemini-1.5-flash-002"

	err := parallelFunctionCalling(&buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("parallelFunctionCalling failed: %v", err.Error())
	}

	funcOut := buf.String()
	testCases := []string{
		`The model suggests to call the function "getCurrentWeather" with args: map[location:New Delhi]`,
		`The model suggests to call the function "getCurrentWeather" with args: map[location:San Francisco]`,
		"weather in New Delhi",
		"weather in San Francisco",
	}

	for _, expOut := range testCases {
		if !strings.Contains(funcOut, expOut) {
			t.Errorf("expected output to contain text %q, got: %q", expOut, funcOut)
		}
	}
}
