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

// Package s3sdk lists GCS objects using the S3 SDK using interoperability mode.
package s3sdk

// [START storage_s3_sdk_list_objects]
import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func listGCSObjects(w io.Writer, bucketName string, googleAccessKeyID string, googleAccessKeySecret string) error {
	// bucketName := "your-gcs-bucket-name"
	// googleAccessKeyID := "Your Google Access Key ID"
	// googleAccessKeySecret := "Your Google Access Key Secret"

	// Create a new client and do the following:
	// 1. Change the endpoint URL to use the Google Cloud Storage XML API endpoint.
	// 2. Use Cloud Storage HMAC Credentials.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Endpoint:    aws.String("https://storage.googleapis.com"),
		Credentials: credentials.NewStaticCredentials(googleAccessKeyID, googleAccessKeySecret, ""),
	}))

	client := s3.New(sess)
	ctx := context.Background()

	result, err := client.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("ListObjectsWithContext: %w", err)
	}

	fmt.Fprintf(w, "Objects:")
	for _, o := range result.Contents {
		fmt.Fprintf(w, "%s\n", aws.StringValue(o.Key))
	}

	return nil
}

// [END storage_s3_sdk_list_objects]
