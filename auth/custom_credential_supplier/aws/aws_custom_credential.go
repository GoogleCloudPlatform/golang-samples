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
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/auth/credentials/externalaccount"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"
)

// customAwsSupplier implements externalaccount.AwsSecurityCredentialsProvider
type customAwsSupplier struct{}

// AwsRegion resolves the region from the AWS SDK default config.
func (s *customAwsSupplier) AwsRegion(ctx context.Context, _ *externalaccount.RequestOptions) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}
	if cfg.Region == "" {
		return "", fmt.Errorf("AWS_REGION not set")
	}
	return cfg.Region, nil
}

// AwsSecurityCredentials retrieves credentials via the AWS SDK.
func (s *customAwsSupplier) AwsSecurityCredentials(ctx context.Context, _ *externalaccount.RequestOptions) (*externalaccount.AwsSecurityCredentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	c, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	return &externalaccount.AwsSecurityCredentials{
		AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey, SessionToken: c.SessionToken,
	}, nil
}

// authenticateWithAwsCredentials authenticates with Google Cloud using the custom supplier.
func authenticateWithAwsCredentials(w io.Writer, bucketName, audience, impersonationURL string) error {
	ctx := context.Background()

	// Initialize credentials with the custom supplier
	creds, err := externalaccount.NewCredentials(&externalaccount.Options{
		Audience:                       audience,
		SubjectTokenType:               "urn:ietf:params:aws:token-type:aws4_request",
		ServiceAccountImpersonationURL: impersonationURL,
		AwsSecurityCredentialsProvider: &customAwsSupplier{},
		Scopes:                         []string{"https://www.googleapis.com/auth/devstorage.read_only"},
	})
	if err != nil {
		return fmt.Errorf("NewCredentials: %w", err)
	}

	// Fetch the OAuth2 token
	token, err := creds.Token(ctx)
	if err != nil {
		return fmt.Errorf("creds.Token: %w", err)
	}

	// Verify access by making a raw HTTP request to the GCS API
	url := fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s", bucketName)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token.Value)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed (%d): %s", resp.StatusCode, string(b))
	}

	fmt.Fprintf(w, "Success! Authenticated and accessed bucket: %s\n", bucketName)
	return nil
}

// [END auth_custom_credential_supplier_aws]

// Note: The code below handles local configuration loading and is not part of the core sample.

func main() {
	loadConfigFromFile()

	audience := os.Getenv("GCP_WORKLOAD_AUDIENCE")
	bucket := os.Getenv("GCS_BUCKET_NAME")
	impersonationURL := os.Getenv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL")

	if audience == "" || bucket == "" {
		fmt.Fprintln(os.Stderr, "Missing required configuration: GCP_WORKLOAD_AUDIENCE, GCS_BUCKET_NAME")
		os.Exit(1)
	}

	if err := authenticateWithAwsCredentials(os.Stdout, bucket, audience, impersonationURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func loadConfigFromFile() {
	file, err := os.Open("custom-credentials-aws-secrets.json")
	if err != nil {
		return // File missing, rely on env vars
	}
	defer file.Close()

	var s struct {
		Audience  string `json:"gcp_workload_audience"`
		Bucket    string `json:"gcs_bucket_name"`
		ImpURL    string `json:"gcp_service_account_impersonation_url"`
		AwsID     string `json:"aws_access_key_id"`
		AwsSecret string `json:"aws_secret_access_key"`
		AwsRegion string `json:"aws_region"`
	}
	if json.NewDecoder(file).Decode(&s) == nil {
		setEnv("AWS_ACCESS_KEY_ID", s.AwsID)
		setEnv("AWS_SECRET_ACCESS_KEY", s.AwsSecret)
		setEnv("AWS_REGION", s.AwsRegion)
		setEnv("GCP_WORKLOAD_AUDIENCE", s.Audience)
		setEnv("GCS_BUCKET_NAME", s.Bucket)
		setEnv("GCP_SERVICE_ACCOUNT_IMPERSONATION_URL", s.ImpURL)
	}
}

func setEnv(k, v string) {
	if v != "" {
		os.Setenv(k, v)
	}
}
