// Copyright 2023 Google LLC
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

package main

// [START logging_update_sink]
import (
	"context"
	"log"

	"cloud.google.com/go/logging/logadmin"
)

func updateSink(projectID string) (*logadmin.Sink, error) {
	ctx := context.Background()
	client, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("logadmin.NewClient: %v", err)
	}
	defer client.Close()
	sink, err := client.UpdateSink(ctx, &logadmin.Sink{
		ID:          "severe-errors-to-gcs",
		Destination: "storage.googleapis.com/logsinks-new-bucket",
		Filter:      "severity >= INFO",
	})
	return sink, err
}

// [END logging_update_sink]
