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

// [START storage_set_object_contexts]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// setObjectContexts patches an object's contexts.
func setObjectContexts(w io.Writer, bucket, object string) error {
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
	// Optional: set a metageneration-match precondition to avoid potential race
	// conditions and data corruptions. The request to update is aborted if the
	// object's metageneration does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})

	// Upsert a context (value is replaced if key already exists;
	// otherwise, a new key-value pair is added).
	// To delete a key, mark it as delete in payload.
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Contexts: &storage.ObjectContexts{
			Custom: map[string]storage.ObjectCustomContextPayload{
				"key1": {Value: "newValue1"},
				"key2": {Delete: true},
				"key3": {Value: "value3"},
			},
		},
	}
	updatedAttrs, err := o.Update(ctx, objectAttrsToUpdate)
	if err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %w", object, err)
	}

	fmt.Fprintf(w, "Updated object contexts for object %v in bucket %v.\n", object, bucket)

	if updatedAttrs.Contexts != nil && len(updatedAttrs.Contexts.Custom) > 0 {
		fmt.Fprintf(w, "Object contexts for %v:\n", updatedAttrs.Name)
		for key, payload := range updatedAttrs.Contexts.Custom {
			fmt.Fprintf(w, "\t%v = %v\n", key, payload.Value)
		}
	} else {
		fmt.Fprintf(w, "No contexts found for %v\n", updatedAttrs.Name)
	}
	return nil
}

// [END storage_set_object_contexts]
