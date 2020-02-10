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

// Package s3sdk lists GCS buckets using the S3 SDK using interoperability mode.
package s3sdk

// [START storage_s3_sdk_list_buckets]
import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func listGCSBuckets(w io.Writer, googleAccessKeyID string, googleAccessKeySecret string) ([]*s3.Bucket, error) {
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

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	result, err := client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Buckets:")
	for _, b := range result.Buckets {
		fmt.Fprintf(w, "%s\n", aws.StringValue(b.Name))
	}

	return result.Buckets, nil
}

// [END storage_s3_sdk_list_buckets]
