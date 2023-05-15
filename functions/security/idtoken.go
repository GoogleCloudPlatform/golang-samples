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
// [START cloudrun_service_to_service_auth]

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/idtoken"
)

// `makeGetRequest` makes a request to the provided `targetURL`
// with an authenticated client using audience `audience`.
func makeGetRequest(w io.Writer, targetURL string, audience string) error {
	// [END functions_bearer_token]
	// Example `audience` value (Cloud Run): https://my-cloud-run-service.run.app/
	// (`targetURL` and `audience` will differ for non-root URLs and GET parameters)
	// [END cloudrun_service_to_service_auth]
	// [START functions_bearer_token]
	// For Cloud Functions, endpoint (`serviceUrl`) and `audience` are the same.
	// Example `audience` value (Cloud Functions): https://<PROJECT>-<REGION>-<PROJECT_ID>.cloudfunctions.net/myFunction
	// (`targetURL` and `audience` will differ for GET parameters)
	// [START cloudrun_service_to_service_auth]
	ctx := context.Background()

	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	client, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return fmt.Errorf("idtoken.NewClient: %w", err)
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return fmt.Errorf("client.Get: %w", err)
	}
	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	return nil
}

// [END cloudrun_service_to_service_auth]
// [END functions_bearer_token]
