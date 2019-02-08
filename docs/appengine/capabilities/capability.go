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

package sample

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/capability"
)

// [START gae_go_capabilities_lookup]
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if !capability.Enabled(ctx, "datastore_v3", "*") {
		http.Error(w, "This service is currently unavailable.", 503)
		return
	}
	// do Datastore lookup ...
}

// [END gae_go_capabilities_lookup]

func example() {
	var ctx context.Context
	// [START gae_go_capabilities_mode]
	if !capability.Enabled(ctx, "datastore_v3", "write") {
		// Datastore is in read-only mode.
	}
	// [END gae_go_capabilities_mode]
}
