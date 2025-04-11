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

package assets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/protobuf/field_mask"
)

var marksAssetName = ""

func getRandomAsset(client *securitycenter.Client, orgID string) (*securitycenterpb.Asset, error) {
	ctx := context.Background()
	req := &securitycenterpb.ListAssetsRequest{
		Parent: fmt.Sprintf("organizations/%s", orgID),
	}

	var randomAsset *securitycenterpb.Asset
	assetsCount := 0
	it := client.ListAssets(ctx, req)
	for {
		result, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("it.Next(): %w", err)
		}
		assetsCount++
		if rand.Float64() < 1.0/float64(assetsCount) {
			randomAsset = result.Asset
		}
	}
	return randomAsset, nil
}

func attemptLease(client *securitycenter.Client, asset *securitycenterpb.Asset, orgID string) error {
	const leaseExpirationKey = "LEASEKEY"
	now := time.Now().UnixNano()

	lease := asset.SecurityMarks.Marks[leaseExpirationKey]
	okToLease := true
	if lease != "" {
		i, err := strconv.ParseInt(lease, 10, 64)
		if err != nil {
			fmt.Printf("strconv.ParseInt(%v, 10, 64): %v", lease, err)
		}
		okToLease = now > i
	}
	if !okToLease {
		return fmt.Errorf("lease by another process still active for %s", asset.Name)
	}
	leaseTime := now + (60 * time.Second).Nanoseconds()
	leaseValue := strconv.FormatInt(leaseTime, 10)
	ctx := context.Background()
	_, err := client.UpdateSecurityMarks(ctx, &securitycenterpb.UpdateSecurityMarksRequest{
		// If not set or empty, all marks would be cleared before
		// adding the new marks below.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{
				fmt.Sprintf("marks.%s", leaseExpirationKey),
				"marks.key_a",
				"marks.key_b",
			},
		},
		SecurityMarks: &securitycenterpb.SecurityMarks{
			Name: fmt.Sprintf("%s/securityMarks", asset.Name),
			// Note keys correspond to the last part of each path.
			Marks: map[string]string{leaseExpirationKey: leaseValue},
		},
	})
	if err != nil {
		return fmt.Errorf("UpdateSecurityMarks: %w", err)
	}
	// Randomize wake-up in case we are in the edge case of two writes
	time.Sleep(time.Duration((100 * rand.Int63n(20))) * time.Millisecond)
	it := client.ListAssets(ctx, &securitycenterpb.ListAssetsRequest{
		Parent: fmt.Sprintf("organizations/%s", orgID),
		Filter: fmt.Sprintf(`name="%s"`, asset.Name),
	})

	result, err := it.Next()
	if err == iterator.Done {
		return fmt.Errorf("didn't find asset %s", asset.Name)
	}
	if err != nil {
		return fmt.Errorf("it.Next: %w", err)
	}
	asset = result.Asset

	if asset.SecurityMarks.Marks[leaseExpirationKey] != leaseValue {
		return fmt.Errorf("simultaneous write by another process for %s", asset.Name)
	}
	marksAssetName = asset.Name
	return nil
}

func initAssetForManipulation() error {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		// Each test checks for GCLOUD_ORGANIZATION. Return nil so we see every skip.
		return nil
	}

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	try := 0
	for try < 3 && marksAssetName == "" {
		try++
		asset, err := getRandomAsset(client, orgID)
		if err != nil {
			continue
		}
		if err := attemptLease(client, asset, orgID); err != nil {
			return fmt.Errorf("attemptLease: %s", asset.Name)
		}
	}
	if marksAssetName == "" {
		return fmt.Errorf("failed to set marksAssetName")
	}
	return nil
}

func setup(t *testing.T) string {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	} else if marksAssetName == "" {
		t.Errorf("marksAssetName wasn't initialized.")
		os.Exit(1)
	}
	return orgID
}

func TestMain(m *testing.M) {
	if err := initAssetForManipulation(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize assets test environment: %v\n", err)
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())
	code := m.Run()
	os.Exit(code)
}

func TestListAllAssets(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		orgID := setup(t)
		err := listAllAssets(buf, orgID)
		if err != nil {
			r.Errorf("listAllAssets(%s) failed: %v", orgID, err)
			return
		}

		want := 59
		got := strings.Count(buf.String(), "\n")
		if got < want {
			r.Errorf("listAllAssets(%s) Not enough results: %d Want >= %d", orgID, got, want)
		}
	})
}

func TestListAllProjectAssets(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		orgID := setup(t)
		err := listAllProjectAssets(buf, orgID)
		if err != nil {
			r.Errorf("listAllAssets(%s) failed: %v", orgID, err)
			return
		}
		want := 3
		got := strings.Count(buf.String(), "\n")
		if got != want {
			r.Errorf("listAllAssets(%s) Unexpected number of results: %d Want: %d", orgID, got, want)
		}
	})
}

func TestListAllProjectAssetsAtTime(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		orgID := setup(t)
		buf := new(bytes.Buffer)
		var nothingInstant = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

		err := listAllProjectAssetsAtTime(buf, orgID, nothingInstant)

		if err != nil {
			r.Errorf("listAllProjectAssetsAtTime(%s, %v) failed: %v", orgID, nothingInstant, err)
			return
		}

		got := strings.Count(buf.String(), "\n")
		if got != 0 {
			r.Errorf("listAllProjectAssetsAtTime(%s, %v) Results not 0: %d", orgID, nothingInstant, got)
		}

		buf.Reset()
		var somethingInstant = time.Now()
		err = listAllProjectAssetsAtTime(buf, orgID, somethingInstant)
		if err != nil {
			r.Errorf("listAllProjectAssetsAtTime(%s, %v) failed: %v", orgID, somethingInstant, err)
			return
		}
		want := 3
		got = strings.Count(buf.String(), "\n")
		if got != want {
			r.Errorf("listAllProjectAssetsAtTime(%s, %v) Unexpected number of projects: %d Want: %d", orgID, somethingInstant, got, want)
		}
	})
}

func TestListAllProjectAssetsAndStateChanges(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		orgID := setup(t)
		err := listAllProjectAssetsAndStateChanges(buf, orgID)
		if err != nil {
			r.Errorf("listAllProjectAssetsAndStateChanges(%s) failed: %v", orgID, err)
			return
		}
		got := strings.Count(buf.String(), "\n")
		want := 3
		if got != want {
			r.Errorf("listAllProjectAssetsAndStateChanges(%s) Unexpected number of results: %d Want: %d", orgID, got, want)
		}
	})
}

func TestAddSecurityMarks(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		err := deleteSecurityMarks(buf, marksAssetName)
		if err != nil {
			r.Errorf("Setup for addSecurityMarks(%s) failed %v", marksAssetName, err)
			return
		}

		err = addSecurityMarks(buf, marksAssetName)
		if err != nil {
			r.Errorf("addSecurityMarks(%s) failed %v", marksAssetName, err)
			return
		}

		got := buf.String()
		if want := "key_a = value_a"; !strings.Contains(got, want) {
			r.Errorf("addSecurityMarks(%s) got: %s want %s", marksAssetName, got, want)
		}

		if want := "key_b = value_b"; !strings.Contains(got, want) {
			r.Errorf("addSecurityMarks(%s) got: %s want %s", marksAssetName, got, want)
		}
	})
}

func TestDeleteSecurityMarks(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		err := addSecurityMarks(buf, marksAssetName)
		if err != nil {
			r.Errorf("Setup for deleteSecurityMarks(%s) failed %v", marksAssetName, err)
			return
		}
		buf.Reset()

		err = deleteSecurityMarks(buf, marksAssetName)
		if err != nil {
			r.Errorf("deleteSecurityMarks(%s) failed %v", marksAssetName, err)
			return
		}

		got := buf.String()
		if dontWant := "key_a = value_a"; strings.Contains(got, dontWant) {
			r.Errorf("deleteSecurityMarks(%s) got: %s dont want %q", marksAssetName, got, dontWant)
		}

		if dontWant := "key_b = value_b"; strings.Contains(got, dontWant) {
			r.Errorf("deleteSecurityMarks(%s) got: %s dont want %q", marksAssetName, got, dontWant)
		}
	})
}

func TestAddDeleteSecurityMarks(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		err := addSecurityMarks(buf, marksAssetName)
		if err != nil {
			r.Errorf("Setup for addDeleteSecurityMarks(%s) failed %v", marksAssetName, err)
			return
		}
		buf.Reset()

		err = addDeleteSecurityMarks(buf, marksAssetName)
		if err != nil {
			r.Errorf("addDeleteSecurityMarks(%s) failed %v", marksAssetName, err)
			return
		}

		got := buf.String()
		if want := "key_a = new_value_a"; !strings.Contains(got, want) {
			r.Errorf("addDeleteSecurityMarks(%s) got: %s want %s", marksAssetName, got, want)
		}

		if dontWant := "key_b = value_b"; strings.Contains(got, dontWant) {
			r.Errorf("addDeleteSecurityMarks(%s) got: %s dont want %q", marksAssetName, got, dontWant)
		}
	})
}

func TestListWithSecurityMarks(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		orgID := setup(t)
		err := addSecurityMarks(buf, marksAssetName)
		if err != nil {
			r.Errorf("Setup for ListWithSecurityMarks(%s) failed %v", orgID, err)
			return
		}

		err = listAssetsWithMarks(buf, orgID)

		if err != nil {
			r.Errorf("listAssetsWithMarks(%s) failed %v", orgID, err)
			return
		}

		got := buf.String()
		if !strings.Contains(got, marksAssetName) {
			r.Errorf("addDeleteSecurityMarks(%s) got: %s want %s", orgID, got, marksAssetName)
		}
	})
}
