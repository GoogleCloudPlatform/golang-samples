// Copyright 2018, Google, LLC.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// [START functions_slack_search]
func kgSearch(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Couldn't parse form", 405)
		return nil, fmt.Errorf("ParseForm: %v", err)
	}
	err = verifyWebhook(r.Form)
	if err != nil {
		return nil, fmt.Errorf("verifyWebhook: %v", err)
	}
	kgSearchResponse, err := makeSearchRequest(r.Form["text"][0])
	return json.Marshal(kgSearchResponse)
}

// [END functions_slack_search]
