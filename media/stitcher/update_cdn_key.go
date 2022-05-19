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

package stitcher

// [START video_stitcher_update_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateCdnKey updates an existing CDN key. This sample updates the uri for an
// existing slate.
func updateCdnKey(w io.Writer, projectID, location, cdnKeyID, hostname, gcdnKeyname, gcdnPrivateKey, akamaiTokenKey string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// cdnKeyID := "my-cdn-key-id"
	// hostname := "cdn.example.com"
	// gcdnKeyname := "my-gcdn-key"
	// gcdnPrivateKey := "VGhpcyBpcyBhIHRlc3Qgc3RyaW5nLg" // Will be converted to []byte
	// akamaiTokenKey := "VGhpcyBpcyBhIHRlc3Qgc3RyaW5nLg" // Will be converted to []byte
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	var cdnKey stitcherstreampb.CdnKey
	cdnKey.Hostname = hostname
	cdnKey.Name = fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectID, location, cdnKeyID)

	var updateMask fieldmaskpb.FieldMask

	if akamaiTokenKey != "" {
		cdnKey.CdnKeyConfig = &stitcherstreampb.CdnKey_AkamaiCdnKey{
			AkamaiCdnKey: &stitcherstreampb.AkamaiCdnKey{
				TokenKey: []byte(akamaiTokenKey),
			},
		}
		updateMask.Paths = []string{"hostname", "akamai_cdn_key"}
	} else if gcdnKeyname != "" {
		cdnKey.CdnKeyConfig = &stitcherstreampb.CdnKey_GoogleCdnKey{
			GoogleCdnKey: &stitcherstreampb.GoogleCdnKey{
				KeyName:    gcdnKeyname,
				PrivateKey: []byte(gcdnPrivateKey),
			},
		}
		updateMask.Paths = []string{"hostname", "google_cdn_key"}
	} else {
		updateMask.Paths = []string{"hostname"}
	}

	req := &stitcherstreampb.UpdateCdnKeyRequest{
		CdnKey:     &cdnKey,
		UpdateMask: &updateMask,
	}
	// Updates the CDN key.
	response, err := client.UpdateCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateCdnKey: %v", err)
	}

	fmt.Fprintf(w, "Updated CDN key: %v", response.Name)
	return nil
}

// [END video_stitcher_update_cdn_key]
