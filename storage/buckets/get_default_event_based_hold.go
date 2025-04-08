// Copyright 2020 Google LLC
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

package buckets

// [START storage_get_default_event_based_hold]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// getDefaultEventBasedHold gets default event-based hold state.
func getDefaultEventBasedHold(w io.Writer, bucketName string) (*storage.BucketAttrs, error) {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Bucket(%q).Attrs: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Default event-based hold enabled? %t\n", attrs.DefaultEventBasedHold)
	return attrs, nil
}

// [END storage_get_default_event_based_hold]
