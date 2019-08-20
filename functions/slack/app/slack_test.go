// Copyright 2019 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slack

import (
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// TestMain sets up the config rather than using the config file
// which contains placeholder values.
func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatalf("testutil.ContextMain failed")
	}
	config = &configuration{
		ProjectID: tc.ProjectID,
		Token:     "mLdpMlnwPr1mcVgusUGR8VBj",
		Key:       "AIzaSyBFrpWs3otLxuWYJkk2AXQ3Xi1OB2oTi0A",
	}

	os.Exit(m.Run())
}

func TestVerifyWebHook(t *testing.T) {
	v := make(url.Values)
	v["token"] = []string{config.Token}
	err := verifyWebHook(v)
	if err != nil {
		t.Errorf("verifyWebHook: %v", err)
	}
	v = make(url.Values)
	v["token"] = []string{"this is not the token"}
	err = verifyWebHook(v)
	if err == nil {
		t.Errorf("got %q, want %q", "nil", "invalid request/credentials")
	}
	v = make(url.Values)
	v["token"] = []string{""}
	err = verifyWebHook(v)
	if err == nil {
		t.Errorf("got %q, want %q", "nil", "empty form token")
	}
}
