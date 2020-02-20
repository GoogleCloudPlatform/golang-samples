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

// [START run_secure_request]
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/compute/metadata"
)

// RenderService represents our upstream render service.
type RenderService struct {
	// URL is the render service address.
	URL string
	// Authenticated determines whether identity token authentication will be used.
	Authenticated bool
}

// NewRequest creates a new HTTP request with IAM ID Token credential.
// This token is automatically handled by private Cloud Run (fully managed) and Cloud Functions.
func (s *RenderService) NewRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(method, s.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}

	// Skip authentication if not using HTTPS, such as for local development.
	if !s.Authenticated {
		return req, nil
	}

	// Query the id_token with ?audience as the serviceURL
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", s.URL)
	token, err := metadata.Get(tokenURL)
	if err != nil {
		return req, fmt.Errorf("metadata.Get: failed to query id_token: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, nil
}

// [END run_secure_request]

// [START run_secure_request_do]

var renderClient = &http.Client{Timeout: 30 * time.Second}

// Render converts the Markdown plaintext to HTML.
func (s *RenderService) Render(in []byte) ([]byte, error) {
	req, err := s.NewRequest(http.MethodPost)
	if err != nil {
		return nil, fmt.Errorf("RenderService.NewRequest: %v", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(in))
	defer req.Body.Close()

	resp, err := renderClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do: %v", err)
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("http.Client.Do: %s (%d): request not OK", http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	return out, nil
}

// [END run_secure_request_do]
