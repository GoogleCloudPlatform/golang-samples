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

// [START videostitcher_update_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateCDNKey updates a CDN key. A CDN key is used to retrieve protected media.
// If isMediaCDN is true, update a Media CDN key. If false, update a Cloud
// CDN key. To create an updated privateKey value for Media CDN, see
// https://cloud.google.com/video-stitcher/docs/how-to/managing-cdn-keys#create-private-key-media-cdn.
func updateCDNKey(w io.Writer, projectID, keyID, privateKey string, isMediaCDN bool) error {
	// projectID := "my-project-id"
	// keyID := "my-cdn-key"
	// privateKey := "my-updated-private-key"
	// isMediaCDN := true
	location := "us-central1"
	hostname := "updated.cdn.example.com"
	keyName := "cdn-key"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	var req *stitcherstreampb.UpdateCdnKeyRequest
	if isMediaCDN {
		req = &stitcherstreampb.UpdateCdnKeyRequest{
			CdnKey: &stitcherstreampb.CdnKey{
				CdnKeyConfig: &stitcherstreampb.CdnKey_MediaCdnKey{
					MediaCdnKey: &stitcherstreampb.MediaCdnKey{
						KeyName:    keyName,
						PrivateKey: []byte(privateKey),
					},
				},
				Name:     fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectID, location, keyID),
				Hostname: hostname,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"hostname", "media_cdn_key",
				},
			},
		}
	} else {
		req = &stitcherstreampb.UpdateCdnKeyRequest{
			CdnKey: &stitcherstreampb.CdnKey{
				CdnKeyConfig: &stitcherstreampb.CdnKey_GoogleCdnKey{
					GoogleCdnKey: &stitcherstreampb.GoogleCdnKey{
						KeyName:    keyName,
						PrivateKey: []byte(privateKey),
					},
				},
				Name:     fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectID, location, keyID),
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
	op, err := client.UpdateCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("client.UpdateCdnKey: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Updated CDN key: %+v", response)
	return nil
}

// [END videostitcher_update_cdn_key]
