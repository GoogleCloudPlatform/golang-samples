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

package regional_secretmanager

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"testing"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	grpcstatus "google.golang.org/grpc/status"
)

func testLocation(tb testing.TB) string {
	tb.Helper()

	v := os.Getenv("GOLANG_REGIONAL_SAMPLES_LOCATION")
	if v == "" {
		tb.Skip("testIamUser: missing GOLANG_REGIONAL_SAMPLES_LOCATION")
	}

	return v
}

func testRegionalClient(tb testing.TB) (*secretmanager.Client, context.Context) {
	tb.Helper()

	ctx := context.Background()

	locationId := testLocation(tb)

	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))

	if err != nil {
		tb.Fatalf("testRegionalClient: failed to create regional client: %v", err)
	}
	return client, ctx
}

func testResourceManagerTagBindingsClient(tb testing.TB, endpoint string) (*resourcemanager.TagBindingsClient, context.Context) {
	tb.Helper()
	ctx := context.Background()

	client, err := resourcemanager.NewTagBindingsClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		tb.Fatalf("testResourceManagerTagBindingsClient: failed to create client: %v", err)
	}
	return client, ctx
}

func testName(tb testing.TB) string {
	tb.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		tb.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
}

func testTopic(tb testing.TB) string {
	tb.Helper()

	v := os.Getenv("GOLANG_SAMPLES_TOPIC_NAME")
	if v == "" {
		tb.Skip("testTopic: missing GOLANG_SAMPLES_TOPIC_NAME")
	}

	return v
}

func testRegionalSecret(tb testing.TB, projectID string) (*secretmanagerpb.Secret, string) {
	tb.Helper()

	secretID := testName(tb)

	locationID := testLocation(tb)
	client, ctx := testRegionalClient(tb)
	secret, err := client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, locationID),
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Annotations: map[string]string{
				"annotationkey": "annotationvalue",
			},
			Labels: map[string]string{
				"labelkey": "labelvalue",
			},
		},
	})
	if err != nil {
		tb.Fatalf("testSecret: failed to create secret: %v", err)
	}

	return secret, secretID
}

func testCleanupRegionalSecret(tb testing.TB, name string) {
	tb.Helper()

	client, ctx := testRegionalClient(tb)

	if err := client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}); err != nil {
		if terr, ok := grpcstatus.FromError(err); !ok || terr.Code() != grpccodes.NotFound {
			tb.Fatalf("testCleanupSecret: failed to delete regional secret: %v", err)
		}
	}
}
func TestCreateRegionalSecretWithLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createRegionalSecretWithLabels"
	locationID := testLocation(t)

	var b bytes.Buffer
	if err := createRegionalSecretWithLabels(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID))

	if got, want := b.String(), "Created secret with labels:"; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithLabels: expected %q to contain %q", got, want)
	}
}

func TestDeleteRegionalSecretLabel(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretId := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := deleteRegionalSecretLabel(&b, tc.ProjectID, locationID, secretId); err != nil {
		t.Fatal(err)
	}

	client, ctx := testRegionalClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{}; reflect.DeepEqual(got, want) {
		t.Errorf("deleteRegionalSecretLabel: expected %q to be %q", got, want)
	}
}

func TestViewRegionalSecretLabels(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := viewRegionalSecretLabels(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found secret"; !strings.Contains(got, want) {
		t.Errorf("viewRegionalSecretLabels: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{"labelkey": "labelvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("viewRegionalSecretLabels: expected %q to be %q", got, want)
	}
}

func TestEditRegionalSecretLabel(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := editRegionalSecretLabel(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)

	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Labels, map[string]string{"labelkey": "updatedlabelvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("createUpdateSecretLabel: expected %q to be %q", got, want)
	}
}

func TestCreateRegionalSecretWithAnnotations(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createRegionalSecretWithAnnotations"
	locationID := testLocation(t)

	var b bytes.Buffer
	if err := createRegionalSecretWithAnnotations(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID))

	if got, want := b.String(), "Created secret with annotations:"; !strings.Contains(got, want) {
		t.Errorf("createSecret: expected %q to contain %q", got, want)
	}
}

func TestEditRegionalSecretAnnotation(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := editRegionalSecretAnnotation(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateSecret: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)

	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Annotations, map[string]string{"annotationkey": "updatedannotationvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("editSecretAnnotations: expected %q to be %q", got, want)
	}
}

func TestViewRegionalSecretAnnotations(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)

	var b bytes.Buffer
	if err := viewRegionalSecretAnnotations(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Found secret"; !strings.Contains(got, want) {
		t.Errorf("viewRegionalSecretAnnotations: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := s.Annotations, map[string]string{"annotationkey": "annotationvalue"}; !reflect.DeepEqual(got, want) {
		t.Errorf("viewRegionalSecretAnnotations: expected %q to be %q", got, want)
	}
}

func testResourceManagerTagsKeyClient(tb testing.TB) (*resourcemanager.TagKeysClient, context.Context) {
	tb.Helper()
	ctx := context.Background()

	client, err := resourcemanager.NewTagKeysClient(ctx)
	if err != nil {
		tb.Fatalf("testResourceManagerTagsKeyClient: failed to create client: %v", err)
	}
	return client, ctx

}

func testResourceManagerTagsValueClient(tb testing.TB) (*resourcemanager.TagValuesClient, context.Context) {
	tb.Helper()
	ctx := context.Background()

	client, err := resourcemanager.NewTagValuesClient(ctx)
	if err != nil {
		tb.Fatalf("testResourceManagerTagsValueClient: failed to create client: %v", err)
	}
	return client, ctx

}

func testCleanupTagKey(tb testing.TB, tagKeyName string) {
	tb.Helper()

	client, ctx := testResourceManagerTagsKeyClient(tb)

	tagKeyOperation, err := client.DeleteTagKey(ctx, &resourcemanagerpb.DeleteTagKeyRequest{
		Name: tagKeyName,
	})
	if err != nil {
		tb.Fatalf("testCleanupTagKey: failed to delete tagKey: %v", err)
		return
	}

	_, err = tagKeyOperation.Wait(ctx)
	if err != nil {
		tb.Fatalf("testCleanupTagKey: failed to delete TagKey after waiting for operation: %v", err)
	}
}

// Polling to clean up the tag value because, after deleting a secret, it takes some time for the tag value to become unbound.
func testCleanupTagValue(tb testing.TB, tagValueName string) {
	tb.Helper()

	client, ctx := testResourceManagerTagsValueClient(tb)

	maxPollingDuration := 10 * time.Minute
	initialDelay := 2 * time.Second
	maxBackoffDelay := 30 * time.Second

	startTime := time.Now()
	attempt := 0

	for time.Since(startTime) < maxPollingDuration {
		attempt++

		tagValueOperation, err := client.DeleteTagValue(ctx, &resourcemanagerpb.DeleteTagValueRequest{
			Name: tagValueName,
		})

		if err != nil {
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.NotFound {
				tb.Logf("Tag value %s already deleted (or never existed) after %v.", tagValueName, time.Since(startTime))
				return
			}

			if ok && s.Code() == codes.FailedPrecondition && strings.Contains(s.Message(), "attached to resources") {
				delay := initialDelay * time.Duration(1<<uint(attempt-1))
				if delay > maxBackoffDelay {
					delay = maxBackoffDelay
				}
				time.Sleep(delay)
				continue
			}

			tb.Errorf("testCleanupTagValue: failed to initiate delete for tag value %s due to unrecoverable error: %v", tagValueName, err)
			return
		}

		_, err = tagValueOperation.Wait(ctx)
		if err != nil {
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.NotFound {
				tb.Logf("Tag value %s was deleted during operation wait (or already gone).", tagValueName)
				return
			}
			if ok && s.Code() == codes.FailedPrecondition && strings.Contains(s.Message(), "attached to resources") {
				delay := initialDelay * time.Duration(1<<uint(attempt-1))
				if delay > maxBackoffDelay {
					delay = maxBackoffDelay
				}
				time.Sleep(delay)
				continue
			}

			tb.Errorf("testCleanupTagValue: failed to delete tag value %s after waiting for operation due to unrecoverable error: %v", tagValueName, err)
			return
		}

		tb.Logf("Successfully deleted tag value %s after %v (attempt %d).", tagValueName, time.Since(startTime), attempt)
		return
	}
	tb.Errorf("testCleanupTagValue: failed to delete tag value %s after %v (max duration reached). It might still be attached.", tagValueName, maxPollingDuration)
}

func testCreateTagKey(tb testing.TB, projectID string) *resourcemanagerpb.TagKey {
	tb.Helper()

	client, ctx := testResourceManagerTagsKeyClient(tb)
	parent := fmt.Sprintf("projects/%s", projectID)
	tagKeyName := testName(tb)
	tagKeyDescription := "creating tag key for secretmanager regional tags sample"

	tagKeyOperation, err := client.CreateTagKey(ctx, &resourcemanagerpb.CreateTagKeyRequest{
		TagKey: &resourcemanagerpb.TagKey{
			Parent:      parent,
			ShortName:   tagKeyName,
			Description: tagKeyDescription,
		},
	})
	if err != nil {
		tb.Fatalf("testCreateTagKey: failed to create tagKey: %v", err)
	}

	createdTagKey, err := tagKeyOperation.Wait(ctx)
	if err != nil {
		tb.Fatalf("testCreateTagKey: failed to create TagKey after waiting for operation: %v", err)
	}

	return createdTagKey
}

func testCreateTagValue(tb testing.TB, tagKeyId string) *resourcemanagerpb.TagValue {
	tb.Helper()

	client, ctx := testResourceManagerTagsValueClient(tb)
	tagValueName := testName(tb)
	tagKeyDescription := "creating TagValue for secretmanager regional tags sample"

	tagKeyOperation, err := client.CreateTagValue(ctx, &resourcemanagerpb.CreateTagValueRequest{
		TagValue: &resourcemanagerpb.TagValue{
			Parent:      tagKeyId,
			ShortName:   tagValueName,
			Description: tagKeyDescription,
		},
	})
	if err != nil {
		tb.Fatalf("testCreateTagValue: failed to create tagValue: %v", err)
	}

	createdTagValue, err := tagKeyOperation.Wait(ctx)
	if err != nil {
		tb.Fatalf("testCreateTagValue: failed to create TagValue after waiting for operation: %v", err)
	}

	return createdTagValue
}

func TestCreateRegionalSecretWithTags(t *testing.T) {
	tc := testutil.SystemTest(t)

	secretID := "createRegionalSecretWithTags"
	locationID := testLocation(t)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.Name)

	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.Name)

	var b bytes.Buffer
	if err := createRegionalSecretWithTags(&b, tc.ProjectID, locationID, secretID, tagKey.GetName(), tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID))

	if got, want := b.String(), "Created secret with tags:"; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithTags: expected %q to contain %q", got, want)
	}
}

func TestCreateRegionalSecretWithDelayedDestroy(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	ttl := 24 * time.Hour

	var b bytes.Buffer
	if err := createRegionalSecretWithDelayedDestroy(&b, tc.ProjectID, secretID, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithDelayedDestroy: expected %q to contain %q", got, want)
	}

	// Verify Delayed Destroy with GetSecret.
	client, ctx := testRegionalClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetVersionDestroyTtl() == nil {
		t.Fatal("GetSecret: VersionDestroyTtl is nil, expected non-nil")
	}
	if secret.GetVersionDestroyTtl().AsDuration() != ttl {
		t.Errorf("VersionDestroyTtl mismatch: got %v, want %v", secret.GetVersionDestroyTtl().AsDuration(), ttl)
	}
}

func TestUpdateRegionalSecretWithDelayedDestroy(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	var b bytes.Buffer
	if err := createRegionalSecretWithDelayedDestroy(&b, tc.ProjectID, secretID, locationID); err != nil {
		t.Fatal(err)
	}

	// Update TTL to 48 hours.
	updatedTTL := 48 * time.Hour
	if err := updateRegionalSecretWithDelayedDestroy(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateRegionalSecretWithDelayedDestroy: expected %q to contain %q", got, want)
	}

	// Verify Delayed Destroy with GetSecret.
	client, ctx := testRegionalClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetVersionDestroyTtl() == nil {
		t.Fatal("GetSecret: VersionDestroyTtl is nil, expected non-nil")
	}
	if secret.GetVersionDestroyTtl().AsDuration() != updatedTTL {
		t.Errorf("VersionDestroyTtl mismatch: got %v, want %v", secret.GetVersionDestroyTtl().AsDuration(), updatedTTL)
	}
}

func TestDeleteRegionalSecretVersionDestroyTTL(t *testing.T) {
	tc := testutil.SystemTest(t)
	locationID := testLocation(t)
	secretID := testName(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	var b bytes.Buffer
	if err := createRegionalSecretWithDelayedDestroy(&b, tc.ProjectID, secretID, locationID); err != nil {
		t.Fatal(err)
	}

	if err := deleteRegionalSecretVersionDestroyTTL(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "removed version_destroy_ttl"; !strings.Contains(got, want) {
		t.Errorf("deleteRegionalSecretVersionDestroyTTL: expected %q to contain %q", got, want)
	}

	// Verify TTL with GetSecret.
	client, ctx := testRegionalClient(t)
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetVersionDestroyTtl() != nil {
		t.Fatal("GetSecret: VersionDestroyTtl is not nil, expected nil")
	}
}

func TestCreateRegionalSecretWithRotation(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	topicName := testTopic(t)
	rotationPeriod := 24 * time.Hour

	var b bytes.Buffer
	if err := createRegionalSecretWithRotation(&b, tc.ProjectID, secretID, locationID, topicName); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithRotation: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetRotation() == nil {
		t.Fatal("GetSecret: Rotation is nil, expected non-nil")
	}
	if secret.GetRotation().GetRotationPeriod().AsDuration() != rotationPeriod {
		t.Errorf("RotationPeriod mismatch: got %v, want %v", secret.GetRotation().GetRotationPeriod().AsDuration(), rotationPeriod)
	}
	if secret.GetRotation().GetNextRotationTime() == nil {
		t.Fatal("GetSecret: NextRotationTime is nil, expected non-nil")
	}
}

func TestUpdateRegionalSecretRotationPeriod(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	topicName := testTopic(t)

	var b bytes.Buffer
	if err := createRegionalSecretWithRotation(&b, tc.ProjectID, secretID, locationID, topicName); err != nil {
		t.Fatal(err)
	}

	// Update rotation period.
	updatedRotationPeriod := 48 * time.Hour
	b.Reset()
	if err := updateRegionalSecretRotationPeriod(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}

	got := b.String()
	if !strings.Contains(got, secretID) {
		t.Errorf("updateRegionalSecretRotationPeriod: output %q did not contain secretId %q", got, secretID)
	}

	// Verify rotation period with GetSecret.
	client, ctx := testRegionalClient(t)
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetRotation() == nil {
		t.Fatal("GetSecret: Rotation is nil, expected non-nil")
	}
	if secret.GetRotation().GetRotationPeriod().AsDuration() != updatedRotationPeriod {
		t.Errorf("RotationPeriod mismatch: got %v, want %v", secret.GetRotation().GetRotationPeriod().AsDuration(), updatedRotationPeriod)
	}
}

func TestDeleteRegionalSecretRotation(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	topicName := testTopic(t)

	var b bytes.Buffer
	if err := createRegionalSecretWithRotation(&b, tc.ProjectID, secretID, locationID, topicName); err != nil {
		t.Fatal(err)
	}

	// Remove rotation.
	if err := deleteRegionalSecretRotation(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}

	got := b.String()
	if !strings.Contains(got, secretID) {
		t.Errorf("deleteRegionalSecretRotation: output %q did not contain secretId %q", got, secretID)
	}

	// Verify rotation is removed with GetSecret.
	client, ctx := testRegionalClient(t)
	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetRotation() != nil {
		t.Errorf("Rotation mismatch: got %v, want nil", secret.GetRotation())
	}
}

func TestCreateRegionalSecretWithTopic(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	topicName := testTopic(t)

	var b bytes.Buffer
	if err := createRegionalSecretWithTopic(&b, tc.ProjectID, secretID, locationID, topicName); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithTopic: expected %q to contain %q", got, want)
	}
	if got, want := b.String(), topicName; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithTopic: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if len(secret.GetTopics()) != 1 || secret.GetTopics()[0].GetName() != topicName {
		t.Errorf("Topics mismatch: got %v, want %s", secret.GetTopics(), topicName)
	}
}

func TestCreateRegionalSecretWithExpireTime(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	// Expire time in 1 hour.
	expire := time.Now().Add(time.Hour)

	var b bytes.Buffer
	if err := createRegionalSecretWithExpireTime(&b, tc.ProjectID, secretID, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created secret"; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithExpireTime: expected %q to contain %q", got, want)
	}

	// Verify expire time with GetSecret.
	client, ctx := testRegionalClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetExpireTime() == nil {
		t.Fatal("GetSecret: ExpireTime is nil, expected non-nil")
	}

	// Allow 1 second difference for precision.
	if diff := secret.GetExpireTime().AsTime().Unix() - expire.Unix(); diff > 1 || diff < -1 {
		t.Errorf("ExpireTime mismatch: got %v, want %v", secret.GetExpireTime().AsTime(), expire)
	}
}

func TestUpdateRegionalSecretExpiration(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	var b bytes.Buffer
	if err := createRegionalSecretWithExpireTime(&b, tc.ProjectID, secretID, locationID); err != nil {
		t.Fatal(err)
	}

	// Update expire time to 2 hours.
	newExpire := time.Now().Add(2 * time.Hour)
	if err := updateRegionalSecretExpiration(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Updated secret"; !strings.Contains(got, want) {
		t.Errorf("updateRegionalSecretExpiration: expected %q to contain %q", got, want)
	}

	// Verify expire time with GetSecret.
	client, ctx := testRegionalClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetExpireTime() == nil {
		t.Fatal("GetSecret: ExpireTime is nil, expected non-nil")
	}

	// Allow 1 second difference for precision.
	if diff := secret.GetExpireTime().AsTime().Unix() - newExpire.Unix(); diff > 1 || diff < -1 {
		t.Errorf("ExpireTime mismatch: got %v, want %v", secret.GetExpireTime().AsTime(), newExpire)
	}
}

func TestRemoveRegionalExpiration(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	var b bytes.Buffer
	if err := createRegionalSecretWithExpireTime(&b, tc.ProjectID, secretID, locationID); err != nil {
		t.Fatal(err)
	}

	// Remove expire time.
	b.Reset()
	if err := deleteRegionalExpiration(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Removed expiration"; !strings.Contains(got, want) {
		t.Errorf("deleteRegionalExpiration: expected %q to contain %q", got, want)
	}

	// Verify expire time is removed with GetSecret.
	client, ctx := testRegionalClient(t)

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secretName,
	})
	if err != nil {
		t.Fatalf("failed to get secret for verification: %v", err)
	}

	if secret.GetExpireTime() != nil {
		t.Errorf("GetSecret: ExpireTime is %v, expected nil", secret.GetExpireTime())
	}
}

func TestBindTagToRegionalSecret(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.GetName())
	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.GetName())

	var b bytes.Buffer
	if err := bindTagToRegionalSecret(&b, tc.ProjectID, secretID, locationID, tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	if got, want := b.String(), "Tag binding created"; !strings.Contains(got, want) {
		t.Errorf("bindTagToRegionalSecret: expected %q to contain %q", got, want)
	}

	// Verify binding exists with API
	ctx := context.Background()
	rmEndpoint := fmt.Sprintf("%s-cloudresourcemanager.googleapis.com:443", locationID)
	tagBindingsClient, ctx := testResourceManagerTagBindingsClient(t, rmEndpoint)
	defer tagBindingsClient.Close()

	parent := "//secretmanager.googleapis.com/" + secretName
	it := tagBindingsClient.ListTagBindings(ctx, &resourcemanagerpb.ListTagBindingsRequest{
		Parent: parent,
	})

	found := false
	for {
		binding, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Failed to list tag bindings for verification: %v", err)
		}
		if binding.TagValue == tagValue.GetName() && path.Base(binding.GetParent()) == path.Base(secretName) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Tag binding for %s with value %s not found after creation", secretName, tagValue.GetName())
	}
}

func TestListRegionalSecretTagBindings(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.GetName())
	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.GetName())

	// Create a secret and bind the tag to it for testing list.
	var b bytes.Buffer
	if err := createRegionalSecretWithTags(&b, tc.ProjectID, locationID, secretID, tagKey.GetName(), tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	b.Reset()
	if err := listRegionalSecretTagBindings(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), tagValue.GetName(); !strings.Contains(got, want) {
		t.Errorf("listRegionalSecretTagBindings: expected %q to contain %q", got, want)
	}
}

func TestDetachRegionalTag(t *testing.T) {
	tc := testutil.SystemTest(t)
	secretID := testName(t)
	locationID := testLocation(t)

	tagKey := testCreateTagKey(t, tc.ProjectID)
	defer testCleanupTagKey(t, tagKey.GetName())
	tagValue := testCreateTagValue(t, tagKey.GetName())
	defer testCleanupTagValue(t, tagValue.GetName())

	// Create a secret and bind the tag to it for testing detach.
	var b bytes.Buffer
	if err := createRegionalSecretWithTags(&b, tc.ProjectID, locationID, secretID, tagKey.GetName(), tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	secretName := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID)
	defer testCleanupRegionalSecret(t, secretName)

	b.Reset()
	if err := detachRegionalTag(&b, secretName, locationID, tagValue.GetName()); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Detached tag value"; !strings.Contains(got, want) {
		t.Errorf("detachRegionalTag: expected %q to contain %q", got, want)
	}

	b.Reset()
	if err := listRegionalSecretTagBindings(&b, secretName, locationID); err != nil {
		t.Fatal(err)
	}
	if got, dontwant := b.String(), tagValue.GetName(); strings.Contains(got, dontwant) {
		t.Errorf("listRegionalSecretTagBindings after detach: expected %q not to contain %q", got, dontwant)
	}
}

func TestDeleteRegionalSecretAnnotation(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, _ := testRegionalSecret(t, tc.ProjectID)
	defer testCleanupRegionalSecret(t, secret.Name)

	locationID := testLocation(t)
	annotationKey := "annotationkey"

	var b bytes.Buffer
	if err := deleteRegionalSecretAnnotation(&b, secret.Name, locationID); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Deleted annotation"; !strings.Contains(got, want) {
		t.Errorf("deleteSecretAnnotation: expected %q to contain %q", got, want)
	}

	client, ctx := testRegionalClient(t)
	s, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: secret.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := s.Annotations[annotationKey]; ok {
		t.Errorf("deleteRegionalSecretAnnotation: key %q still present after deletion", annotationKey)
	}
}
