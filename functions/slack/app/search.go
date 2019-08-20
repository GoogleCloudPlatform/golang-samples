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
	"fmt"
	"log"
	"net/http"
	"os"

	slack "github.com/nlopes/slack"
	"google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
)

// KGSearch uses the Knowledge Graph API
func KGSearch(w http.ResponseWriter, r *http.Request) {
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}

	d := json.NewDecoder(cfgFile)
	config = &configType{}
	err = d.Decode(config)
	if err != nil {
		log.Fatalf("Decode: %v", err)
	}

	slackClient = slack.New(config.Key)
	kgService, err := kgsearch.NewService(r.Context(), option.WithAPIKey(config.Key))
	if err != nil {
		log.Fatalf("kgsearch.NewClient: %v", err)
	}
	kgService.Entities = kgsearch.NewEntitiesService(kgService)
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}
	err = verifyWebhook(r.Form)
	if err != nil {
		log.Fatalf("verifyWebhook: %v", err)
	}
	kgSearchResponse, err := makeSearchRequest(kgService.Entities, r.Form["text"][0])
	if err != nil {
		log.Fatalf("makeSearchRequest: %v", err)
	}
	resp, err := json.Marshal(kgSearchResponse)
	if err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
	fmt.Fprint(w, string(resp))
}

// [END functions_slack_search]
