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

// [START video_stitcher_create_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherpb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// createCdnKey creates a CDN key. A CDN key is used to retrieve protected media.
// If akamaiTokenKey != "", then this is an Akamai CDN key, or else this is a
// Cloud CDN key.
func createCdnKey(w io.Writer, projectID, cdnKeyID, hostname, gcdnKeyname, gcdnPrivateKey, akamaiTokenKey string) error {
	// projectID := "my-project-id"
	// cdnKeyID := "my-cdn-key"
	// hostname := "cdn.example.com"
	// gcdnKeyname := "gcdn-key"
	// gcdnPrivateKey := "VGhpcyBpcyBhIHRlc3Qgc3RyaW5nLg=="
	// akamaiTokenKey := "VGhpcyBpcyBhIHRlc3Qgc3RyaW5nLg=="
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	var req *stitcherpb.CreateCdnKeyRequest
	if akamaiTokenKey != "" {
		req = &stitcherpb.CreateCdnKeyRequest{
			Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, location),
			CdnKeyId: cdnKeyID,
			CdnKey: &stitcherpb.CdnKey{
				CdnKeyConfig: &stitcherpb.CdnKey_AkamaiCdnKey{
					AkamaiCdnKey: &stitcherpb.AkamaiCdnKey{
						TokenKey: []byte(akamaiTokenKey),
					},
				},
				Hostname: hostname,
			},
		}
	} else {
		req = &stitcherpb.CreateCdnKeyRequest{
			Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, location),
			CdnKeyId: cdnKeyID,
			CdnKey: &stitcherpb.CdnKey{
				CdnKeyConfig: &stitcherpb.CdnKey_GoogleCdnKey{
					GoogleCdnKey: &stitcherpb.GoogleCdnKey{
						KeyName:    gcdnKeyname,
						PrivateKey: []byte(gcdnPrivateKey),
					},
				},
				Hostname: hostname,
			},
		}
	}

	// Creates the CDN key.
	response, err := client.CreateCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateCdnKey: %v", err)
	}

	fmt.Fprintf(w, "CDN key: %v", response.GetName())
	return nil
}

// [END video_stitcher_create_cdn_key]
