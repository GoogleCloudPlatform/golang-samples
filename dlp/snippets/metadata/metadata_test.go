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

// Package metadata contains example snippets using the DLP info types API.
package metadata

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

const (
	termListFileName = "term_list.txt"
	filePathToUpload = "./testdata/term_list_storedInfotype.txt"
	bucket_prefix    = "test"
)

func deleteStoredInfoTypeAfterTest(t *testing.T, name string) error {
	t.Helper()
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &dlppb.DeleteStoredInfoTypeRequest{
		Name: name,
	}
	err = client.DeleteStoredInfoType(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func filesForUpdateStoredInfoType(t *testing.T, projectID, bucketName string) (string, string, error) {
	t.Helper()
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", "", err
	}
	defer client.Close()

	dirPath := "update-stored-infoType-data/"

	// Check if the directory already exists in the bucket.
	dirExists := false
	query := &storage.Query{Prefix: dirPath}
	it := client.Bucket(bucketName).Objects(ctx, query)
	_, err = it.Next()
	if err == nil {
		dirExists = true
	}

	// If the directory doesn't exist, create it.
	if !dirExists {
		obj := client.Bucket(bucketName).Object(dirPath)
		if _, err := obj.NewWriter(ctx).Write([]byte("")); err != nil {
			log.Fatalf("[INFO] [filesForUpdateStoredInfoType] Failed to create directory: %v", err)
		}
		log.Printf("[INFO] [filesForUpdateStoredInfoType] Directory '%s' created successfully in bucket '%s'.\n", dirPath, bucketName)
	} else {
		log.Printf("[INFO] [filesForUpdateStoredInfoType] Directory '%s' already exists in bucket '%s'.\n", dirPath, bucketName)
	}

	// file upload code

	// Open local file.
	file, err := os.ReadFile(filePathToUpload)
	if err != nil {
		log.Fatalf("[INFO] [filesForUpdateStoredInfoType] Failed to read file: %v", err)
	}

	// Get a reference to the bucket
	bucket := client.Bucket(bucketName)

	// Upload the file
	object := bucket.Object(termListFileName)
	writer := object.NewWriter(ctx)
	_, err = writer.Write(file)
	if err != nil {
		log.Fatalf("[INFO] [filesForUpdateStoredInfoType] Failed to write file: %v", err)
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}
	log.Printf("[INFO] [filesForUpdateStoredInfoType] File uploaded successfully: %v\n", termListFileName)

	// Check if the file exists in the bucket
	_, err = bucket.Object(termListFileName).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			log.Printf("[INFO] [filesForUpdateStoredInfoType] File %v does not exist in bucket %v\n", termListFileName, bucketName)
		} else {
			log.Fatalf("[INFO] [filesForUpdateStoredInfoType] Failed to check file existence: %v", err)
		}
	} else {
		log.Printf("[INFO] [filesForUpdateStoredInfoType] File %v exists in bucket %v\n", termListFileName, bucketName)
	}

	fileSetUrl := fmt.Sprint("gs://" + bucketName + "/" + termListFileName)
	gcsUri := fmt.Sprint("gs://" + bucketName)

	return fileSetUrl, gcsUri, err
}

func createStoredInfoTypeForTesting(t *testing.T, projectID, outputPath string) (string, error) {
	t.Helper()
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()
	u := uuid.New().String()[:8]
	displayName := "dlp-stored-info-type-for-test-" + u
	description := "Dictionary of GitHub usernames used in commits"

	cloudStoragePath := &dlppb.CloudStoragePath{
		Path: outputPath,
	}

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

	largeCustomDictionaryConfig := &dlppb.LargeCustomDictionaryConfig{
		OutputPath: cloudStoragePath,
		Source: &dlppb.LargeCustomDictionaryConfig_BigQueryField{
			BigQueryField: bigQueryField,
		},
	}

	storedInfoTypeConfig := &dlppb.StoredInfoTypeConfig{
		DisplayName: displayName,
		Description: description,
		Type: &dlppb.StoredInfoTypeConfig_LargeCustomDictionary{
			LargeCustomDictionary: largeCustomDictionaryConfig,
		},
	}

	req := &dlppb.CreateStoredInfoTypeRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		Config:           storedInfoTypeConfig,
		StoredInfoTypeId: "go-sample-test-stored-infoType-" + u,
	}
	resp, err := client.CreateStoredInfoType(ctx, req)
	if err != nil {
		return "nil", err
	}

	return resp.Name, nil
}
