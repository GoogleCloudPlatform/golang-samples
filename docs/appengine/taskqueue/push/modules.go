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

package counter

// [START push_queues_and_backends]
import (
	"net/http"
	"net/url"

	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

func pushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	key := r.FormValue("key")

	// Create a task pointed at a backend.
	t := taskqueue.NewPOSTTask("/path/to/my/worker/", url.Values{
		"key": {key},
	})
	host, err := appengine.ModuleHostname(ctx, "backend1", "", "")
	if err != nil {
		// Handle err
	}
	t.Header = http.Header{
		"Host": {host},
	}

	// Add the task to the default queue.
	if _, err := taskqueue.Add(ctx, t, ""); err != nil {
		// Handle err
	}
}

// [END push_queues_and_backends]
