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

package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

func setup(ctx context.Context, t *testing.T) (*firestore.Client, string) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	return client, projectID
}

func TestListen(t *testing.T) {
	ctx := context.Background()
	client, projectID := setup(ctx, t)
	defer client.Close()

	cityCollection := []struct {
		city, name, state string
	}{
		{city: "SF", name: "San Francisco", state: "CA"},
		{city: "LA", name: "Los Angeles", state: "CA"},
		{city: "DC", name: "Washington D.C."},
	}

	for _, c := range cityCollection {
		if _, err := client.Collection("cities").Doc(c.city).Set(ctx, map[string]string{
			"name":  c.name,
			"state": c.state,
		}); err != nil {
			log.Fatalf("Set: %v", err)
		}
	}
	if err := listenDocument(ioutil.Discard, projectID); err != nil {
		t.Errorf("listenDocument: %v", err)
	}
}
func TestListenMultiple(t *testing.T) {
	ctx := context.Background()
	client, projectID := setup(ctx, t)
	defer client.Close()

	if err := listenMultiple(ioutil.Discard, projectID); err != nil {
		t.Errorf("listenMultiple: %v", err)
	}
}

func TestListenChanges(t *testing.T) {
	ctx := context.Background()
	client, projectID := setup(ctx, t)
	defer client.Close()

	if err := listenChanges(ioutil.Discard, projectID); err != nil {
		t.Errorf("listenChanges: %v", err)
	}
}

func TestListenErrors(t *testing.T) {
	ctx := context.Background()
	client, projectID := setup(ctx, t)
	defer client.Close()

	if err := listenErrors(ioutil.Discard, projectID); err != nil {
		t.Errorf("listenErrors: %v", err)
	}
}
