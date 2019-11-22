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

// Package cloudrunci_test holds tests for the cloudrunci package.
package cloudrunci_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestManagedPlatformRequest(t *testing.T) {
	testutil.EndToEndTest(t)
	p := cloudrunci.ManagedPlatform{}

	req, err := p.NewRequest("GET", "http://example.com")
	if err != nil {
		t.Errorf("ManagedPlatform.Request: %q", err)
	}
	authzHeader := req.Header.Get("Authorization")
	if authzHeader == "" {
		t.Errorf("ManagedPlatform.Request: missing authentication header: %q", err)
	}
}

func TestGKEPlatformRequest(t *testing.T) {
	testutil.EndToEndTest(t)
	p := cloudrunci.GKEPlatform{}

	req, err := p.NewRequest("GET", "http://example.com")
	if err != nil {
		t.Errorf("GKEPlatform.Request: %q", err)
	}
	authzHeader := req.Header.Get("Authorization")
	if authzHeader != "" {
		t.Errorf("GKEPlatform.Request: unexpected authentication header: %q", err)
	}
}

func TestKubernetesPlatformRequest(t *testing.T) {
	testutil.EndToEndTest(t)
	p := cloudrunci.KubernetesPlatform{}

	req, err := p.NewRequest("GET", "http://example.com")
	if err != nil {
		t.Errorf("KubernetesPlatform.Request: %q", err)
	}
	authzHeader := req.Header.Get("Authorization")
	if authzHeader != "" {
		t.Errorf("KubernetesPlatform.Request: unexpected authentication header: %q", err)
	}
}
