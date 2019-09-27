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

// [START functions_log_stackdriver]

// Package log contains examples for handling Cloud Functions logs.
package log

import (
	"context"
	"log"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ProcessLogEntry processes a Pub/Sub message from Stackdriver.
func ProcessLogEntry(ctx context.Context, m PubSubMessage) error {
	log.Printf("Log entry data: %s", string(m.Data))
	return nil
}

// [END functions_log_stackdriver]
