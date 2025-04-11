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

package objects

// [START storage_download_file_requester_pays]
import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// downloadUsingRequesterPays downloads an object using billing project.
func downloadUsingRequesterPays(w io.Writer, bucket, object, billingProjectID string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	// billingProjectID := "billing_account_id"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	b := client.Bucket(bucket).UserProject(billingProjectID)
	src := b.Object(object)

	// Open local file.
	f, err := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := src.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := rc.Close(); err != nil {
		return fmt.Errorf("Reader.Close: %w", err)
	}
	fmt.Fprintf(w, "Downloaded using %v as billing project.\n", billingProjectID)
	return nil
}

// [END storage_download_file_requester_pays]
