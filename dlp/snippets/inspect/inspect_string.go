// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package inspect contains example snippets using the DLP Inspect API.
package inspect

// [START dlp_inspect_string]
import (
	"context"
	"log"

	"cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// inspectString inspects the a given string, and prints results.
func inspectString(projectID, textToInspect string) error {
	// projectID := "my-project-id"
	// textToInspect := "My name is Gary and my email is gary@example.com"
	ctx := context.Background()

	// Initialize client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Construct the request to be processed by the client.
	// Set the item for the request to inspect.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: textToInspect,
		},
	}

	// Set the inspection configuration for the request.
	config := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			{Name: "PHONE_NUMBER"},
			{Name: "EMAIL_ADDRESS"},
			{Name: "CREDIT_CARD_NUMBER"},
		},
		IncludeQuote: true,
	}

	// Create and send the request.
	req := &dlppb.InspectContentRequest{
		Parent:        "projects/" + projectID,
		Item:          item,
		InspectConfig: config,
	}
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		return err
	}

	// Process the results.
	result := resp.Result
	log.Printf("Findings: %d\n", len(result.Findings))
	for _, f := range result.Findings {
		log.Printf("\tQoute: %s\n", f.Quote)
		log.Printf("\tInfo type: %s\n", f.InfoType.Name)
		log.Printf("\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_string]
