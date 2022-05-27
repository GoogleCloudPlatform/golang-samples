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

// [START video_stitcher_update_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherpb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateCdnKey updates a CDN key.
func updateCdnKey(w io.Writer, projectID, cdnKeyID, hostname, gcdnKeyname, gcdnPrivateKey, akamaiTokenKey string) error {
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

	var req *stitcherpb.UpdateCdnKeyRequest
	if akamaiTokenKey != "" {
		req = &stitcherpb.UpdateCdnKeyRequest{
			CdnKey: &stitcherpb.CdnKey{
				CdnKeyConfig: &stitcherpb.CdnKey_AkamaiCdnKey{
					AkamaiCdnKey: &stitcherpb.AkamaiCdnKey{
						TokenKey: []byte(akamaiTokenKey),
					},
				},
				Name:     fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectID, location, cdnKeyID),
				Hostname: hostname,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"hostname", "akamai_cdn_key",
				},
			},
		}
	} else {
		req = &stitcherpb.UpdateCdnKeyRequest{
			CdnKey: &stitcherpb.CdnKey{
				CdnKeyConfig: &stitcherpb.CdnKey_GoogleCdnKey{
					GoogleCdnKey: &stitcherpb.GoogleCdnKey{
						KeyName:    gcdnKeyname,
						PrivateKey: []byte(gcdnPrivateKey),
					},
				},
				Name:     fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectID, location, cdnKeyID),
				Hostname: hostname,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"hostname", "google_cdn_key",
				},
			},
		}
	}

	// Updates the CDN key.
	response, err := client.UpdateCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("client.UpdateCdnKey: %v", err)
	}

	fmt.Fprintf(w, "Updated CDN key: %+v", response)
	return nil
}

// [END video_stitcher_update_cdn_key]
