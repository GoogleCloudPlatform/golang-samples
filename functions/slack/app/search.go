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

// [START functions_slack_search]

package slack

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// KGSearch uses the Knowledge Graph API
func KGSearch(w http.ResponseWriter, r *http.Request) {
	setup(r.Context())
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}

	d := json.NewDecoder(cfgFile)
	config = &configuration{}
	if err = d.Decode(config); err != nil {
		log.Fatalf("Decode: %v", err)
	}
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}
	if err = r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}
	if err = verifyWebHook(r.Form); err != nil {
		log.Fatalf("verifyWebhook: %v", err)
	}
	kgSearchResponse, err := makeSearchRequest(r.Form["text"][0])
	if err != nil {
		log.Fatalf("makeSearchRequest: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(kgSearchResponse); err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
}

// [END functions_slack_search]
