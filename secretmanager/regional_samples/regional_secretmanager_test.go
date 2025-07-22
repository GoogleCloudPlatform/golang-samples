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
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func testName(tb testing.TB) string {
	tb.Helper()

	u, err := uuid.NewV4()
	if err != nil {
		tb.Fatalf("testName: failed to generate uuid: %v", err)
	}
	return u.String()
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
		if terr, ok := status.FromError(err); !ok || terr.Code() != codes.NotFound {
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
	tagKeyName := "sm_secret_regional_tag_key_test"
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
	tagValueName := "sm_secret_regional_tag_value_test"
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

	secretID := testName(t)
	locationID := testLocation(t)

	var b bytes.Buffer
	if err := createRegionalSecretWithDelayedDestroy(&b, tc.ProjectID, locationID, secretID, 86400); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, fmt.Sprintf("projects/%s/locations/%s/secrets/%s", tc.ProjectID, locationID, secretID))

	if got, want := b.String(), "Created secret with version destroy ttl:"; !strings.Contains(got, want) {
		t.Errorf("createRegionalSecretWithDelayedDestroy: expected %q to contain %q", got, want)
	}
}

func TestUpdateRegionalSecretWithDelayedDestroy(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	locationID := testLocation(t)

	var b bytes.Buffer
	if err := updateRegionalSecretWithDelayedDestroy(&b, tc.ProjectID, locationID, secretID, 86400); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, secret.Name)

	if got, want := b.String(), "Updated secret:"; !strings.Contains(got, want) {
		t.Errorf("updateRegionalSecretWithDelayedDestroy: expected %q to contain %q", got, want)
	}
}

func TestDisableRegionalSecretDelayedDestroy(t *testing.T) {
	tc := testutil.SystemTest(t)

	secret, secretID := testRegionalSecret(t, tc.ProjectID)
	locationID := testLocation(t)

	var b bytes.Buffer
	if err := disableRegionalSecretDelayedDestroy(&b, tc.ProjectID, locationID, secretID); err != nil {
		t.Fatal(err)
	}
	defer testCleanupRegionalSecret(t, secret.Name)

	if got, want := b.String(), "Updated secret:"; !strings.Contains(got, want) {
		t.Errorf("disableRegionalSecretDelayedDestroy: expected %q to contain %q", got, want)
	}
}
