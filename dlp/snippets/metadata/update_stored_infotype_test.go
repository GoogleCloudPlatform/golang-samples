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
package metadata

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

func TestUpdateStoredInfoType(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	outputBucket := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, bucket_prefix)

	outputPath := fmt.Sprintf("gs://%s/", outputBucket)

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, bucket_prefix)

	fileSetUrl, gcsUri, err := filesForUpdateStoredInfoType(t, tc.ProjectID, bucketName)
	if err != nil {
		t.Fatal(err)
	}

	infoTypeId, err := createStoredInfoTypeForTesting(t, tc.ProjectID, outputPath)
	if err != nil {
		t.Fatal(err)
	}

	infoTypeId = strings.TrimPrefix(infoTypeId, fmt.Sprint("projects/"+tc.ProjectID+"/locations/global/storedInfoTypes/"))

	duration := time.Duration(30) * time.Second
	time.Sleep(duration)

	if err := updateStoredInfoType(&buf, tc.ProjectID, gcsUri, fileSetUrl, infoTypeId); err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if want := "output: "; !strings.Contains(got, want) {
		t.Errorf("error from create stored infoType %q", got)
	}

	name := strings.TrimPrefix(got, "output: ")

	defer deleteStoredInfoTypeAfterTest(t, name)
}
