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

package objects

// [START storage_list_objects_with_context_filter]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// listObjectsWithFilter lists objects using a context filter.
func listObjectsWithContextFilter(w io.Writer, bucket, filter string) error {
	// bucket := "bucket-name"
	// filter := "contexts.\"keyA\"=\"valueA\""

	/*
	 * More examples of filters:
	 * List any object that has a context with the specified key attached
	 * filter := "contexts.\"KEY\":*";
	 *
	 * List any object that that does not have a context with the specified key attached
	 * filter := "-contexts.\"KEY\":*";
	 *
	 * List any object that has a context with the specified key and value attached
	 * filter := "contexts.\"KEY\"=\"VALUE\"";
	 *
	 * List any object that does not have a context with the specified key and value attached
	 * filter := "-contexts.\"KEY\"=\"VALUE\"";
	 */

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := client.Bucket(bucket).Objects(ctx, &storage.Query{
		Filter: filter,
	})
	fmt.Fprintln(w, "Filtered objects: ")
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(%q).Objects(): %w", bucket, err)
		}
		fmt.Fprintf(w, "\t%v\n", attrs.Name)
	}
	return nil
}

// [END storage_list_objects_with_context_filter]
