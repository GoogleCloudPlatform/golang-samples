// Copyright 2022 Google LLC
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

package videostitcher

import (
	"context"
	"fmt"
	"log"

	"strconv"
	"strings"
	"testing"
	"time"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

const (
	location            = "us-central1" // All samples use this location
	slateIDPrefix       = "go-test-slate"
	slateURI            = "https://storage.googleapis.com/cloud-samples-data/media/ForBiggerEscapes.mp4"
	updatedSlateURI     = "https://storage.googleapis.com/cloud-samples-data/media/ForBiggerJoyrides.mp4"
	deleteSlateResponse = "Deleted slate"

	deleteCDNKeyResponse = "Deleted CDN key"
	mediaCDNKeyIDPrefix  = "go-test-media-cdn"
	cloudCDNKeyIDPrefix  = "go-test-cloud-cdn"
	akamaiCDNKeyIDPrefix = "go-test-akamai-cdn"
	hostname             = "cdn.example.com"
	updatedHostname      = "updated.cdn.example.com"
	keyName              = "my-key"

	vodURI = "https://storage.googleapis.com/cloud-samples-data/media/hls-vod/manifest.m3u8"

	liveConfigIDPrefix       = "go-test-live-config"
	deleteLiveConfigResponse = "Deleted live config"
	liveURI                  = "https://storage.googleapis.com/cloud-samples-data/media/hls-live/manifest.m3u8"
	liveAdTagURI             = "https://pubads.g.doubleclick.net/gampad/ads?iu=/21775744923/external/single_ad_samples&sz=640x480&cust_params=sample_ct%3Dlinear&ciu_szs=300x250%2C728x90&gdfp_req=1&output=vast&unviewed_position_start=1&env=vp&impl=s&correlator="
)

// To run the tests, do the following:
// Export the following env vars:
// *   GOOGLE_APPLICATION_CREDENTIALS
// *   GOLANG_SAMPLES_PROJECT_ID
// Enable the following API on the test project:
// *   Video Stitcher API

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatal("couldn't initialize test")
		return
	}
	cleanStaleResources(tc.ProjectID)
	m.Run()
}

func cleanStaleResources(projectID string) {
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		log.Fatalf("stitcher.NewVideoStitcherClient")
		return
	}
	defer client.Close()

	// Live configs
	req := &stitcherstreampb.ListLiveConfigsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := client.ListLiveConfigs(ctx, req)

	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Can't find next live config: %s", err)
			continue
		}
		if strings.Contains(response.GetName(), liveConfigIDPrefix) {

			arr := strings.Split(response.GetName(), "-")
			t := arr[len(arr)-1]
			if isResourceStale(t) == true {
				req := &stitcherstreampb.DeleteLiveConfigRequest{
					Name: response.GetName(),
				}
				// Deletes the live config.
				op, err := client.DeleteLiveConfig(ctx, req)
				if err != nil {
					log.Printf("cleanStaleResources DeleteLiveConfig: %s", err)
				}
				err = op.Wait(ctx)
				if err != nil {
					log.Printf("cleanStaleResources Wait: %s", err)
				}
			}
		}
	}

	// Slates
	req2 := &stitcherstreampb.ListSlatesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it2 := client.ListSlates(ctx, req2)

	for {
		response, err := it2.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Can't find next slate: %s", err)
			continue
		}
		if strings.Contains(response.GetName(), slateIDPrefix) {

			arr := strings.Split(response.GetName(), "-")
			t := arr[len(arr)-1]
			if isResourceStale(t) == true {
				req := &stitcherstreampb.DeleteSlateRequest{
					Name: response.GetName(),
				}
				// Deletes the slate.
				op, err := client.DeleteSlate(ctx, req)
				if err != nil {
					log.Printf("cleanStaleResources DeleteSlate: %s", err)
				}
				err = op.Wait(ctx)
				if err != nil {
					log.Printf("cleanStaleResources Wait: %s", err)
				}
			}
		}
	}

	// CDN keys
	req3 := &stitcherstreampb.ListCdnKeysRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it3 := client.ListCdnKeys(ctx, req3)

	for {
		response, err := it3.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Can't find next CDN key: %s", err)
			continue
		}
		if strings.Contains(response.GetName(), mediaCDNKeyIDPrefix) ||
			strings.Contains(response.GetName(), cloudCDNKeyIDPrefix) ||
			strings.Contains(response.GetName(), akamaiCDNKeyIDPrefix) {

			arr := strings.Split(response.GetName(), "-")
			t := arr[len(arr)-1]
			if isResourceStale(t) == true {
				req := &stitcherstreampb.DeleteCdnKeyRequest{
					Name: response.GetName(),
				}
				// Deletes the CDN key.
				op, err := client.DeleteCdnKey(ctx, req)
				if err != nil {
					log.Printf("cleanStaleResources DeleteCdnKey: %s", err)
				}
				err = op.Wait(ctx)
				if err != nil {
					log.Printf("cleanStaleResources Wait: %s", err)
				}
			}
		}
	}
}

func isResourceStale(timestamp string) bool {
	const threeHoursInSecs = 3 * 60 * 60
	past, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Printf("isResourceStale timestamp: %s, err: %s", timestamp, err)
		return false
	}

	now := time.Now().Unix()
	if past < (now - threeHoursInSecs) {
		return true
	}
	return false
}

func getUUID() (string, error) {
	t := time.Now()
	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("uuid err: %v", err)
	}
	uuid := u.String()
	return fmt.Sprintf("%s-%d", strings.ReplaceAll(uuid, "-", ""), t.Unix()), nil
}

func getUUID64() (string, error) {
	u1, err1 := uuid.NewRandom()
	u2, err2 := uuid.NewRandom()
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("getUUID64 err: %v, %v", err1, err2)
	}
	uuid := fmt.Sprintf("%s%s", u1.String(), u2.String())
	return strings.ReplaceAll(uuid, "-", ""), nil
}

func createTestSlate(slateID string, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	tc := testutil.SystemTest(t)
	req := &stitcherstreampb.CreateSlateRequest{
		Parent:  fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		SlateId: slateID,
		Slate: &stitcherstreampb.Slate{
			Uri: slateURI,
		},
	}
	op, err := client.CreateSlate(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = op.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestSlate(slateName string, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherstreampb.DeleteSlateRequest{
		Name: slateName,
	}
	op, err := client.DeleteSlate(ctx, req)
	if err != nil {
		t.Errorf("client.DeleteSlate: %v", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		t.Error(err)
	}
}

func createTestLiveConfig(slateID, liveConfigID string, t *testing.T) {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	tc := testutil.SystemTest(t)
	req := &stitcherstreampb.CreateLiveConfigRequest{
		Parent:       fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		LiveConfigId: liveConfigID,
		LiveConfig: &stitcherstreampb.LiveConfig{
			SourceUri:       liveURI,
			AdTagUri:        liveAdTagURI,
			AdTracking:      stitcherstreampb.AdTracking_SERVER,
			StitchingPolicy: stitcherstreampb.LiveConfig_CUT_CURRENT,
			DefaultSlate:    fmt.Sprintf("projects/%s/locations/%s/slates/%s", tc.ProjectID, location, slateID),
		},
	}
	op, err := client.CreateLiveConfig(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = op.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestLiveConfig(liveConfigName string, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	// Delete the live config.
	req := &stitcherstreampb.DeleteLiveConfigRequest{
		Name: liveConfigName,
	}
	op, err := client.DeleteLiveConfig(ctx, req)
	if err != nil {
		t.Errorf("client.DeleteLiveConfig: %v", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		t.Error(err)
	}
}

func createTestMediaCDNKey(keyID string, t *testing.T) {
	t.Helper()
	ctx := context.Background()

	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	// Create a random private key for the CDN key. It is not validated.
	mediaCDNPrivateKey, err := getUUID64()
	if err != nil {
		t.Fatalf("getUUID64 err: %v", err)
	}

	tc := testutil.SystemTest(t)
	req := &stitcherstreampb.CreateCdnKeyRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
		CdnKeyId: keyID,
		CdnKey: &stitcherstreampb.CdnKey{
			CdnKeyConfig: &stitcherstreampb.CdnKey_MediaCdnKey{
				MediaCdnKey: &stitcherstreampb.MediaCdnKey{
					KeyName:    keyName,
					PrivateKey: []byte(mediaCDNPrivateKey),
				},
			},
			Hostname: hostname,
		},
	}
	op, err := client.CreateCdnKey(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	_, err = op.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestCDNKey(CDNKeyName string, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		t.Fatalf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherstreampb.DeleteCdnKeyRequest{
		Name: CDNKeyName,
	}
	op, err := client.DeleteCdnKey(ctx, req)
	if err != nil {
		t.Errorf("client.DeleteCdnKey: %v", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		t.Error(err)
	}
}
