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

package buckets

// [START storage_get_bucket_metadata]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// getBucketMetadata gets the bucket metadata.
func getBucketMetadata(w io.Writer, bucketName string) (*storage.BucketAttrs, error) {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Bucket(%q).Attrs: %w", bucketName, err)
	}
	fmt.Fprintf(w, "BucketName: %v\n", attrs.Name)
	fmt.Fprintf(w, "Location: %v\n", attrs.Location)
	fmt.Fprintf(w, "LocationType: %v\n", attrs.LocationType)
	fmt.Fprintf(w, "StorageClass: %v\n", attrs.StorageClass)
	fmt.Fprintf(w, "Turbo replication (RPO): %v\n", attrs.RPO)
	fmt.Fprintf(w, "TimeCreated: %v\n", attrs.Created)
	fmt.Fprintf(w, "Metageneration: %v\n", attrs.MetaGeneration)
	fmt.Fprintf(w, "PredefinedACL: %v\n", attrs.PredefinedACL)
	if attrs.Encryption != nil {
		fmt.Fprintf(w, "DefaultKmsKeyName: %v\n", attrs.Encryption.DefaultKMSKeyName)
	}
	if attrs.Website != nil {
		fmt.Fprintf(w, "IndexPage: %v\n", attrs.Website.MainPageSuffix)
		fmt.Fprintf(w, "NotFoundPage: %v\n", attrs.Website.NotFoundPage)
	}
	fmt.Fprintf(w, "DefaultEventBasedHold: %v\n", attrs.DefaultEventBasedHold)
	if attrs.RetentionPolicy != nil {
		fmt.Fprintf(w, "RetentionEffectiveTime: %v\n", attrs.RetentionPolicy.EffectiveTime)
		fmt.Fprintf(w, "RetentionPeriod: %v\n", attrs.RetentionPolicy.RetentionPeriod)
		fmt.Fprintf(w, "RetentionPolicyIsLocked: %v\n", attrs.RetentionPolicy.IsLocked)
	}
	fmt.Fprintf(w, "ObjectRetentionMode: %v\n", attrs.ObjectRetentionMode)
	fmt.Fprintf(w, "RequesterPays: %v\n", attrs.RequesterPays)
	fmt.Fprintf(w, "VersioningEnabled: %v\n", attrs.VersioningEnabled)
	if attrs.Logging != nil {
		fmt.Fprintf(w, "LogBucket: %v\n", attrs.Logging.LogBucket)
		fmt.Fprintf(w, "LogObjectPrefix: %v\n", attrs.Logging.LogObjectPrefix)
	}
	if attrs.CORS != nil {
		fmt.Fprintln(w, "CORS:")
		for _, v := range attrs.CORS {
			fmt.Fprintf(w, "\tMaxAge: %v\n", v.MaxAge)
			fmt.Fprintf(w, "\tMethods: %v\n", v.Methods)
			fmt.Fprintf(w, "\tOrigins: %v\n", v.Origins)
			fmt.Fprintf(w, "\tResponseHeaders: %v\n", v.ResponseHeaders)
		}
	}
	if attrs.Labels != nil {
		fmt.Fprintf(w, "\n\n\nLabels:")
		for key, value := range attrs.Labels {
			fmt.Fprintf(w, "\t%v = %v\n", key, value)
		}
	}
	return attrs, nil
}

// [END storage_get_bucket_metadata]
