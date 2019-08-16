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

// [START functions_slack_setup]

import (
	"context"
	"encoding/json"
	"log"
	"os"

	slack "github.com/nlopes/slack"
	"google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
)

type configType struct {
	ProjectID string `json:"PROJECT_ID"`
	Token     string `json:"SLACK_TOKEN"`
	Key       string `json:"KG_API_KEY"`
}

// // GCSEvent is the payload of a GCS event.
// type GCSEvent struct {
// 	Bucket         string    `json:"bucket"`
// 	Name           string    `json:"name"`
// 	Metageneration string    `json:"metageneration"`
// 	ResourceState  string    `json:"resourceState"`
// 	TimeCreated    time.Time `json:"timeCreated"`
// 	Updated        time.Time `json:"updated"`
// }

type attachment struct {
	color     string
	title     string
	titleLink string
	text      string
	imageURL  string
}

// SlackMessage is the a Slack message event.
type SlackMessage struct {
	responseType string
	text         string
	attachments  []attachment
}

var (
	slackClient *slack.Client
	kgService   *kgsearch.EntitiesService
	config      *configType
)

func setup(ctx context.Context) {

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
	kgService, err := kgsearch.NewService(ctx, option.WithAPIKey(config.Key))
	if err != nil {
		log.Fatalf("kgsearch.NewClient: %v", err)
	}
	kgService.Entities = kgsearch.NewEntitiesService(kgService)
}

// [END functions_slack_setup]
