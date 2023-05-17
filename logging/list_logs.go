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

// Package snippets contains examples of using the Cloud Logging API.
package snippets

// [START logging_list_logs]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

// listLogs lists all available logs in the project.
func listLogs(w io.Writer, projectID string) error {
	ctx := context.Background()

	client, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	iter := client.Logs(ctx)
	for {
		log, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("List logs failed: %w", err)
		}
		fmt.Fprintf(w, "%s\n", log)
	}
	return nil
}

// [END logging_list_logs]
