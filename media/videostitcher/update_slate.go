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

// [START videostitcher_update_slate]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "cloud.google.com/go/video/stitcher/apiv1/stitcherpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateSlate updates an existing slate. This sample updates the uri for an
// existing slate.
func updateSlate(w io.Writer, projectID, slateID, slateURI string) error {
	// projectID := "my-project-id"
	// slateID := "my-slate-id"
	// slateURI := "https://my-updated-slate-uri/test.mp4"
	location := "us-central1"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return fmt.Errorf("stitcher.NewVideoStitcherClient: %w", err)
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
	op, err := client.UpdateSlate(ctx, req)
	if err != nil {
		return fmt.Errorf("client.UpdateSlate: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Updated slate: %+v", response)
	return nil
}

// [END videostitcher_update_slate]
