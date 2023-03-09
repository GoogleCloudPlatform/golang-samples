// Copyright 2020 Google LLC
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

package firestore

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var duration = 20 * time.Second

func setup(ctx context.Context, t *testing.T) (*firestore.Client, string, string) {
	tc := testutil.SystemTest(t)
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}
	collection := tc.ProjectID + "-collection-cities"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	return client, projectID, collection
}

func TestListen(t *testing.T) {
	ctx := context.Background()
	client, projectID, collection := setup(ctx, t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	// Delete all docs first to make sure setup works.
	docs, err := client.Collection(collection).Documents(ctx).GetAll()
	if err == nil {
		for _, doc := range docs {
			doc.Ref.Delete(ctx)
		}
	}
	cityCollection := []struct {
		city, name, state string
	}{
		{city: "SF", name: "San Francisco", state: "CA"},
		{city: "LA", name: "Los Angeles", state: "CA"},
		{city: "DC", name: "Washington D.C."},
	}

	for _, c := range cityCollection {
		if _, err := client.Collection(collection).Doc(c.city).Set(ctx, map[string]string{
			"name":  c.name,
			"state": c.state,
		}); err != nil {
			t.Fatalf("Set: %v", err)
		}
	}
	if err := listenDocument(ctx, ioutil.Discard, projectID, collection); err != nil {
		t.Errorf("listenDocument: %v", err)
	}
}
func TestListenMultiple(t *testing.T) {
	ctx := context.Background()
	client, projectID, collection := setup(ctx, t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	if err := listenMultiple(ctx, ioutil.Discard, projectID, collection); err != nil {
		t.Errorf("listenMultiple: %v", err)
	}
}

func TestListenChanges(t *testing.T) {
	ctx := context.Background()
	client, projectID, collection := setup(ctx, t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		c := make(chan *bytes.Buffer)
		go func() {
			defer close(c)
			err := listenChanges(ctx, buf, projectID, collection)
			if err != nil {
				t.Errorf("listenChanges: %v", err)
			}
			c <- buf
		}()
		// Add some changes to data in parallel.
		time.Sleep(time.Second)
		var pop int64 = 3900000
		if _, err := client.Collection(collection).Doc("LA").Update(ctx, []firestore.Update{
			{Path: "population", Value: pop},
		}); err != nil {
			log.Fatalf("Doc.Update: %v", err)
		}

		<-c

		time.Sleep(time.Second)

		want := "population:3900000"

		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("listenChanges got\n----\n%s\n----\nWant to contain:\n----\n%s\n----", got, want)
		}
	})
}

func TestListenErrors(t *testing.T) {
	ctx := context.Background()
	client, projectID, collection := setup(ctx, t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	if err := listenErrors(ctx, ioutil.Discard, projectID, collection); err != nil {
		t.Errorf("listenErrors: %v", err)
	}
}
