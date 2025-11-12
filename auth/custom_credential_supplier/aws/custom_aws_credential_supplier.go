// Copyright 2025 Google LLC
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

package main

// This sample uses cloud.google.com/go/auth, not the older golang.org/x/oauth2/google.
// [START auth_custom_credential_supplier_aws]

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/auth/credentials/externalaccount"
	"cloud.google.com/go/auth/oauth2adapt"
	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// CustomAwsSupplier implements externalaccount.AwsSecurityCredentialsProvider
// using the official AWS SDK for Go v2.
type CustomAwsSupplier struct{}

// AwsRegion resolves the AWS region using the AWS SDK's default configuration chain.
// In EKS, this typically picks up the AWS_REGION environment variable automatically.
func (s *CustomAwsSupplier) AwsRegion(ctx context.Context, opts *externalaccount.RequestOptions) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("AWS SDK failed to load config for region: %w", err)
	}

	if cfg.Region == "" {
		// Fallback: In some minimal EKS setups, AWS_REGION might not be set by default
		// even if credentials work. You might want to hardcode a default here if acceptable.
		return "", fmt.Errorf("AWS region could not be resolved by SDK; ensure AWS_REGION is set")
	}

	return cfg.Region, nil
}

// AwsSecurityCredentials retrieves credentials using the AWS SDK's default provider chain.
// This supports EKS IRSA (IAM Roles for Service Accounts), EC2 IMDS, environment variables, etc.
func (s *CustomAwsSupplier) AwsSecurityCredentials(ctx context.Context, opts *externalaccount.RequestOptions) (*externalaccount.AwsSecurityCredentials, error) {
	// Load the default AWS configuration.
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("AWS SDK failed to load config for credentials: %w", err)
	}

	// Retrieve the actual credentials values.
	// The SDK handles caching and refreshing these automatically.
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("AWS SDK failed to retrieve credentials: %w", err)
	}

	return &externalaccount.AwsSecurityCredentials{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
	}, nil
}

func main() {
	ctx := context.Background()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Read GCP configuration from environment
	gcpAudience := os.Getenv("GCP_WORKLOAD_AUDIENCE")
	saImpersonationURL := os.Getenv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL")
	targetProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if gcpAudience == "" || targetProjectID == "" {
		log.Fatal("Missing required environment variables: GCP_WORKLOAD_AUDIENCE, GOOGLE_CLOUD_PROJECT")
	}

	log.Println("Initializing Custom AWS Supplier with AWS SDK...")
	// 1. Instantiate the custom supplier.
	customSupplier := &CustomAwsSupplier{}

	// 2. Configure the credentials options.
	opts := &externalaccount.Options{
		Audience:                       gcpAudience,
		SubjectTokenType:               "urn:ietf:params:aws:token-type:aws4_request",
		ServiceAccountImpersonationURL: saImpersonationURL,
		AwsSecurityCredentialsProvider: customSupplier,
		Scopes:                         []string{"https://www.googleapis.com/auth/cloud-platform"},
	}

	// 3. Create the credentials.
	creds, err := externalaccount.NewCredentials(opts)
	if err != nil {
		log.Fatalf("Failed to create external credentials: %v", err)
	}

	// Adapt to the older interface required by current Google Cloud clients
	oauth2Creds := oauth2adapt.Oauth2CredentialsFromAuthCredentials(creds)

	// 4. Use the credentials with the Storage client.
	log.Println("Creating Storage client...")
	storageClient, err := storage.NewClient(ctx, option.WithCredentials(oauth2Creds))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	// Example: List buckets to verify authentication.
	log.Printf("Attempting to list buckets in project: %s\n", targetProjectID)
	it := storageClient.Buckets(ctx, targetProjectID)
	count := 0
	for {
		bkt, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list buckets: %v", err)
		}
		fmt.Printf(" - %s\n", bkt.Name)
		count++
		if count >= 10 {
			fmt.Println("... (stopping after 10)")
			break
		}
	}
	log.Println("Successfully listed buckets.")
}

// [END auth_custom_credential_supplier_aws]
