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

// [START data_catalog_quickstart]

// The datacatalog_quickstart application demonstrates how to define a tag
// template, populate values in the template, and attach a tag based on the
// template to a BigQuery table.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	datacatalog "cloud.google.com/go/datacatalog/apiv1"
	"cloud.google.com/go/datacatalog/apiv1/datacatalogpb"
	"github.com/googleapis/gax-go/v2"
)

func main() {
	projectID := flag.String("project_id", "", "Cloud Project ID, used for session creation.")
	location := flag.String("location", "us-central1", "data catalog region to use for the quickstart")
	table := flag.String("table", "myproject.mydataset.mytable", "bigquery table to tag in project.dataset.table format")

	flag.Parse()

	ctx := context.Background()
	client, err := datacatalog.NewClient(ctx)
	if err != nil {
		log.Fatalf("datacatalog.NewClient: %v", err)
	}
	defer client.Close()

	// Create the tag template.
	tmpl, err := createQuickstartTagTemplate(ctx, client, *projectID, *location)
	if err != nil {
		log.Fatalf("createQuickstartTagTemplate: %v", err)
	}
	fmt.Printf("Created tag template: %s\n", tmpl.GetName())

	// Convert a BigQuery resource identifier into the equivalent datacatalog
	// format.
	resource, err := convertBigQueryResourceRepresentation(*table)
	if err != nil {
		log.Fatalf("couldn't parse --table flag (%s): %v", *table, err)
	}

	// Lookup the entry metadata for the BQ table resource.
	entry, err := LookupEntry(ctx, client, &datacatalogpb.LookupEntryRequest{
		TargetName: &datacatalogpb.LookupEntryRequest_LinkedResource{
			LinkedResource: resource,
		},
	})
	if err != nil {
		log.Fatalf("LookupEntry: %v", err)
	}
	fmt.Printf("Successfully looked up table entry: %s\n", entry.GetName())

	// Create a tag based on the template, and apply it to the entry.
	tag, err := createQuickstartTag(ctx, client, "my-quickstart-tag", tmpl.GetName(), entry.GetName())
	if err != nil {
		log.Fatalf("couldn't create tag: %v", err)
	}
	fmt.Printf("Created tag: %s", tag.GetName())
}

// createQuickstartTagTemplate registers a tag template in datacatalog.
func createQuickstartTagTemplate(ctx context.Context, client *datacatalog.Client, projectID, location string) (*datacatalogpb.TagTemplate, error) {
	loc := fmt.Sprintf("projects/%s/locations/%s", projectID, location)

	// Define the tag template.
	template := &datacatalogpb.TagTemplate{
		DisplayName: "Quickstart Tag Template",
		Fields: map[string]*datacatalogpb.TagTemplateField{
			"source": {
				DisplayName: "Source of data asset",
				Type: &datacatalogpb.FieldType{
					TypeDecl: &datacatalogpb.FieldType_PrimitiveType_{
						PrimitiveType: datacatalogpb.FieldType_STRING,
					},
				},
			},
			"num_rows": {
				DisplayName: "Number of rows in data asset",
				Type: &datacatalogpb.FieldType{
					TypeDecl: &datacatalogpb.FieldType_PrimitiveType_{
						PrimitiveType: datacatalogpb.FieldType_DOUBLE,
					},
				},
			},
			"has_pii": {
				DisplayName: "Has PII",
				Type: &datacatalogpb.FieldType{
					TypeDecl: &datacatalogpb.FieldType_PrimitiveType_{
						PrimitiveType: datacatalogpb.FieldType_BOOL,
					},
				},
			},
			"pii_type": {
				DisplayName: "PII Type",
				Type: &datacatalogpb.FieldType{
					TypeDecl: &datacatalogpb.FieldType_EnumType_{
						EnumType: &datacatalogpb.FieldType_EnumType{
							AllowedValues: []*datacatalogpb.FieldType_EnumType_EnumValue{
								{DisplayName: "EMAIL"},
								{DisplayName: "SOCIAL SECURITY NUMBER"},
								{DisplayName: "NONE"},
							},
						},
					},
				},
			},
		},
	}

	//Construct the creation request using the template definition.
	req := &datacatalogpb.CreateTagTemplateRequest{
		Parent:        loc,
		TagTemplateId: "quickstart_tag_template",
		TagTemplate:   template,
	}

	// [END data_catalog_quickstart]
	// To aid testing, we add some uniqueness to the template ID.
	req.TagTemplateId = fmt.Sprintf("%s_%d", req.GetTagTemplateId(), time.Now().UnixNano())
	// [START data_catalog_quickstart]
	return client.CreateTagTemplate(ctx, req)

}

// createQuickstartTag populates tag values according to the template, and attaches
// the tag to the designeated entry.
func createQuickstartTag(ctx context.Context, client *datacatalog.Client, tagID, templateName, entryName string) (*datacatalogpb.Tag, error) {
	tag := &datacatalogpb.Tag{
		Name:     fmt.Sprintf("%s/tags/%s", entryName, tagID),
		Template: templateName,
		Fields: map[string]*datacatalogpb.TagField{
			"source": {
				Kind: &datacatalogpb.TagField_StringValue{StringValue: "Copied from tlc_yellow_trips_2018"},
			},
			"num_rows": {
				Kind: &datacatalogpb.TagField_DoubleValue{DoubleValue: 113496874},
			},
			"has_pii": {
				Kind: &datacatalogpb.TagField_BoolValue{BoolValue: false},
			},
			"pii_type": {
				Kind: &datacatalogpb.TagField_EnumValue_{
					EnumValue: &datacatalogpb.TagField_EnumValue{
						DisplayName: "NONE",
					},
				},
			},
		},
	}

	req := &datacatalogpb.CreateTagRequest{
		Parent: entryName,
		Tag:    tag,
	}
	return client.CreateTag(ctx, req)
}

// convertBigQueryResourceRepresentation converts a table identifier in standard sql form
// (project.datadata.table) into the representation used within data catalog.
func convertBigQueryResourceRepresentation(table string) (string, error) {
	parts := strings.Split(table, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("specified table string is not in expected project.dataset.table format: %s", table)
	}
	return fmt.Sprintf("//bigquery.googleapis.com/projects/%s/datasets/%s/tables/%s", parts[0], parts[1], parts[2]), nil
}

// LookupEntry provides a simply retry wrapper around the LookupEntry RPC.
//
// There's a potential propagation delay from when an entity is created until it appears in data catalog,
// so we wrap the lookup in a retry with a short context deadline to avoid unnecessary waiting for datacatalog
// to pick up new resources.
func LookupEntry(ctx context.Context, client *datacatalog.Client, req *datacatalogpb.LookupEntryRequest) (*datacatalogpb.Entry, error) {

	cCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// gax provides a basic backoff implementation for retries.
	bo := gax.Backoff{
		Initial: time.Second,
	}

	var entry *datacatalogpb.Entry
	var err error

	for {
		entry, err = client.LookupEntry(cCtx, req)
		if err != nil {
			if err = gax.Sleep(cCtx, bo.Pause()); err != nil {
				return nil, err
			}
			continue
		}
		return entry, err
	}
}

// [END data_catalog_quickstart]
