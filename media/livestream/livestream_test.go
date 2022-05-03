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
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	deleteChannelEventResponse = "Deleted channel event"
	deleteChannelResponse      = "Deleted channel"
	deleteInputResponse        = "Deleted input"
	startChannelResponse       = "Started channel"
	stopChannelResponse        = "Stopped channel"
	location                   = "us-central1"
	inputID                    = "my-go-test-input"
	backupInputID              = "my-go-test-backup-input"
	channelID                  = "my-go-test-channel"
	eventID                    = "my-go-test-channel-event"
)

// To run the tests, do the following:
// Export the following env vars:
// *   GOOGLE_APPLICATION_CREDENTIALS
// *   GOLANG_SAMPLES_PROJECT_ID
// Enable the following API on the test project:
// *   Live Stream API

// TestLiveStream tests major operations on inputs, channels, and channel
// events.
func TestLiveStream(t *testing.T) {
	tc := testutil.SystemTest(t)

	bucketName := tc.ProjectID + "-golang-samples-livestream-test"
	outputURI := "gs://" + bucketName + "/test-output-channel/"

	testInputs(t)
	t.Logf("\ntestInputs() completed\n")

	testChannels(t, outputURI)
	t.Logf("\ntestChannels() completed\n")

	testChannelEvents(t, outputURI)
	t.Logf("\ntestChannelEvents() completed\n")
}

// testInputs tests major operations on inputs. Create, list, update,
// and get operations check if the input resource name is returned. The
// delete operation checks for a hard-coded string response.
func testInputs(t *testing.T) {
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
}

// testChannels tests major operations on channels. Create, list, update,
// and get operations check if the channel resource name is returned. The
// delete operation checks for a hard-coded string response.
func testChannels(t *testing.T, outputURI string) {
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
}

// testChannelEvents tests event operations on channels. Create, list, and get
// operations check if the channel event resource name is returned. The delete
// operation checks for a hard-coded string response.
func testChannelEvents(t *testing.T, outputURI string) {
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
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		inputName := fmt.Sprintf("projects/%s/locations/%s/inputs/%s", tc.ProjectID, location, inputID)
		if err := createInput(buf, tc.ProjectID, location, inputID); err != nil {
			r.Errorf("createInput got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, inputName) {
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

	// Tests

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
}
