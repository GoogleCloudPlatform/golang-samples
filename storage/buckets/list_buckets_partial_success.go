// Copyright 2025 Google LLC
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

// [START storage_list_buckets_partial_success]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// listBucketsPartialSuccess lists buckets in the project. If ReturnPartialSuccess
// is true, the iterator will return reachable buckets and a list of
// unreachable bucket resource names.
func listBucketsPartialSuccess(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	it := client.Buckets(ctx, projectID)
	// Enable returning unreachable buckets.
	it.ReturnPartialSuccess = true

	fmt.Fprintln(w, "Reachable buckets:")
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Errors here usually indicate a problem with the overall list operation
			// or connection, such as a network issue, rather than individual
			// buckets being unreachable. Unreachable buckets due to issues like
			// regional outages or permission issues are typically reported via the
			// Unreachable() method below.
			return err
		}
		fmt.Fprintf(w, "- %v\n", battrs.Name)
	}

	// Retrieve the list of buckets that were unreachable.
	unreachable := it.Unreachable()
	if len(unreachable) > 0 {
		fmt.Fprintln(w, "\nUnreachable buckets:")
		for _, r := range unreachable {
			fmt.Fprintf(w, "- %v\n", r)
		}
	} else {
		fmt.Fprintln(w, "\nNo unreachable buckets.")
	}

	return nil
}

// [END storage_list_buckets_partial_success]
