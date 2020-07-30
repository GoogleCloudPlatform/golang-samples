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
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func createPubSubCE() *cloudevents.Event {
	pubsubEvent := &PubSub{
		Message: PubSubMessage{
			Data: []byte("foo"),
			ID:   "id",
		},
		Subscription: "sub",
	}
	data, err := json.Marshal(pubsubEvent)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	ce := &cloudevents.Event{
		Context: cloudevents.EventContextV1{
			ID:              "321-CBA",
			Type:            "unit.test.client.response",
			Time:            &cloudevents.Timestamp{Time: time.Now()},
			Source:          *cloudevents.ParseURIRef("/unit/test/client"),
			DataContentType: cloudevents.StringOfApplicationJSON(),
		}.AsV1(),
		DataEncoded: data,
	}
	return ce
}

func TestEventsPubSubReceive(t *testing.T) {
	t.Skip("test requires Go 1.13+. See: https://github.com/GoogleCloudPlatform/golang-samples/issues/1224")
	ce := createPubSubCE()

	// Basic test
	HelloPubSub(context.Background(), *ce)
	if ce == nil {
		t.Error()
	}

	// TODO send CE with HTTP recorder
	// payload := strings.NewReader("foo")
	// req := httptest.NewRequest("POST", "/", payload)
	// rr := httptest.NewRecorder()
	// HelloPubSub(rr, *ce)

	// if e != nil {
	// 	t.Errorf("HelloPubSub: %q", e)
	// }
	// if want := "Hello, foo! ID: 321-CBA"; got != want {
	// 	t.Errorf("HelloPubSub: got %q, want %q", got, want)
	// }
}

func TestEventsPubSubLocalSend(t *testing.T) {
	// The default client is HTTP.
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/")
	ce := createPubSubCE()

	// Send that Event.
	if result := c.Send(ctx, *ce); cloudevents.IsUndelivered(result) {
		log.Fatalf("failed to send, %v", result)
		t.Error("Failed to send event.")
	}
}
