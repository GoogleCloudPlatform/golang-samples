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

// Package v3 contains samples for Google Cloud Translation API v3.
package v3

// [START translate_v3_delete_glossary]
import (
	"context"
	"fmt"
	"io"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

// deleteGlossary deletes a glossary.
func deleteGlossary(w io.Writer, projectID string, location string, glossaryID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// glossaryID := "glossary-id"

	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return fmt.Errorf("NewTranslationClient: %w", err)
	}
	defer client.Close()

	req := &translatepb.DeleteGlossaryRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/glossaries/%s", projectID, location, glossaryID),
	}

	// The DeleteGlossary operation is async.
	op, err := client.DeleteGlossary(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteGlossary: %w", err)
	}
	fmt.Fprintf(w, "Processing operation name: %q\n", op.Name())

	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Deleted: %v\n", resp.GetName())

	return nil
}

// [END translate_v3_delete_glossary]
