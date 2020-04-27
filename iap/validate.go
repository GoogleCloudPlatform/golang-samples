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

	"google.golang.org/api/idtoken"
)

// validateJWTFromAppEngine validates a JWT found in the
// "x-goog-iap-jwt-assertion" header.
func validateJWTFromAppEngine(iapJWT, projectNumber, projectID string) error {
	ctx := context.Background()
	aud := fmt.Sprintf("/projects/%s/apps/%s", projectNumber, projectID)

	payload, err := idtoken.Validate(ctx, iapJWT, aud)
	if err != nil {
		return err
	}

	// payload contains the JWT claims for further inspection or validation
	_ = payload

	return nil
}

// validateJWTFromComputeEngine validates a JWT found in the
// "x-goog-iap-jwt-assertion" header.
func validateJWTFromComputeEngine(iapJWT, projectNumber, backendServiceID string) error {
	ctx := context.Background()
	aud := fmt.Sprintf("/projects/%s/global/backendServices/%s", projectNumber, backendServiceID)

	payload, err := idtoken.Validate(ctx, iapJWT, aud)
	if err != nil {
		return err
	}

	// payload contains the JWT claims for further inspection or validation
	_ = payload

	return nil
}
// [END iap_validate_jwt]
