// Copyright 2019 Google LLC
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

package howto

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
)

func TestCreateClientEvent(t *testing.T) {
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		tc := testutil.SystemTest(t)
		buf := &bytes.Buffer{}

		requestID := fmt.Sprintf("requestId-%s", uuid.Must(uuid.NewV4()).String())
		eventID := fmt.Sprintf("eventId-%s", uuid.Must(uuid.NewV4()).String())
		relatedJobNames := []string{testJob.Name}

		if _, err := createClientEvent(buf, tc.ProjectID, requestID, eventID, relatedJobNames); err != nil {
			log.Fatalf("createClientEvent: %v", err)
		}
		want := "Client event created: "
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("getJob got %q, want to contain %q", got, want)
		}
	})
}
