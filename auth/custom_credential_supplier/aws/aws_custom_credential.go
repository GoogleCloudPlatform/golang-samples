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

// [START auth_custom_credential_supplier_aws]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"cloud.google.com/go/auth/credentials/externalaccount"
	"github.com/aws/aws-sdk-go-v2/config"
)

// customAwsSupplier implements externalaccount.AwsSecurityCredentialsProvider
// using the official AWS SDK for Go v2.
type customAwsSupplier struct{}

// AwsRegion resolves the AWS region using the AWS SDK's default configuration chain.
// It prioritizes the AWS_REGION environment variable to match standard behavior.
func (s *customAwsSupplier) AwsRegion(ctx context.Context, opts *externalaccount.RequestOptions) (string, error) {
	// Explicitly check environment variable first.
	if region := os.Getenv("AWS_REGION"); region != "" {
		return region, nil
	}
	if region := os.Getenv("AWS_DEFAULT_REGION"); region != "" {
		return region, nil
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("AWS SDK failed to load config for region: %w", err)
	}

	if cfg.Region == "" {
		return "", fmt.Errorf("AWS region could not be resolved by SDK; ensure AWS_REGION is set")
	}

	return cfg.Region, nil
}

// AwsSecurityCredentials retrieves credentials using the AWS SDK's default provider chain.
func (s *customAwsSupplier) AwsSecurityCredentials(ctx context.Context, opts *externalaccount.RequestOptions) (*externalaccount.AwsSecurityCredentials, error) {
	// Load the default AWS configuration.
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("AWS SDK failed to load config for credentials: %w", err)
	}

	// Retrieve the actual credentials values.
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

// authenticateWithAwsCredentials demonstrates how to use a custom AWS credential supplier
// to authenticate with Google Cloud and verify access to a specific bucket.
//
// impersonationURL is optional. If provided, the credential will exchange the federated
// token for a Service Account token. If empty, the federated token is used directly.
func authenticateWithAwsCredentials(w io.Writer, bucketName, audience, impersonationURL string) error {
	// bucketName := "sample-bucket"
	// audience := "//iam.googleapis.com/projects/sample-project/locations/global/workloadIdentityPools/sample-pool/providers/sample-provider"
	// [Optional] impersonationURL := "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/myserviceaccount@iam.gserviceaccount.com:generateAccessToken"

	ctx := context.Background()

	// 1. Initialize Custom AWS Supplier
	supplier := &customAwsSupplier{}

	// 2. Configure the credentials options
	opts := &externalaccount.Options{
		Audience:                       audience,
		SubjectTokenType:               "urn:ietf:params:aws:token-type:aws4_request",
		ServiceAccountImpersonationURL: impersonationURL,
		AwsSecurityCredentialsProvider: supplier,
		Scopes:                         []string{"https://www.googleapis.com/auth/devstorage.read_write"},
	}

	// 3. Create the credentials object
	creds, err := externalaccount.NewCredentials(opts)
	if err != nil {
		return fmt.Errorf("externalaccount.NewCredentials: %w", err)
	}

	// 4. Authenticate and Check Bucket Existence via Raw HTTP
	bucketURL := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s", bucketName)
	fmt.Fprintf(w, "Request URL: %s\n", bucketURL)
	fmt.Fprintln(w, "Attempting to make authenticated request to Google Cloud Storage...")

	// Retrieve a valid token
	token, err := creds.Token(ctx)
	if err != nil {
		return fmt.Errorf("creds.Token: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", bucketURL, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Value)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	// Check for failure
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("bucket '%s' does not exist (404)", bucketName)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse and display success output
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("json.Decode: %w", err)
	}

	fmt.Fprintln(w, "\n--- SUCCESS! ---")
	fmt.Fprintln(w, "Successfully authenticated and retrieved bucket data:")
	prettyJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintln(w, string(prettyJSON))

	return nil
}

// [END auth_custom_credential_supplier_aws]

func main() {
	gcpAudience := os.Getenv("GCP_WORKLOAD_AUDIENCE")
	// Service account impersonation is optional
	saImpersonationURL := os.Getenv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL")
	gcsBucketName := os.Getenv("GCS_BUCKET_NAME")

	if gcpAudience == "" || gcsBucketName == "" {
		fmt.Fprintln(os.Stderr, "Missing required environment variables: GCP_WORKLOAD_AUDIENCE, GCS_BUCKET_NAME")
		os.Exit(1)
	}

	// Pass os.Stdout as the writer to print to console
	if err := authenticateWithAwsCredentials(os.Stdout, gcsBucketName, gcpAudience, saImpersonationURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
