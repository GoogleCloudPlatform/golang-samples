// Copyright 2021 Google LLC
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

package retry

// [START storage_configure_retries]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/googleapis/gax-go/v2"
)

// configureRetries configures a custom retry strategy for an API call.
func configureRetries(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Configure retries for all operations using this ObjectHandle. Retries may
	// also be configured on the BucketHandle or Client types.
	o := client.Bucket(bucket).Object(object).Retryer(
		// Use WithBackoff to control the timing of the exponential backoff. As
		// configured here, there will be an initial delay between retries of 2s,
		// a maximum delay of 60s, and the pause length will be multiplied by 3 in
		// each iteration.
		storage.WithBackoff(gax.Backoff{
			Initial:    2 * time.Second,
			Max:        60 * time.Second,
			Multiplier: 3,
		}),
		// Use WithPolicy to configure the idempotency policy. By default, only
		// idempotent operations are retried. RetryAlways will retry the operation
		// even if it is non-idempotent.
		storage.WithPolicy(storage.RetryAlways),
	)

	// Use context timeouts to set an overall deadline on the call, including all
	// potential retries.
	ctx, cancel := context.WithTimeout(ctx, 500*time.Second)
	defer cancel()

	// Delete an object using the specified policy.
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	fmt.Fprintf(w, "Blob %v deleted with a customized retry strategy.\n", object)
	return nil
}

// [END storage_configure_retries]
