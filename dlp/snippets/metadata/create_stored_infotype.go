// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package metadata

// [START dlp_create_stored_infotype]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// createStoredInfoType creates a custom stored info type based on your input data.
func createStoredInfoType(w io.Writer, projectID, outputPath string) error {
	// projectId := "my-project-id"
	// outputPath := "gs://" + "your-bucket-name" + "path/to/directory"

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// Specify the name you want to give the dictionary.
	displayName := "Github Usernames"

	// Specify a description of the dictionary.
	description := "Dictionary of GitHub usernames used in commits"

	// Specify the path to the location in a Cloud Storage
	// bucket to store the created dictionary.
	cloudStoragePath := &dlppb.CloudStoragePath{
		Path: outputPath,
	}

	// Specify your term list is stored in BigQuery.
	bigQueryField := &dlppb.BigQueryField{
		Table: &dlppb.BigQueryTable{
			ProjectId: "bigquery-public-data",
			DatasetId: "samples",
			TableId:   "github_nested",
		},
		Field: &dlppb.FieldId{
			Name: "actor",
		},
	}

	// Specify the configuration of the large custom dictionary.
	largeCustomDictionaryConfig := &dlppb.LargeCustomDictionaryConfig{
		OutputPath: cloudStoragePath,
		Source: &dlppb.LargeCustomDictionaryConfig_BigQueryField{
			BigQueryField: bigQueryField,
		},
	}

	// Specify the configuration for stored infoType.
	storedInfoTypeConfig := &dlppb.StoredInfoTypeConfig{
		DisplayName: displayName,
		Description: description,
		Type: &dlppb.StoredInfoTypeConfig_LargeCustomDictionary{
			LargeCustomDictionary: largeCustomDictionaryConfig,
		},
	}

	// Combine configurations into a request for the service.
	req := &dlppb.CreateStoredInfoTypeRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		Config:           storedInfoTypeConfig,
		StoredInfoTypeId: "github-usernames",
	}

	// Send the request and receive response from the service.
	resp, err := client.CreateStoredInfoType(ctx, req)
	if err != nil {
		return err
	}

	// Print the result.
	fmt.Fprintf(w, "output: %v", resp.Name)
	return nil

}

// [END dlp_create_stored_infotype]
