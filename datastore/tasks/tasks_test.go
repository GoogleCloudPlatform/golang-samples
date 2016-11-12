// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var client *datastore.Client

func TestMain(m *testing.M) {
	ctx := context.Background()
	if tc, ok := testutil.ContextMain(m); ok {
		var err error
		client, err = datastore.NewClient(ctx, tc.ProjectID)
		if err != nil {
			log.Fatalf("datastore.NewClient: %v", err)
		}
		defer client.Close()
	}
	os.Exit(m.Run())
}

func makeDesc() string {
	return fmt.Sprintf("t-%d", time.Now().Unix())
}

func TestAddMarkDelete(t *testing.T) {
	testutil.SystemTest(t)
	ctx := context.Background()

	desc := makeDesc()

	k, err := AddTask(ctx, client, desc)
	if err != nil {
		t.Fatal(err)
	}

	if err := MarkDone(ctx, client, k.ID); err != nil {
		t.Fatal(err)
	}

	if err := DeleteTask(ctx, client, k.ID); err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	t.Skip("Flaky. Eventual consistency. Re-enable once the datastore emulator works with gRPC.")

	testutil.SystemTest(t)
	ctx := context.Background()

	desc := makeDesc()

	k, err := AddTask(ctx, client, desc)
	if err != nil {
		t.Fatal(err)
	}

	foundTask := listAndGetTask(t, desc)
	if got, want := foundTask.id, k.ID; got != want {
		t.Errorf("k.ID: got %d, want %d", got, want)
	}

	if err := MarkDone(ctx, client, foundTask.id); err != nil {
		t.Fatal(err)
	}

	foundTask = listAndGetTask(t, desc)
	if !foundTask.Done {
		t.Error("foundTask.Done: got false, want true")
	}

	if err := DeleteTask(ctx, client, foundTask.id); err != nil {
		t.Fatal(err)
	}
}

func listAndGetTask(t *testing.T, desc string) *Task {
	ctx := context.Background()

	tasks, err := ListTasks(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	var foundTask *Task
	for _, t := range tasks {
		if t.Desc == desc {
			foundTask = t
		}
	}
	if foundTask == nil {
		t.Fatalf("Did not find task %s in list.", desc)
	}

	return foundTask
}
