// Copyright 2023 Google LLC
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
	"fmt"
	"strings"
	"testing"
	"time"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/documentai/apiv1/documentaipb"
	"github.com/google/uuid"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	testLocation := "us"
	testFilePath := "../testdata/invoice.pdf"
	testMimeType := "application/pdf"
	testProcessorType := "OCR_PROCESSOR"
	testPrefix := "golang-test"

	ctx := context.Background()

	endpoint := fmt.Sprintf("%s-documentai.googleapis.com:443", testLocation)
	client, err := documentai.NewDocumentProcessorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		t.Fatalf("error creating Document AI client: %v", err)
	}
	defer client.Close()

	testProcessorName, err := CreateTestProcessor(t, client, tc.ProjectID, testLocation, testProcessorType, testPrefix)

	testProcessorID := strings.Split(testProcessorName, "/")[5]

	stdOut, stdErr, err := m.Run(nil, 10*time.Minute,
		"--project_id", tc.ProjectID,
		"--location", testLocation,
		"--processor_id", testProcessorID,
		"--file_path", testFilePath,
		"--mime_type", testMimeType,
	)

	if err != nil {
		t.Errorf("stdout: %v", string(stdOut))
		t.Errorf("stderr: %v", string(stdErr))
		t.Errorf("execution failed: %v", err)
	}

	got := string(stdOut)

	if want := "Document Text:"; !strings.Contains(got, want) {
		t.Errorf("quickstart got %q, want %q", got, want)
	}
	if want := "Invoice"; !strings.Contains(got, want) {
		t.Errorf("quickstart got %q, want %q", got, want)
	}

	err = DeleteTestProcessor(t, client, testProcessorName)
	if err != nil {
		t.Errorf("Post-test cleanup failed: %v", err)
	}

}

// CreateTestProcessor creates a new processor with the given prefix
func CreateTestProcessor(t *testing.T, client *documentai.DocumentProcessorClient, projectID, location, processorType, prefix string) (string, error) {
	t.Helper()
	ctx := context.Background()

	processorName := UniqueProcessorName(prefix)

	req := &documentaipb.CreateProcessorRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Processor: &documentaipb.Processor{
			Type:        processorType,
			DisplayName: processorName,
		},
	}
	resp, err := client.CreateProcessor(ctx, req)
	if err != nil {
		t.Fatalf("error creating Document AI processor: %v", err)
	}

	return resp.GetName(), nil
}

// UniqueProcessorName returns a unique name with the test prefix
// Any process created with this prefix may be deleted by DeleteExpiredProcessors
func UniqueProcessorName(prefix string) string {
	return strings.Join([]string{prefix, uuid.New().String()}, "-")
}

func DeleteTestProcessor(t *testing.T, client *documentai.DocumentProcessorClient, processorName string) error {
	t.Helper()
	ctx := context.Background()

	req := &documentaipb.DeleteProcessorRequest{
		Name: processorName,
	}

	op, err := client.DeleteProcessor(ctx, req)
	if err != nil {
		t.Fatalf("error deleting Document AI processor: %v", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		t.Fatalf("error deleting Document AI processor: %v", err)
	}
	return nil
}
