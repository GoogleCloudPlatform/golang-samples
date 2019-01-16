// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
)

const topicName = "functions-test-topic"

func TestPublishMessage(t *testing.T) {
	// TODO: Use testutil.
	projectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	ctx := context.Background()
	var err error
	client, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("pubsub.NewClient: %v", err)
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

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", strings.NewReader(fmt.Sprintf(`{"topic":%q}`, topicName)))
	PublishMessage(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PublishMessage got response code %v, want %v", rr.Code, http.StatusOK)
	}

	want := "Published"
	if got := rr.Body.String(); !strings.Contains(got, want) {
		t.Errorf("PublishMessage got %q, want to contain %q", got, want)
	}
}
