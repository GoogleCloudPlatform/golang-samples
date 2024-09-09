// Copyright 2024 Google LLC
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

package snippets

// [START auth_cloud_aws_provider]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/auth/credentials/externalaccount"
	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type customAwsProvider struct {
	awsRegion string
}

func (acp customAwsProvider) AwsRegion(ctx context.Context, opts *externalaccount.RequestOptions) (string, error) {
	return acp.awsRegion, nil
}

func (acp customAwsProvider) AwsSecurityCredentials(ctx context.Context, opts *externalaccount.RequestOptions) (*externalaccount.AwsSecurityCredentials, error) {
	// Load the AWS default config and retrieve the credentials.
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	val, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the AWS credentials into the external account libraries version of AWS security credentials.
	awsCredentials := &externalaccount.AwsSecurityCredentials{
		AccessKeyID:     val.AccessKeyID,
		SecretAccessKey: val.SecretAccessKey,
		SessionToken:    val.SessionToken,
	}

	return awsCredentials, nil
}

// authenticateExplicitWithAdc uses Application Default Credentials
// to print storage buckets.
func authenticateExplicitWithAwsProvider(w io.Writer) error {
	ctx := context.Background()

	// Set the active AWS region, i.e. "us-east-2"
	awsRegion := "us-east-2"
	awsCredentialsProvider := customAwsProvider{awsRegion: awsRegion}

	// Set the scopes for the credential.
	// If you are authenticating to a Cloud API, you can let the library include the default scope,
	// https://www.googleapis.com/auth/cloud-platform, because IAM is used to provide fine-grained
	// permissions for Cloud.
	// For more information on scopes to use,
	// see: https://developers.google.com/identity/protocols/oauth2/scopes
	scopes := []string{"https://www.googleapis.com/auth/cloud-platform"}

	// Set external account credential options.
	options := externalaccount.Options{
		Audience:                       "//iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/WORKLOAD_POOL/providers/WORKLOAD_PROVIDER", // Set the Workload Identity audience, replacing PROJECT_NUMBER, WORKLOAD_POOL, and WORKLOAD_PROVIDER.
		SubjectTokenType:               "urn:ietf:params:aws:token-type:aws4_request",                                                                                   // Set the AWS subject token type.
		AwsSecurityCredentialsProvider: awsCredentialsProvider,                                                                                                          // Set the AWS credentials provider.
		Scopes:                         scopes,                                                                                                                          // Set the scopes                                                                                                    // Set the AWS credentials provider.
	}
	credentials, err := externalaccount.NewCredentials(&options)
	if err != nil {
		return fmt.Errorf("failed to generate credentials: %w", err)
	}

	// Construct the Storage client.
	client, err := storage.NewClient(ctx, option.WithAuthCredentials(credentials))
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// Replace this with your GCP project ID.
	projectId := "PROJECT_ID"

	it := client.Buckets(ctx, projectId)
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Bucket: %v\n", bucketAttrs.Name)
	}

	fmt.Fprintf(w, "Listed all storage buckets.\n")

	return nil
}

// [END auth_cloud_aws_provider]
