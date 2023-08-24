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

// [START dlp_update_stored_infotype]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateStoredInfoType uses the Data Loss Prevention API to update stored infoType
// detector by changing the source term list from one stored in Bigquery
// to one stored in Cloud Storage.
func updateStoredInfoType(w io.Writer, projectID, gcsUri, fileSetUrl, infoTypeId string) error {
	// projectId := "your-project-id"
	// gcsUri := "gs://" + "your-bucket-name" + "/path/to/your/file.txt"
	// fileSetUrl := "your-cloud-storage-file-set"
	// infoTypeId := "your-stored-info-type-id"

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

	// Set path in Cloud Storage.
	cloudStoragePath := &dlppb.CloudStoragePath{
		Path: gcsUri,
	}
	cloudStorageFileSet := &dlppb.CloudStorageFileSet{
		Url: fileSetUrl,
	}

	// Configuration for a custom dictionary created from a data source of any size
	largeCustomDictionaryConfig := &dlppb.LargeCustomDictionaryConfig{
		OutputPath: cloudStoragePath,
		Source: &dlppb.LargeCustomDictionaryConfig_CloudStorageFileSet{
			CloudStorageFileSet: cloudStorageFileSet,
		},
	}

	// Set configuration for stored infoTypes.
	storedInfoTypeConfig := &dlppb.StoredInfoTypeConfig{
		Type: &dlppb.StoredInfoTypeConfig_LargeCustomDictionary{
			LargeCustomDictionary: largeCustomDictionaryConfig,
		},
	}

	// Set mask to control which fields get updated.
	fieldMask := &fieldmaskpb.FieldMask{
		Paths: []string{"large_custom_dictionary.cloud_storage_file_set.url"},
	}
	// Construct the job creation request to be sent by the client.
	req := &dlppb.UpdateStoredInfoTypeRequest{
		Name:       fmt.Sprint("projects/" + projectID + "/storedInfoTypes/" + infoTypeId),
		Config:     storedInfoTypeConfig,
		UpdateMask: fieldMask,
	}

	// Use the client to send the API request.
	resp, err := client.UpdateStoredInfoType(ctx, req)
	if err != nil {
		return err
	}

	// Print the result.
	fmt.Fprintf(w, "output: %v", resp.Name)
	return nil
}

// [END dlp_update_stored_infotype]
