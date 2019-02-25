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

package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestIntentManagement(t *testing.T) {
	t.Skip("flaky")
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	parent := fmt.Sprintf("projects/%s/agents", projectID)

	intentDisplayNames := [...]string{fmt.Sprintf("intent-%s", time.Now()), fmt.Sprintf("intent-%s", time.Now())}

	initialIntents, err := ListIntents(projectID)

	if err != nil {
		t.Error("Unsuccessful initial ListIntents")
	}

	for _, displayName := range intentDisplayNames {
		trainingPhraseParts := []string{fmt.Sprintf("%s-phrase-%s", displayName, time.Now()), fmt.Sprintf("%s-phrase-%s", displayName, time.Now())}
		messageTexts := []string{fmt.Sprintf("%s-message-%s", displayName, time.Now()), fmt.Sprintf("%s-message-%s", displayName, time.Now())}
		err = CreateIntent(projectID, displayName, trainingPhraseParts, messageTexts)
		if err != nil {
			t.Errorf("Unsuccessful intent creation: %s under %s", displayName, parent)
		}
	}

	intermediateIntents, err := ListIntents(projectID)

	if err != nil {
		t.Error("Unsuccessful intermediate ListIntents")
	}

	if len(intermediateIntents) != len(initialIntents)+len(intentDisplayNames) {
		t.Errorf("len(intermediateIntents) = %d; want %d", len(intermediateIntents), len(initialIntents)+len(intentDisplayNames))
	}

	for _, displayName := range intentDisplayNames {
		var intentID string
		for _, intent := range intermediateIntents {
			if intent.GetDisplayName() == displayName {
				route := strings.Split(intent.GetName(), "/")
				intentID = route[len(route)-1]
			}
		}
		if intentID == "" {
			t.Error("intentID empty; want non-empty")
		}

		err = DeleteIntent(projectID, intentID)
		if err != nil {
			t.Errorf("Unsuccessful intent deletion %s/intent/%s", parent, intentID)
		}
	}

	finalIntents, err := ListIntents(projectID)

	if err != nil {
		t.Error("Unsuccessful final ListIntents")
	}

	if len(finalIntents) != len(initialIntents) {
		t.Errorf("Actual len(finalIntents) = %d; want %d", len(finalIntents), len(initialIntents))
	}
}
