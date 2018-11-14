package main

// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// This sample shows how to list Google Cloud Storage (GCS) buckets
//  using the AWS S3 SDK with the GCS interoperable XML API.
//
// GCS Credentials are passed in using the following environment variables:
//
//     * AWS_ACCESS_KEY_ID
//     * AWS_SECRET_ACCESS_KEY
//
// Learn how to get GCS interoperable credentials at
// https://cloud.google.com/storage/docs/migrating#keys.

// [START storage_s3_sdk_list_buckets]
import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func list_gcs_buckets() ([]*s3.Bucket, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("auto"),
		Endpoint: aws.String("https://storage.googleapis.com"),
	}))
	client := s3.New(sess)
	ctx := context.Background()

	result, err := client.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Buckets:")
	for _, b := range result.Buckets {
		fmt.Printf("%s\n", aws.StringValue(b.Name))
	}

	return result.Buckets, nil
}

// [START storage_s3_sdk_list_buckets]

func main() {
	list_gcs_buckets()
}
