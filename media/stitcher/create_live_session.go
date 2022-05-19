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

// [START video_stitcher_create_live_session]
import (
	"context"
	"fmt"
	"io"

	stitcher "cloud.google.com/go/video/stitcher/apiv1"
	stitcherstreampb "google.golang.org/genproto/googleapis/cloud/video/stitcher/v1"
)

// createLiveSession creates a livestream session in which to insert ads. Live
// sessions are ephemeral resources that expire after a few minutes. This
// function returns the play URI.
func createLiveSession(w io.Writer, projectID, location, sourceURI, adTagURI, slateID string) (string, error) {
	// projectID := "my-project-id"
	// location := "us-central1"
	// sourceURI := "my-manifest.m3u8"
	// adTagURI := "ad-tag-uri"
	// slateID := "my-slate-id"
	ctx := context.Background()
	client, err := stitcher.NewVideoStitcherClient(ctx)
	if err != nil {
		return "", fmt.Errorf("NewVideoStitcherClient: %v", err)
	}
	defer client.Close()

	req := &stitcherstreampb.CreateLiveSessionRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		LiveSession: &stitcherstreampb.LiveSession{
			SourceUri: sourceURI,
			AdTagMap: map[string]*stitcherstreampb.AdTag{
				"default": &stitcherstreampb.AdTag{
					Uri: adTagURI,
				},
			},
			DefaultSlateId: slateID,
		},
	}
	// Creates the session.
	response, err := client.CreateLiveSession(ctx, req)
	if err != nil {
		return "", fmt.Errorf("CreateLiveSession: %v", err)
	}

	fmt.Fprintf(w, "Live session: %v", response.Name)
	return response.PlayUri, nil
}

// [END video_stitcher_create_live_session]
