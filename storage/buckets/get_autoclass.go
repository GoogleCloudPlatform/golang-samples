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

package buckets

// [START storage_get_autoclass]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// getAutoclass gets the Autoclass configuration for a bucket.
// See https://cloud.google.com/storage/docs/using-autoclass for more information.
func getAutoclass(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	fmt.Fprintf(w, "Autoclass enabled was set to %v on bucket %q at %v", attrs.Autoclass.Enabled, bucketName, attrs.Autoclass.ToggleTime)
	return nil
}

// [END storage_get_autoclass]
