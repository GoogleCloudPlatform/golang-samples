// Copyright 2023 Google LLC
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

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/googleapi"
)

// Global variables used in testing
var location string
var projectId string
var r *rand.Rand
var buf bytes.Buffer

// Setup for all tests
func setupTests(t *testing.T) {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	location = "us-central1"
	tc := testutil.SystemTest(t)
	projectId = tc.ProjectID
}

// Setup and teardown functions for CaPoolTests
func setupCaPool(t *testing.T) (string, func(t *testing.T)) {
	caPoolId := fmt.Sprintf("test-ca-pool-%v-%v", time.Now().Format("2006-01-02"), r.Int())

	if err := createCaPool(&buf, projectId, location, caPoolId); err != nil {
		t.Fatal("setupCaPool got err:", err)
	}

	// Return a function to teardown the test
	return caPoolId, func(t *testing.T) {
		if err := deleteCaPool(&buf, projectId, location, caPoolId); err != nil {
			var gerr *googleapi.Error
			if errors.As(err, &gerr) {
				if gerr.Code == 404 {
					t.Log("setupCaPool teardown - skipped CA Pool deletion (not found)")
				} else {
					t.Errorf("setupCaPool teardown got err: %v", err)
				}
			}
		}
	}
}

// Setup and teardown functions for CaTests
func setupCa(t *testing.T, caPoolId string) (string, func(t *testing.T)) {
	caId := fmt.Sprintf("test-ca-%v-%v", time.Now().Format("2006-01-02"), r.Int())
	caCommonName := fmt.Sprintf("CN - %s", caId)
	org := "ORGANIZATION"
	caDuration := int64(2592000) // 30 days

	if err := createCa(&buf, projectId, location, caPoolId, caId, caCommonName, org, caDuration); err != nil {
		t.Fatal("setupCa got err:", err)
	}

	// Return a function to teardown the test
	return caId, func(t *testing.T) {
		if err := deleteCaPerm(projectId, location, caPoolId, caId); err != nil {
			var gerr *googleapi.Error
			if errors.As(err, &gerr) {
				if gerr.Code == 404 {
					t.Log("setupCa teardown - skipped CA Pool deletion (not found)")
				} else {
					t.Errorf("setupCa teardown got err: %v", err)
				}
			}
		}
	}
}

// Helper function to permanently remove CAs without 30d grace period
func deleteCaPerm(projectID string, location string, caPoolId string, caId string) error {
	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return err
	}
	defer caClient.Close()

	fullCaName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s/certificateAuthorities/%s",
		projectID, location, caPoolId, caId)

	req := &privatecapb.DeleteCertificateAuthorityRequest{
		Name:                     fullCaName,
		IgnoreActiveCertificates: true,
		IgnoreDependentResources: true,
		SkipGracePeriod:          true,
	}

	op, err := caClient.DeleteCertificateAuthority(ctx, req)
	if err != nil {
		return err
	}

	if _, err = op.Wait(ctx); err != nil {
		return err
	}

	return nil
}

func TestCreateCaPool(t *testing.T) {
	setupTests(t)

	t.Run("createCaPool", func(t *testing.T) {
		caPoolId := fmt.Sprintf("test-ca-pool-%v-%v", time.Now().Format("2006-01-02"), r.Int())

		buf.Reset()
		if err := createCaPool(&buf, projectId, location, caPoolId); err != nil {
			t.Fatal("createCaPool got err:", err)
		}

		expectedResult := "CA Pool created"
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("createCaPool got %q, want %q", got, expectedResult)
		}

		if err := deleteCaPool(&buf, projectId, location, caPoolId); err != nil {
			t.Fatal("createCaPool teardown (deleteCaPool) got err:", err)
		}
	})

	t.Run("deleteCaPool", func(t *testing.T) {
		caPoolId, teardownCaPoolTests := setupCaPool(t)
		defer teardownCaPoolTests(t)

		buf.Reset()
		if err := deleteCaPool(&buf, projectId, location, caPoolId); err != nil {
			t.Fatal("deleteCaPool got err:", err)
		}

		expectedResult := "CA Pool deleted"
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("deleteCaPool got %q, want %q", got, expectedResult)
		}
	})
}

func TestCreateCa(t *testing.T) {
	setupTests(t)
	caPoolId, teardownCaPool := setupCaPool(t)
	defer teardownCaPool(t)

	t.Run("createCa", func(t *testing.T) {
		caId := fmt.Sprintf("test-ca-%v-%v", time.Now().Format("2006-01-02"), r.Int())
		caCommonName := fmt.Sprintf("CN - %s", caId)
		org := "ORGANIZATION"
		caDuration := int64(2592000) // 30 days

		buf.Reset()
		if err := createCa(&buf, projectId, location, caPoolId, caId, caCommonName, org, caDuration); err != nil {
			t.Fatal("createCa got err:", err)
		}

		expectedResult := fmt.Sprintf("CA %s created", caId)
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("createCa got %q, want %q", got, expectedResult)
		}

		if err := deleteCaPerm(projectId, location, caPoolId, caId); err != nil {
			t.Fatal("createCa teardown (deleteCaPerm) got err:", err)
		}
	})

	t.Run("deleteCa", func(t *testing.T) {
		caId, _ := setupCa(t, caPoolId)

		buf.Reset()
		if err := deleteCa(&buf, projectId, location, caPoolId, caId); err != nil {
			t.Fatal("deleteCa got err:", err)
		}

		expectedResult := fmt.Sprintf("Successfully deleted Certificate Authority: %s.", caId)
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("deleteCa got %q, want %q", got, expectedResult)
		}

		// Certificate Authority needs to be undeleted first, so we can delete it again permanently
		// without 30d grace period to be able to clean up CA Pool afterwards
		if err := unDeleteCa(&buf, projectId, location, caPoolId, caId); err != nil {
			t.Error("createCa teardown (unDeleteCa) got err:", err)
		}

		// We need to make sure it's completely deleted (without graceperiod before we finish tests)
		if err := deleteCaPerm(projectId, location, caPoolId, caId); err != nil {
			t.Fatal("deleteCa teardown got err:", err)
		}
	})

	t.Run("enableDisableCa", func(t *testing.T) {
		caId, teardownCa := setupCa(t, caPoolId)
		defer teardownCa(t)

		buf.Reset()
		if err := enableCa(&buf, projectId, location, caPoolId, caId); err != nil {
			t.Fatal("enableCa got err:", err)
		}

		expectedResult := fmt.Sprintf("Successfully enabled Certificate Authority: %s.", caId)
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("enableCa got %q, want %q", got, expectedResult)
		}

		buf.Reset()
		if err := disableCa(&buf, projectId, location, caPoolId, caId); err != nil {
			t.Fatal("disableCa got err:", err)
		}

		expectedResult = fmt.Sprintf("Successfully disabled Certificate Authority: %s.", caId)
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("disableCa got %q, want %q", got, expectedResult)
		}
	})
}
