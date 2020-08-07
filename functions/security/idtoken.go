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

// Package security contains samples for securely calling functions.
package security

// [START functions_bearer_token]
// [START run_service_to_service_auth]

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/idtoken"
)

// makeGetRequest makes a request to the provided targetURL with an authenticated client.
func makeGetRequest(w io.Writer, targetURL string) error {
	// functionURL := "https://TARGET_URL"
	ctx := context.Background()

	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	client, err := idtoken.NewClient(ctx, targetURL)
	if err != nil {
		return fmt.Errorf("idtoken.NewClient: %v", err)
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return fmt.Errorf("client.Get: %v", err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	return nil
}

// [END run_service_to_service_auth]
// [END functions_bearer_token]
