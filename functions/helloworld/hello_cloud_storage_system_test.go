// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build ignore
// Disabled until system tests are working on Kokoro.

// TODO: use testutil.SystemTest

// [START functions_storage_system_test]

package helloworld

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
)

func TestHelloGCSSystem(t *testing.T) {
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME")

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}

	// Create a file.
	startTime := time.Now().UTC().Format(time.RFC3339)
	oh := client.Bucket(bucketName).Object("TestHelloGCSSystem")
	w := oh.NewWriter(ctx)
	fmt.Fprintf(w, "Content of the file")
	w.Close()

	// Wait for logs to be consistent.
	time.Sleep(20 * time.Second)

	// Check logs.
	want := "created"
	if got := readLogs(t, startTime); !strings.Contains(got, want) {
		t.Errorf("HelloGCS logged %q, want to contain %q", got, want)
	}

	// Modify the file.
	startTime = time.Now().UTC().Format(time.RFC3339)
	_, err = oh.Update(ctx, storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{"Content-Type": "text/html"},
	})
	if err != nil {
		t.Errorf("Update: %v", err)
	}

	// Wait for logs to be consistent.
	time.Sleep(20 * time.Second)

	// Check logs.
	want = "updated"
	if got := readLogs(t, startTime); !strings.Contains(got, want) {
		t.Errorf("HelloGCS logged %q, want to contain %q", got, want)
	}

	// Delete the file.
	startTime = time.Now().UTC().Format(time.RFC3339)
	if err := oh.Delete(ctx); err != nil {
		t.Errorf("Delete: %v", err)
	}

	// Wait for logs to be consistent.
	time.Sleep(20 * time.Second)

	// Check logs.
	want = "deleted"
	if got := readLogs(t, startTime); !strings.Contains(got, want) {
		t.Errorf("HelloGCS logged %q, want to contain %q", got, want)
	}
}

func readLogs(t *testing.T, startTime string) string {
	t.Helper()
	cmd := exec.Command("gcloud", "alpha", "functions", "logs", "read", "HelloGCS", "--start-time", startTime)
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec.Command: %v", err)
	}
	return string(got)
}

// [END functions_storage_system_test]
