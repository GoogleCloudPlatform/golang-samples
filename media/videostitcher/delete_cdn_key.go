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

// [START videostitcher_delete_cdn_key]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
)

// deleteCDNKey deletes a CDN key.
func deleteCDNKey(w io.Writer, projectID, keyID string) error {
	// projectID := "my-project-id"
	// keyID := "my-cdn-key"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/locations/%s/cdnKeys/%s", projectID, location, keyID)

	req := &stitcherstreampb.DeleteCdnKeyRequest{
		Name: name,
	}
	// Deletes the CDN key.
	op, err := client.DeleteCdnKey(ctx, req)
	if err != nil {
		return fmt.Errorf("client.DeleteCdnKey: %w", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Deleted CDN key")
	return nil
}

// [END videostitcher_delete_cdn_key]
