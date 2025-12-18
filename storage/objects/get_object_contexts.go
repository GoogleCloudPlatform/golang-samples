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

// [START storage_get_object_contexts]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// getObjectContexts gets an object's contexts.
func getObjectContexts(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).Attrs: %w", object, err)
	}

	if attrs.Contexts != nil && len(attrs.Contexts.Custom) > 0 {
		fmt.Fprintf(w, "Object contexts for %v:\n", attrs.Name)
		for key, payload := range attrs.Contexts.Custom {
			fmt.Fprintf(w, "\t%v = %v\n", key, payload.Value)
		}
	} else {
		fmt.Fprintf(w, "No contexts found for %v\n", attrs.Name)
	}
	return nil
}

// [END storage_get_object_contexts]
