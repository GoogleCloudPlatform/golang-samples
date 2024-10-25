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

package livestream

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

const (
	deleteChannelClipResponse  = "Deleted channel clip"
	deleteChannelEventResponse = "Deleted channel event"
	deleteChannelResponse      = "Deleted channel"
	deleteInputResponse        = "Deleted input"
	deleteAssetResponse        = "Deleted asset"
	startChannelResponse       = "Started channel"
	stopChannelResponse        = "Stopped channel"
	location                   = "us-central1"
	inputID                    = "my-go-test-input"
	backupInputID              = "my-go-test-backup-input"
	channelID                  = "my-go-test-channel"
	clipID                     = "my-go-test-channel-clip"
	eventID                    = "my-go-test-channel-event"
	assetID                    = "my-go-test-asset"
	poolID                     = "default" // only 1 pool supported per location
)

var bucketName string
var outputURI string
var assetURI string

// To run the tests, do the following:
// Export the following env vars:
// *   GOOGLE_APPLICATION_CREDENTIALS
// *   GOLANG_SAMPLES_PROJECT_ID
// Enable the following API on the test project:
// *   Live Stream API

// TestMain tests major operations on inputs, channels, channel
// events, assets, and pools.
func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	bucketName := testutil.TestBucket(ctx, t, tc.ProjectID, "golang-samples-livestream")
	outputURI = "gs://" + bucketName + "/test-output-channel/"
	assetURI = "gs://cloud-samples-data/media/ForBiggerEscapes.mp4"
	cleanStaleAssets(tc)
}

func cleanStaleAssets(tc testutil.Context) {
	ctx := context.Background()
	var threeHoursInSec int64 = 60 * 60 * 3
	timeNowSec := time.Now().Unix()

	client, err := livestream.NewClient(ctx)
	if err != nil {
		fmt.Printf("NewClient: %v", err)
	}
	defer client.Close()

	req := &livestreampb.ListAssetsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", tc.ProjectID, location),
	}

	it := client.ListAssets(ctx, req)
	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("ListAssets: %v", err)
			continue
		}
		req := &livestreampb.GetAssetRequest{
			Name: response.Name,
		}
		asset, err := client.GetAsset(ctx, req)
		if err != nil {
			fmt.Printf("GetAsset: %v", err)
			continue
		}
		if asset.GetCreateTime().GetSeconds() < timeNowSec-threeHoursInSec {
			fmt.Printf("%v - delete asset", asset.GetCreateTime().GetSeconds())
			req := &livestreampb.DeleteAssetRequest{
				Name: asset.GetName(),
			}
			// No need to wait for delete ops to finish, as this is a background
			// cleanup.
			_, err := client.DeleteAsset(ctx, req)
			if err != nil {
				fmt.Printf("DeleteAsset: %v", err)
				continue
			}
		}
	}
}

// TestInputs tests major operations on inputs. Create, list, update,
// and get operations check if the input resource name is returned. The
// delete operation checks for a hard-coded string response.
func TestInputs(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Test setup

	// Stop and delete the default channel if it exists
	if err := getChannel(buf, tc.ProjectID, location, channelID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := stopChannel(buf, tc.ProjectID, location, channelID); err != nil {
				// Ignore the error when the channel is already stopped
			}
		})

		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteChannel(buf, tc.ProjectID, location, channelID); err != nil {
				r.Errorf("deleteChannel got err: %v", err)
			}
		})
	}

	// Delete the default input if it exists
	if err := getInput(buf, tc.ProjectID, location, inputID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteInput(buf, tc.ProjectID, location, inputID); err != nil {
				r.Errorf("deleteInput got err: %v", err)
			}
		})
	}

	// Delete the default backup input if it exists
	if err := getInput(buf, tc.ProjectID, location, backupInputID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteInput(buf, tc.ProjectID, location, backupInputID); err != nil {
				r.Errorf("deleteInput got err: %v", err)
			}
		})
	}

	// Create a new backup input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := createInput(buf, tc.ProjectID, location, backupInputID); err != nil {
			r.Errorf("createInput got err: %v", err)
		}
	})

	// Tests

	// Create a new input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		inputName := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", tc.ProjectID, location, inputID)
		if err := createInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("createInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, inputName) {
			r.Errorf("createInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, inputName)
		}
	})
	buf.Reset()

	// List the inputs for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		inputName := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", tc.ProjectID, location, inputID)
		if err := listInputs(buf, tc.ProjectID, location); err != nil {
			r.Errorf("listInputs got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, inputName) {
			r.Errorf("listInputs got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, inputName)
		}
	})
	buf.Reset()

	// Update an existing input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		inputName := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", tc.ProjectID, location, inputID)
		if err := updateInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("updateInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, inputName) {
			r.Errorf("updateInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, inputName)
		}
	})
	buf.Reset()

	// Get the updated input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		inputName := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", tc.ProjectID, location, inputID)
		if err := getInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("getInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, inputName) {
			r.Errorf("getInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, inputName)
		}
	})

	// Delete the input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("deleteInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteInputResponse) {
			r.Errorf("deleteInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteInputResponse)
		}
	})
	t.Logf("\nTestInputs() completed\n")
}

// TestChannels tests major operations on channels. Create, list, update,
// and get operations check if the channel resource name is returned. The
// delete operation checks for a hard-coded string response.
func TestChannels(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Test setup

	// Stop and delete the default channel if it exists
	if err := getChannel(buf, tc.ProjectID, location, channelID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := stopChannel(buf, tc.ProjectID, location, channelID); err != nil {
				// Ignore the error when the channel is already stopped
			}
		})

		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteChannel(buf, tc.ProjectID, location, channelID); err != nil {
				r.Errorf("deleteChannel got err: %v", err)
			}
		})
	}

	// Delete the default input if it exists
	if err := getInput(buf, tc.ProjectID, location, inputID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteInput(buf, tc.ProjectID, location, inputID); err != nil {
				r.Errorf("deleteInput got err: %v", err)
			}
		})
	}
	buf.Reset()

	// Create a new input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		inputName := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", tc.ProjectID, location, inputID)
		if err := createInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("createInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, inputName) {
			r.Errorf("createInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, inputName)
		}
	})
	buf.Reset()

	// Tests

	// Create a new channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		channelName := fmt.Sprintf("projects/%s/locations/%s/channels/%s", tc.ProjectID, location, channelID)
		if err := createChannel(buf, tc.ProjectID, location, channelID, inputID, outputURI); err != nil {
			r.Errorf("createChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, channelName) {
			r.Errorf("createChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, channelName)
		}
	})
	buf.Reset()

	// List the channels for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		channelName := fmt.Sprintf("projects/%s/locations/%s/channels/%s", tc.ProjectID, location, channelID)
		if err := listChannels(buf, tc.ProjectID, location); err != nil {
			r.Errorf("listChannels got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, channelName) {
			r.Errorf("listChannels got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, channelName)
		}
	})
	buf.Reset()

	// Update an existing channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		channelName := fmt.Sprintf("projects/%s/locations/%s/channels/%s", tc.ProjectID, location, channelID)
		if err := updateChannel(buf, tc.ProjectID, location, channelID, inputID); err != nil {
			r.Errorf("updateChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, channelName) {
			r.Errorf("updateChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, channelName)
		}
	})
	buf.Reset()

	// Get the updated channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		channelName := fmt.Sprintf("projects/%s/locations/%s/channels/%s", tc.ProjectID, location, channelID)
		if err := getChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("getChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, channelName) {
			r.Errorf("getChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, channelName)
		}
	})
	buf.Reset()

	// Start the channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := startChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("startChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, startChannelResponse) {
			r.Errorf("startChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, startChannelResponse)
		}
	})
	buf.Reset()

	// Stop the channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := stopChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("stopChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, stopChannelResponse) {
			r.Errorf("stopChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, stopChannelResponse)
		}
	})
	buf.Reset()

	// Delete the channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("deleteChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteChannelResponse) {
			r.Errorf("deleteChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteChannelResponse)
		}
	})
	buf.Reset()

	// Create a new channel with backup input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		channelName := fmt.Sprintf("projects/%s/locations/%s/channels/%s", tc.ProjectID, location, channelID)
		if err := createChannelWithBackupInput(buf, tc.ProjectID, location, channelID, inputID, backupInputID, outputURI); err != nil {
			r.Errorf("createChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, channelName) {
			r.Errorf("createChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, channelName)
		}
	})
	buf.Reset()

	// Clean up

	// Delete the channel with backup input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("deleteChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteChannelResponse) {
			r.Errorf("deleteChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteChannelResponse)
		}
	})

	// Delete the input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("deleteInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteInputResponse) {
			r.Errorf("deleteInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteInputResponse)
		}
	})

	// Delete the backup input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteInput(buf, tc.ProjectID, location, backupInputID); err != nil {
			r.Errorf("deleteInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteInputResponse) {
			r.Errorf("deleteInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteInputResponse)
		}
	})
	t.Logf("\nTestChannels() completed\n")
}

// TestChannelEventsAndClips tests event and clip operations on channels. Create, list, and get
// operations check if the channel event or channel clip resource name is returned. The delete
// operation checks for a hard-coded string response.
func TestChannelEventsAndClips(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Test setup
	// Stop and delete the default channel if it exists
	if err := getChannel(buf, tc.ProjectID, location, channelID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := stopChannel(buf, tc.ProjectID, location, channelID); err != nil {
				// Ignore the error when the channel is already stopped.
			}
		})

		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteChannel(buf, tc.ProjectID, location, channelID); err != nil {
				r.Errorf("deleteChannel got err: %v", err)
			}
		})
	}

	// Delete the default input if it exists
	if err := getInput(buf, tc.ProjectID, location, inputID); err == nil {
		testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
			if err := deleteInput(buf, tc.ProjectID, location, inputID); err != nil {
				r.Errorf("deleteInput got err: %v", err)
			}
		})
	}

	// Create a new input.
	uri := ""
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		inputName := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", tc.ProjectID, location, inputID)
		if err := createInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("createInput got err: %v", err)
		}
		got := buf.String()
		re := regexp.MustCompile(`Uri: (.*)`)
		match := re.FindStringSubmatch(got)
		uri = match[1]
		if !strings.Contains(got, inputName) {
			r.Errorf("createInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, inputName)
		}
	})

	// Create a new channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		channelName := fmt.Sprintf("projects/%s/locations/%s/channels/%s", tc.ProjectID, location, channelID)
		if err := createChannel(buf, tc.ProjectID, location, channelID, inputID, outputURI); err != nil {
			r.Errorf("createChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, channelName) {
			r.Errorf("createChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, channelName)
		}
	})

	// Start the channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := startChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("startChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, startChannelResponse) {
			r.Errorf("startChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, startChannelResponse)
		}
	})

	buf.Reset()

	// Tests for events

	// Create a new channel event.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		eventName := fmt.Sprintf("projects/%s/locations/%s/channels/%s/events/%s", tc.ProjectID, location, channelID, eventID)
		if err := createChannelEvent(buf, tc.ProjectID, location, channelID, eventID); err != nil {
			r.Errorf("createChannelEvent got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, eventName) {
			r.Errorf("createChannelEvent got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, eventName)
		}
	})
	buf.Reset()

	// List the channel events for a given channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		eventName := fmt.Sprintf("projects/%s/locations/%s/channels/%s/events/%s", tc.ProjectID, location, channelID, eventID)
		if err := listChannelEvents(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("listChannelEvents got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, eventName) {
			r.Errorf("listChannelEvents got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, eventName)
		}
	})
	buf.Reset()

	// Get the channel event.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		eventName := fmt.Sprintf("projects/%s/locations/%s/channels/%s/events/%s", tc.ProjectID, location, channelID, eventID)
		if err := getChannelEvent(buf, tc.ProjectID, location, channelID, eventID); err != nil {
			r.Errorf("getChannelEvent got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, eventName) {
			r.Errorf("getChannelEvent got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, eventName)
		}
	})
	buf.Reset()

	// Delete the channel event.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteChannelEvent(buf, tc.ProjectID, location, channelID, eventID); err != nil {
			r.Errorf("deleteChannelEvent got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteChannelEventResponse) {
			r.Errorf("deleteChannelEvent got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteChannelEventResponse)
		}
	})

	// Tests for clips
	// Send a test stream for the clip.
	cmd := exec.Command("ffmpeg", "-re", "-f", "lavfi", "-t", "45", "-i",
		"testsrc=size=1280x720 [out0]; sine=frequency=500 [out1]", "-vcodec",
		"h264", "-acodec", "aac", "-f", "flv", uri)

	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec.Command: %v", err)
	}

	// Create a new channel clip.
	clipOutputUri := fmt.Sprintf("%sclips", outputURI)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		clipName := fmt.Sprintf("projects/%s/locations/%s/channels/%s/clips/%s", tc.ProjectID, location, channelID, clipID)
		if err := createChannelClip(buf, tc.ProjectID, channelID, clipID, clipOutputUri); err != nil {
			r.Errorf("createChannelClip got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, clipName) {
			r.Errorf("createChannelClip got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, clipName)
		}
	})
	buf.Reset()

	// List the channel clips for a given channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		clipName := fmt.Sprintf("projects/%s/locations/%s/channels/%s/clips/%s", tc.ProjectID, location, channelID, clipID)
		if err := listChannelClips(buf, tc.ProjectID, channelID); err != nil {
			r.Errorf("listChannelClips got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, clipName) {
			r.Errorf("listChannelClips got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, clipName)
		}
	})
	buf.Reset()

	// Get the channel clip.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		clipName := fmt.Sprintf("projects/%s/locations/%s/channels/%s/clips/%s", tc.ProjectID, location, channelID, clipID)
		if err := getChannelClip(buf, tc.ProjectID, channelID, clipID); err != nil {
			r.Errorf("getChannelClip got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, clipName) {
			r.Errorf("getChannelClip got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, clipName)
		}
	})
	buf.Reset()

	// Delete the channel clip.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteChannelClip(buf, tc.ProjectID, channelID, clipID); err != nil {
			r.Errorf("deleteChannelClip got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteChannelClipResponse) {
			r.Errorf("deleteChannelClip got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteChannelClipResponse)
		}
	})

	// Clean up

	// Stop the channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := stopChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("stopChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, stopChannelResponse) {
			r.Errorf("stopChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, stopChannelResponse)
		}
	})

	// Delete the channel.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteChannel(buf, tc.ProjectID, location, channelID); err != nil {
			r.Errorf("deleteChannel got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteChannelResponse) {
			r.Errorf("deleteChannel got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteChannelResponse)
		}
	})

	// Delete the input.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("deleteInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteInputResponse) {
			r.Errorf("deleteInput got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteInputResponse)
		}
	})
	t.Logf("TestChannelEventsAndClips() completed\n")
}

// TestAssets tests major operations on assets. Create, list,
// and get operations check if the asset resource name is returned. The
// delete operation checks for a hard-coded string response.
func TestAssets(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Tests

	// Create a new asset.

	testAssetID := fmt.Sprintf("%s-%s", assetID, strconv.FormatInt(time.Now().Unix(), 10))
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		testAssetName := fmt.Sprintf("projects/%s/locations/%s/assets/%s", tc.ProjectID, location, testAssetID)
		if err := createAsset(buf, tc.ProjectID, location, testAssetID, assetURI); err != nil {
			r.Errorf("createAsset got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, testAssetName) {
			r.Errorf("createAsset got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, testAssetName)
		}
	})
	buf.Reset()

	// List the assets for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		testAssetName := fmt.Sprintf("projects/%s/locations/%s/assets/%s", tc.ProjectID, location, testAssetID)
		if err := listAssets(buf, tc.ProjectID, location); err != nil {
			r.Errorf("listAssets got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, testAssetName) {
			r.Errorf("listAssets got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, testAssetName)
		}
	})
	buf.Reset()

	// Get the asset.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		testAssetName := fmt.Sprintf("projects/%s/locations/%s/assets/%s", tc.ProjectID, location, testAssetID)
		if err := getAsset(buf, tc.ProjectID, location, testAssetID); err != nil {
			r.Errorf("getAsset got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, testAssetName) {
			r.Errorf("getAsset got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, testAssetName)
		}
	})
	buf.Reset()

	// Delete the asset.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteAsset(buf, tc.ProjectID, location, testAssetID); err != nil {
			r.Errorf("deleteAsset got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteAssetResponse) {
			r.Errorf("deleteAsset got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteAssetResponse)
		}
	})
	buf.Reset()
	t.Logf("\nTestAssets() completed\n")
}

// TestPools tests major operations on pool. Get and update
// operations check if the pool resource name is returned.
func TestPools(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Tests

	// Get the pool.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		poolName := fmt.Sprintf("projects/%s/locations/%s/pools/%s", tc.ProjectID, location, poolID)
		if err := getPool(buf, tc.ProjectID, location, poolID); err != nil {
			r.Errorf("getPool got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, poolName) {
			r.Errorf("getPool got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, poolName)
		}
	})
	buf.Reset()

	// Update an existing pool. Set the updated peer network to "", which
	// is the same as the default otherwise the test will take a long time
	// to complete.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		poolName := fmt.Sprintf("projects/%s/locations/%s/pools/%s", tc.ProjectID, location, poolID)
		if err := updatePool(buf, tc.ProjectID, location, poolID, ""); err != nil {
			r.Errorf("updatePool got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, poolName) {
			r.Errorf("updatePool got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, poolName)
		}
	})
	buf.Reset()
	t.Logf("\nTestPools() completed\n")
}
