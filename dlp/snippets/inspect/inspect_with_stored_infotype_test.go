// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inspect

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInspectWithStoredInfotype(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	outputBucketPathForStoredInfotype := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
	outputPath := fmt.Sprintf("gs://" + outputBucketPathForStoredInfotype + "/")

	infoTypeId, err := createStoredInfoTypeForTesting(t, tc.ProjectID, outputPath)
	if err != nil {
		t.Fatal(err)
	}

	duration := time.Duration(45) * time.Second
	time.Sleep(duration)

	textToDeidentify := "This commit was made by kewin2010"
	if err := inspectWithStoredInfotype(&buf, tc.ProjectID, infoTypeId, textToDeidentify); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Quote: kewin2010"; !strings.Contains(got, want) {
		t.Errorf("TestInspectWithStoredInfotype got %q, want %q", got, want)
	}
	if want := "Info type: GITHUB_LOGINS"; !strings.Contains(got, want) {
		t.Errorf("TestInspectWithStoredInfotype got %q, want %q", got, want)
	}

	defer deleteStoredInfoTypeAfterTest(t, infoTypeId)
}
