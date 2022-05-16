// Copyright 2022 Google LLC
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
	"log"
	"os"

	"google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
)

var (
	entitiesService *kgsearch.EntitiesService
	kgKey           string
	slackSecret     string
)

func setup(ctx context.Context) {
	kgKey = os.Getenv("KG_API_KEY")
	slackSecret = os.Getenv("SLACK_SECRET")

	if entitiesService == nil {
		kgService, err := kgsearch.NewService(ctx, option.WithAPIKey(kgKey))
		if err != nil {
			log.Fatalf("kgsearch.NewService: %v", err)
		}
		entitiesService = kgsearch.NewEntitiesService(kgService)
	}
}

// [END functions_slack_setup]
