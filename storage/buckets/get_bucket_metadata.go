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

// Sample buckets creates a bucket, lists buckets and deletes a bucket
// using the Google Storage API. More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.

package main

// [START storage_get_bucket_metadata]
import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
)

func getBucketMetadata(w io.Writer, client *storage.Client, bucketName string) (*storage.BucketAttrs, error) {
	// bucketName := "bucket-name"
	ctx := context.Background()

	// Initialize client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "BucketName: %v", attrs.Name)
	fmt.Fprintln(w, "Location: %v", attrs.Location)
	fmt.Fprintln(w, "LocationType: %v", attrs.LocationType)
	fmt.Fprintln(w, "StorageClass: %v", attrs.StorageClass)
	fmt.Fprintln(w, "TimeCreated: %v", attrs.Created)
	fmt.Fprintln(w, "Metageneration: %v", attrs.MetaGeneration)
	fmt.Fprintln(w, "PredefinedACL: %v", attrs.PredefinedACL)
	fmt.Fprintln(w, "DefaultKmsKeyName: %v", attrs.Encryption.DefaultKMSKeyName)
	fmt.Fprintln(w, "IndexPage: %v", attrs.Website.MainPageSuffix)
	fmt.Fprintln(w, "NotFoundPage: %v", attrs.Website.NotFoundPage)
	fmt.Fprintln(w, "DefaultEventBasedHold: %v", attrs.DefaultEventBasedHold)
	fmt.Fprintln(w, "RetentionEffectiveTime: %v", attrs.RetentionPolicy.EffectiveTime)
	fmt.Fprintln(w, "RetentionPeriod: %v", attrs.RetentionPolicy.RetentionPeriod)
	fmt.Fprintln(w, "RetentionPolicyIsLocked: %v", attrs.RetentionPolicy.IsLocked)
	fmt.Fprintln(w, "RequesterPays: %v", attrs.RequesterPays)
	fmt.Fprintln(w, "VersioningEnabled: %v", attrs.VersioningEnabled)
	fmt.Fprintln(w, "LogBucket: %v", attrs.Logging.LogBucket)
	fmt.Fprintln(w, "LogObjectPrefix: %v", attrs.Logging.LogObjectPrefix)
	fmt.Fprintln(w, "\n\n\nLabels:")
	for key, value := range attrs.Labels {
		fmt.Fprintln(w, "\t%v = %v", key, value)
	}

	return attrs, nil
}

// [END storage_get_bucket_metadata]
