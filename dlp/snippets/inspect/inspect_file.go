// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package inspect contains example snippets using the DLP Inspect API.
package inspect

// [START dlp_inspect_file]
import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
	"log"
)

// inspectFile inspects the file at a given filePath, and prints results.
func inspectFile(projectID, filepath, fileType string) error {
	// projectID := "my-project-id";
	// filePath := "path/to/image.png";
	// fileType := "IMAGE"
	ctx := context.Background()

	// Initialize client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	// Construct the request to be processed by the client.
	// Set the item for the request to inspect.
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	var itemType dlppb.ByteContentItem_BytesType
	switch fileType {
	case "IMAGE":
		itemType = dlppb.ByteContentItem_IMAGE
	case "TEXT_UTF8":
		itemType = dlppb.ByteContentItem_TEXT_UTF8
	default:
		return fmt.Errorf("invalid ByteType for ByteContentItem: '%s'", fileType)
	}
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_ByteItem{
			ByteItem: &dlppb.ByteContentItem{
				Type: itemType,
				Data: data,
			},
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
		return fmt.Errorf("InspectContent: %v", err)
	}

	// Process the results.
	log.Printf("Findings: %d\n", len(resp.Result.Findings))
	for _, f := range resp.Result.Findings {
		log.Printf("\tQoute: %s\n", f.Quote)
		log.Printf("\tInfo type: %s\n", f.InfoType.Name)
		log.Printf("\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_file]
