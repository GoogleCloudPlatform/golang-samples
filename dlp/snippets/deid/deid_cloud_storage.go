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
package deid

// [START dlp_deidentify_cloud_storage]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

func deidentifyCloudStorage(w io.Writer, projectID, gcsUri, tableId, datasetId, outputDirectory, deidentifyTemplateId, structuredDeidentifyTemplateId, imageRedactTemplateId string) error {
	// projectId := "my-project-id"
	// gcsUri := "gs://" + "your-bucket-name" + "/path/to/your/file.txt"
	// tableId := "your-bigquery-table-id"
	// datasetId := "your-bigquery-dataset-id"
	// outputDirectory := "your-output-directory"
	// deidentifyTemplateId := "your-deidentify-template-id"
	// structuredDeidentifyTemplateId := "your-structured-deidentify-template-id"
	// imageRedactTemplateId := "your-image-redact-template-id"

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
	cloudStorageOptions := &dlppb.CloudStorageOptions{
		FileSet: &dlppb.CloudStorageOptions_FileSet{
			Url: gcsUri,
		},
	}

	// Define the storage config options for cloud storage options.
	storageConfig := &dlppb.StorageConfig{
		Type: &dlppb.StorageConfig_CloudStorageOptions{
			CloudStorageOptions: cloudStorageOptions,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	infoTypes := []*dlppb.InfoType{
		{Name: "PERSON_NAME"},
		{Name: "EMAIL_ADDRESS"},
	}

	// inspectConfig holds the configuration settings for data inspection and analysis
	// within the context of the Google Cloud Data Loss Prevention (DLP) API.
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes:    infoTypes,
		IncludeQuote: true,
	}

	// Types of files to include for de-identification.
	fileTypesToTransform := []dlppb.FileType{
		dlppb.FileType_CSV,
		dlppb.FileType_IMAGE,
		dlppb.FileType_TEXT_FILE,
	}

	// Specify the BigQuery table to be inspected.
	table := &dlppb.BigQueryTable{
		ProjectId: projectID,
		DatasetId: datasetId,
		TableId:   tableId,
	}

	// transformationDetailsStorageConfig holds configuration settings for storing transformation
	// details in the context of the Google Cloud Data Loss Prevention (DLP) API.
	transformationDetailsStorageConfig := &dlppb.TransformationDetailsStorageConfig{
		Type: &dlppb.TransformationDetailsStorageConfig_Table{
			Table: table,
		},
	}

	transformationConfig := &dlppb.TransformationConfig{
		DeidentifyTemplate:           deidentifyTemplateId,
		ImageRedactTemplate:          imageRedactTemplateId,
		StructuredDeidentifyTemplate: structuredDeidentifyTemplateId,
	}

	// Action to execute on the completion of a job.
	deidentify := &dlppb.Action_Deidentify{
		TransformationConfig:               transformationConfig,
		TransformationDetailsStorageConfig: transformationDetailsStorageConfig,
		Output: &dlppb.Action_Deidentify_CloudStorageOutput{
			CloudStorageOutput: outputDirectory,
		},
		FileTypesToTransform: fileTypesToTransform,
	}

	action := &dlppb.Action{
		Action: &dlppb.Action_Deidentify_{
			Deidentify: deidentify,
		},
	}

	// Configure the inspection job we want the service to perform.
	inspectJobConfig := &dlppb.InspectJobConfig{
		StorageConfig: storageConfig,
		InspectConfig: inspectConfig,
		Actions: []*dlppb.Action{
			action,
		},
	}

	// Construct the job creation request to be sent by the client.
	req := &dlppb.CreateDlpJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		Job: &dlppb.CreateDlpJobRequest_InspectJob{
			InspectJob: inspectJobConfig,
		},
	}

	// Send the request.
	resp, err := client.CreateDlpJob(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "error after resp: %v", err)
		return err
	}

	// Print the results.
	fmt.Fprint(w, "Job created successfully: ", resp.Name)
	return nil

}

// [END dlp_deidentify_cloud_storage]
