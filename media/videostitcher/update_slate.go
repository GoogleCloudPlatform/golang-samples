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

// [START video_stitcher_update_slate]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateSlate updates an existing slate. This sample updates the uri for an
// existing slate.
func updateSlate(w io.Writer, projectID, location, slateID, slateURI string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// slateID := "my-slate-id"
	// slateURI := "my-updated-slate-uri"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherstreampb.UpdateSlateRequest{
		Slate: &stitcherstreampb.Slate{
			Name: fmt.Sprintf("projects/%s/locations/%s/slates/%s", projectID, location, slateID),
			Uri:  slateURI,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"uri",
			},
		},
	}
	// Updates the slate.
	response, err := client.UpdateSlate(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateSlate: %v", err)
	}

	fmt.Fprintf(w, "Updated slate: %v", response.Name)
	return nil
}

// [END video_stitcher_update_slate]
