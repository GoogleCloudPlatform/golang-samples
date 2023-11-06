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

package videostitcher

// [START videostitcher_create_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// createCDNKey creates a CDN key. A CDN key is used to retrieve protected media.
// If isMediaCDN is true, create a Media CDN key. If false, create a Cloud
// CDN key. To create a privateKey value for Media CDN, see
// https://cloud.google.com/video-stitcher/docs/how-to/managing-cdn-keys#create-private-key-media-cdn.
func createCDNKey(w io.Writer, projectID, keyID, privateKey string, isMediaCDN bool) error {
	// projectID := "my-project-id"
	// keyID := "my-cdn-key"
	// privateKey := "my-private-key"
	// isMediaCDN := true
	location := "us-central1"
	hostname := "cdn.example.com"
	keyName := "cdn-key"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	var req *stitcherstreampb.CreateCdnKeyRequest
	if isMediaCDN {
		req = &stitcherstreampb.CreateCdnKeyRequest{
			Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, location),
			CdnKeyId: keyID,
			CdnKey: &stitcherstreampb.CdnKey{
				CdnKeyConfig: &stitcherstreampb.CdnKey_MediaCdnKey{
					MediaCdnKey: &stitcherstreampb.MediaCdnKey{
						KeyName:    keyName,
						PrivateKey: []byte(privateKey),
					},
				},
				Hostname: hostname,
			},
		}
	} else {
		req = &stitcherstreampb.CreateCdnKeyRequest{
			Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, location),
			CdnKeyId: keyID,
			CdnKey: &stitcherstreampb.CdnKey{
				CdnKeyConfig: &stitcherstreampb.CdnKey_GoogleCdnKey{
					GoogleCdnKey: &stitcherstreampb.GoogleCdnKey{
						KeyName:    keyName,
						PrivateKey: []byte(privateKey),
					},
				},
				Hostname: hostname,
			},
		}
	}

	// Creates the CDN key.
	op, err := client.CreateCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("client.CreateCdnKey: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "CDN key: %v", response.GetName())
	return nil
}

// [END videostitcher_create_cdn_key]
