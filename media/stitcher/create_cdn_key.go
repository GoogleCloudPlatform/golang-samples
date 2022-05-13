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

// [START video_stitcher_create_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// createCdnKey creates a CDN key. A CDN key is used to retrieve protected media.
func createCdnKey(w io.Writer, projectID, location, cdnKeyID, hostname, gcdnKeyname, gcdnPrivateKey, akamaiTokenKey string) error {
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

	if akamaiTokenKey != "" {
		cdnKey.CdnKeyConfig = &stitcherstreampb.CdnKey_AkamaiCdnKey{
			AkamaiCdnKey: &stitcherstreampb.AkamaiCdnKey{
				TokenKey: []byte(akamaiTokenKey),
			},
		}
	} else if gcdnKeyname != "" {
		cdnKey.CdnKeyConfig = &stitcherstreampb.CdnKey_GoogleCdnKey{
			GoogleCdnKey: &stitcherstreampb.GoogleCdnKey{
				KeyName:    gcdnKeyname,
				PrivateKey: []byte(gcdnPrivateKey),
			},
		}
	}

	req := &stitcherstreampb.CreateCdnKeyRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		CdnKeyId: cdnKeyID,
		CdnKey:   &cdnKey,
	}
	// Creates the CDN key.
	response, err := client.CreateCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateCdnKey: %v", err)
	}

	fmt.Fprintf(w, "CDN key: %v", response.Name)
	return nil
}

// [END video_stitcher_create_cdn_key]
