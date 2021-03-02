// Copyright 2021 Google LLC
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

// package subscriptions is a tool to manage Google Cloud Pub/Sub subscriptions by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package pslite

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsublite"
)

const (
	testRegion = "us-central1"
)

func testAdminClient(t *testing.T) (*pubsublite.AdminClient, context.Context) {
	t.Helper()

	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, testRegion)
	if err != nil {
		t.Fatalf("testClient: failed to create client: %v", err)
	}
	return client, ctx
}

func TestCreateTopic(t *testing.T) {
	admin, _ := testAdminClient(t)
	admin.CreateTopic()

}
