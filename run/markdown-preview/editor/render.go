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

// [START cloudrun_secure_request]
// [START run_secure_request]
import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

// RenderService represents our upstream render service.
type RenderService struct {
	// URL is the render service address.
	URL string
	// tokenSource provides an identity token for requests to the Render Service.
	tokenSource oauth2.TokenSource
}

// NewRequest creates a new HTTP request to the Render service.
// If authentication is enabled, an Identity Token is created and added.
func (s *RenderService) NewRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(method, s.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a TokenSource if none exists.
	if s.tokenSource == nil {
		s.tokenSource, err = idtoken.NewTokenSource(ctx, s.URL)
		if err != nil {
			return nil, fmt.Errorf("idtoken.NewTokenSource: %w", err)
		}
	}

	// Retrieve an identity token. Will reuse tokens until refresh needed.
	token, err := s.tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("TokenSource.Token: %w", err)
	}
	token.SetAuthHeader(req)

	return req, nil
}

// [END run_secure_request]
// [END cloudrun_secure_request]

// [START cloudrun_secure_request_do]
// [START run_secure_request_do]

var renderClient = &http.Client{Timeout: 30 * time.Second}

// Render converts the Markdown plaintext to HTML.
func (s *RenderService) Render(in []byte) ([]byte, error) {
	req, err := s.NewRequest(http.MethodPost)
	if err != nil {
		return nil, fmt.Errorf("RenderService.NewRequest: %w", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(in))
	defer req.Body.Close()

	resp, err := renderClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do: %w", err)
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("http.Client.Do: %s (%d): request not OK", http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	return out, nil
}

// [END run_secure_request_do]
// [END cloudrun_secure_request_do]
