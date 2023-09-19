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

const (
	location            = "us-central1" // All samples use this location
	slateID             = "my-go-test-slate"
	slateURI            = "https://storage.googleapis.com/cloud-samples-data/media/ForBiggerEscapes.mp4"
	updatedSlateURI     = "https://storage.googleapis.com/cloud-samples-data/media/ForBiggerJoyrides.mp4"
	deleteSlateResponse = "Deleted slate"

	deleteCDNKeyResponse = "Deleted CDN key"
	mediaCDNKeyID        = "my-go-test-media-cdn"
	cloudCDNKeyID        = "my-go-test-cloud-cdn"
	akamaiCDNKeyID       = "my-go-test-akamai-cdn"
	hostname             = "cdn.example.com"
	updatedHostname      = "updated.example.com"
	keyName              = "my-key"

	vodURI = "https://storage.googleapis.com/cloud-samples-data/media/hls-vod/manifest.m3u8"

	liveConfigID             = "my-go-test-live-config"
	deleteLiveConfigResponse = "Deleted live config"
	liveURI                  = "https://storage.googleapis.com/cloud-samples-data/media/hls-live/manifest.m3u8"
)

// To run the tests, do the following:
// Export the following env vars:
// *   GOOGLE_APPLICATION_CREDENTIALS
// *   GOLANG_SAMPLES_PROJECT_ID
// Enable the following API on the test project:
// *   Video Stitcher API
