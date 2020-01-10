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

// [START functions_slack_setup]

package slack

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
)

type configuration struct {
	ProjectID string `json:"PROJECT_ID"`
	Token     string `json:"SLACK_TOKEN"`
	Key       string `json:"KG_API_KEY"`
	Secret    string `json:"SLACK_SIGNING_SECRET"`
}

var (
	entitiesService *kgsearch.EntitiesService
	config          *configuration
)

func setup(ctx context.Context) {
	if config == nil {
		cfgFile, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("os.Open: %v", err)
		}

		d := json.NewDecoder(cfgFile)
		config = &configuration{}
		if err = d.Decode(config); err != nil {
			log.Fatalf("Decode: %v", err)
		}
	}

	if entitiesService == nil {
		kgService, err := kgsearch.NewService(ctx, option.WithAPIKey(config.Key))
		if err != nil {
			log.Fatalf("kgsearch.NewClient: %v", err)
		}
		entitiesService = kgsearch.NewEntitiesService(kgService)
	}
}

// [END functions_slack_setup]
