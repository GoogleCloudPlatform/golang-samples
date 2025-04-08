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

package iap

// [START iap_validate_jwt]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/idtoken"
)

// validateJWTFromAppEngine validates a JWT found in the
// "x-goog-iap-jwt-assertion" header.
func validateJWTFromAppEngine(w io.Writer, iapJWT, projectNumber, projectID string) error {
	// iapJWT := "YmFzZQ==.ZW5jb2RlZA==.and0" // req.Header.Get("X-Goog-IAP-JWT-Assertion")
	// projectNumber := "123456789"
	// projectID := "your-project-id"
	ctx := context.Background()
	aud := fmt.Sprintf("/projects/%s/apps/%s", projectNumber, projectID)

	payload, err := idtoken.Validate(ctx, iapJWT, aud)
	if err != nil {
		return fmt.Errorf("idtoken.Validate: %w", err)
	}

	// payload contains the JWT claims for further inspection or validation
	fmt.Fprintf(w, "payload: %v", payload)

	return nil
}

// validateJWTFromComputeEngine validates a JWT found in the
// "x-goog-iap-jwt-assertion" header.
func validateJWTFromComputeEngine(w io.Writer, iapJWT, projectNumber, backendServiceID string) error {
	// iapJWT := "YmFzZQ==.ZW5jb2RlZA==.and0" // req.Header.Get("X-Goog-IAP-JWT-Assertion")
	// projectNumber := "123456789"
	// backendServiceID := "backend-service-id"
	ctx := context.Background()
	aud := fmt.Sprintf("/projects/%s/global/backendServices/%s", projectNumber, backendServiceID)

	payload, err := idtoken.Validate(ctx, iapJWT, aud)
	if err != nil {
		return fmt.Errorf("idtoken.Validate: %w", err)
	}

	// payload contains the JWT claims for further inspection or validation
	fmt.Fprintf(w, "payload: %v", payload)

	return nil
}

// [END iap_validate_jwt]
