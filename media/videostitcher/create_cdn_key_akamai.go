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

// [START videostitcher_create_cdn_key_akamai]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// createCDNKeyAkamai creates an Akamai CDN key. A CDN key is used to retrieve
// protected media.
func createCDNKeyAkamai(w io.Writer, projectID, keyID, akamaiTokenKey string) error {
	// projectID := "my-project-id"
	// keyID := "my-cdn-key"
	// akamaiTokenKey := "my-private-token-key"
	location := "us-central1"
	hostname := "cdn.example.com"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	req := &stitcherstreampb.CreateCdnKeyRequest{
		Parent:   fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		CdnKeyId: keyID,
		CdnKey: &stitcherstreampb.CdnKey{
			CdnKeyConfig: &stitcherstreampb.CdnKey_AkamaiCdnKey{
				AkamaiCdnKey: &stitcherstreampb.AkamaiCdnKey{
					TokenKey: []byte(akamaiTokenKey),
				},
			},
			Hostname: hostname,
		},
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

// [END videostitcher_create_cdn_key_akamai]
