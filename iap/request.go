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

// [START iap_make_request]
import (
	"context"
	"net/http"

	"google.golang.org/api/idtoken"
)

// makeIAPRequest makes a request to an application protected by Identity-Aware
// Proxy with the given iapClientID.
func makeIAPRequest(request *http.Request, iapClientID string) error {
	ctx := context.Background()
	client, err := idtoken.NewClient(ctx, iapClientID)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// Use the response
	_ = response

	return nil
}
// [END iap_make_request]
