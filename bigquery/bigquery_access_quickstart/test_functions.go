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

package bigqueryaccess

import (
	"context"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/exp/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestClient(t *testing.T) (*bigquery.Client, error) {
	t.Helper()

	ctx := context.Background()
	tc := testutil.SystemTest(t)

	// Creates a client.
	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client, nil
}

func TestCleanup(t *testing.T, client *bigquery.Client, datasetName string) {
	t.Helper()

	ctx := context.Background()

	if err := client.Dataset(datasetName).DeleteWithContents(ctx); err != nil {
		t.Errorf("Failed to delete table: %v", err)
	}

	if err := client.Close(); err != nil {
		t.Fatalf("Failed to close Big Query client: %v", err)
	}
}
