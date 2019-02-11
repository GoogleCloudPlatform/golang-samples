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
	"net/http"

	"google.golang.org/appengine/log"
)

// [START communication_between_modules_1]
import "google.golang.org/appengine"

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	module := appengine.ModuleName(ctx)
	instance := appengine.InstanceID()

	log.Infof(ctx, "Received on module %s; instance %s", module, instance)
}

// [END communication_between_modules_1]

func handler2(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// [START communication_between_modules_2]
	hostname, err := appengine.ModuleHostname(ctx, "my-backend", "", "")
	if err != nil {
		// ...
	}
	url := "http://" + hostname + "/"
	// ...
	// [END communication_between_modules_2]

	_ = url
}
