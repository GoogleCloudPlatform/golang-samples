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

// Simple CLI to run the executeWorkflow function.
// Used for one-off testing and development.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: Must include 2 arguments for projectID, locationID")
		os.Exit(1)
	}
	projectID := os.Args[1]
	locationID := os.Args[2]
	workflowID := ""
	if len(os.Args) > 3 {
		workflowID = os.Args[3]
	}

	res, err := executeWorkflow(projectID, locationID, workflowID)
	if err != nil {
		log.Fatalf("executeWorkflow: %v", err)
	}
	var jsonStringArr []string
	err = json.Unmarshal([]byte(res), &jsonStringArr)

	fmt.Print(jsonStringArr)
}
