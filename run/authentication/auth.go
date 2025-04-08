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

// Package authentication contains the authentication samples for Cloud Run.
package authentication

// [START cloudrun_service_to_service_auth]
import (
	"fmt"
	"net/http"

	"cloud.google.com/go/compute/metadata"
)

// makeGetRequest makes a GET request to the specified Cloud Run endpoint in
// serviceURL (must be a complete URL) by authenticating with the ID token
// obtained from the Metadata API.
func makeGetRequest(serviceURL string) (*http.Response, error) {
	// query the id_token with ?audience as the serviceURL
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", serviceURL)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("metadata.Get: failed to query id_token: %+v", err)
	}
	req, err := http.NewRequest("GET", serviceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", idToken))
	return http.DefaultClient.Do(req)
}

// [END cloudrun_service_to_service_auth]
