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

// Package iap contains Identity-Aware Proxy samples.
package iap

// [START iap_make_request]
import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/api/idtoken"
)

// makeIAPRequest makes a request to an application protected by Identity-Aware
// Proxy with the given audience.
func makeIAPRequest(w io.Writer, request *http.Request, audience string) error {
	// request, err := http.NewRequest("GET", "http://example.com", nil)
	// audience := "IAP_CLIENT_ID.apps.googleusercontent.com"
	ctx := context.Background()

	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return fmt.Errorf("idtoken.NewClient: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer response.Body.Close()
	if _, err := io.Copy(w, response.Body); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	return nil
}

// [END iap_make_request]
