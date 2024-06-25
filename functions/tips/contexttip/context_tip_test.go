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

package contexttip

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestPublishMessage(t *testing.T) {
	tc := testutil.SystemTest(t)
	os.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	ctx := context.Background()
	var err error
	client, err = pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("pubsub.NewClient: %v", err)
	}

	topicName := os.Getenv("FUNCTIONS_TOPIC_NAME")
	if topicName == "" {
		topicName = "functions-test-topic"
	}

	topic := client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		t.Fatalf("topic(%s).Exists: %v", topicName, err)
	}
	if !exists {
		_, err = client.CreateTopic(context.Background(), topicName)
		if err != nil {
			t.Fatalf("topic(%s).CreateTopic: %v", topicName, err)
		}
	}

	payload := strings.NewReader(fmt.Sprintf(`{"topic":%q, "message": %q}`, topicName, "my_message"))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", payload)
	PublishMessage(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PublishMessage: got response code %v, want %v", rr.Code, http.StatusOK)
	}

	want := "published"
	if got := rr.Body.String(); !strings.Contains(got, want) {
		t.Errorf("PublishMessage: got %q, want to contain %q", got, want)
	}
}
