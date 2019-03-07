package main

// Copyright 2019 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START storage_s3_sdk_list_buckets]
import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func list_gcs_buckets(googleAccessKeyID string, googleAccessKeySecret string) ([]*s3.Bucket, error) {
	// Create a new client and do the following:
	// 1. Change the endpoint URL to use the Google Cloud Storage XML API endpoint.
    // 2. Use Cloud Storage HMAC Credentials.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("auto"),
		Endpoint: aws.String("https://storage.googleapis.com"),
		Credentials: credentials.NewStaticCredentials(googleAccessKeyID, googleAccessKeySecret, ""),
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

// [END storage_s3_sdk_list_buckets]

func main() {
	var googleAccessKeyID string
	var googleAccessKeySecret string

	flag.StringVar(&googleAccessKeyID, "googleAccessKeyID", "", "Your Cloud Storage HMAC Access Key ID.")
	flag.StringVar(&googleAccessKeySecret, "googleAccessKeySecret", "", "Your Cloud Storage HMAC Access Key Secret.")
	flag.Parse()

	list_gcs_buckets(googleAccessKeyID, googleAccessKeySecret)
}
