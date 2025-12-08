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

// setObjectContexts sets an object's contexts.
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

	// To set contexts on a new object during write (followed by writing
	// data and closing the writer to set contexts):
	//
	// writer := o.NewWriter(ctx)
	// writer.Contexts = &storage.ObjectContexts{
	// 	Custom: map[string]storage.ObjectCustomContextPayload{
	// 		"keyOnWrite": {Value: "valueOnWrite"},
	// 	},
	// }

	// Optional: set a metageneration-match precondition to avoid potential race
	// conditions and data corruptions. The request to update is aborted if the
	// object's metageneration does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})

	// Upsert a context (value is replaced if key already exists;
	// otherwise, a new key-value pair is added):
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Contexts: &storage.ObjectContexts{
			Custom: map[string]storage.ObjectCustomContextPayload{
				"key1": {Value: "value1"},
				"key2": {Value: "value2"},
			},
		},
	}

	// To delete all existing contexts:
	// objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
	// 	Contexts: &storage.ObjectContexts{
	// 		Custom: map[string]storage.ObjectCustomContextPayload{},
	// 	},
	// }

	// To delete a specific key from the context:
	// objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
	// 	Contexts: &storage.ObjectContexts{
	// 		Custom: map[string]storage.ObjectCustomContextPayload{
	// 			"keyToBeDeleted": {Delete: true},
	// 		},
	// 	},
	// }

	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %w", object, err)
	}
	fmt.Fprintf(w, "Updated object contexts for object %v in bucket %v.\n", object, bucket)
	return nil
}

// [END storage_set_object_contexts]
