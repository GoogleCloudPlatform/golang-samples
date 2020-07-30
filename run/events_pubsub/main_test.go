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

func TestHelloPubSubCloudEvent(t *testing.T) {
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

	got, e := HelloPubSub(context.Background(), *ce)
	if e != nil {
		t.Errorf("HelloPubSub: %q", e)
	}
	if want := "Hello, foo! ID: 321-CBA"; got != want {
		t.Errorf("HelloPubSub: got %q, want %q", got, want)
	}
}
