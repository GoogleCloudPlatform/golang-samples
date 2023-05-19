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

// Sample simplelog writes some entries, lists them, then deletes the log.
package main

// [START logging_delete_log]
import (
	"context"
	"log"

	"cloud.google.com/go/logging/logadmin"
)

func deleteLog(projectID string) error {
	ctx := context.Background()
	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logadmin client: %v", err)
	}
	defer adminClient.Close()

	const name = "log-example"
	if err := adminClient.DeleteLog(ctx, name); err != nil {
		return err
	}
	return nil
}

// [END logging_delete_log]
