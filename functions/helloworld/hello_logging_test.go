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

package helloworld

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestHelloLogging(t *testing.T) {
	rStd, wStd, _ := os.Pipe()
	stdLogger.SetOutput(wStd)

	rErr, wErr, _ := os.Pipe()
	logger.SetOutput(wErr)

	HelloLogging(nil, nil)

	wStd.Close()
	wErr.Close()

	stdout, err := ioutil.ReadAll(rStd)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := "log entry"
	if got := string(stdout); !strings.Contains(got, want) {
		t.Errorf("Stdout got %q, want to contain %q", got, want)
	}

	stderr, err := ioutil.ReadAll(rErr)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want = "error"
	if got := string(stderr); !strings.Contains(got, want) {
		t.Errorf("Stderr got %q, want to contain %q", got, want)
	}
}
