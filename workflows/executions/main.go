// Copyright 2025 Google LLC
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

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	// Get environment variables.
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	workflowID := os.Getenv("WORKFLOW")
	locationID := os.Getenv("LOCATION")

	buf := bytes.Buffer{}

	// Execute workflow.
	if err := executeWorkflow(&buf, projectID, workflowID, locationID); err != nil {
		log.Fatalf("Error when executing workflow: %v", err)
	}

	fmt.Println(buf.String())
}
