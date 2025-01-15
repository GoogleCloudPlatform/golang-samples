// Copyright 2024 Google LLC
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

// [START storage_set_object_retention_policy]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// setObjectRetentionPolicy sets the object retention policy of an object.
func setObjectRetentionPolicy(w io.Writer, bucket, object string) error {
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

	// Update the object to set the retention policy.
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Retention: &storage.ObjectRetention{
			Mode:        "Unlocked",
			RetainUntil: time.Now().Add(time.Hour * 24 * 10),
		},
	}
	attrs, err := o.Update(ctx, objectAttrsToUpdate)
	if err != nil {
		return fmt.Errorf("Object(%q).Update: %w", object, err)
	}
	fmt.Fprintf(w, "Retention policy for object %s was set to %v\n", object, attrs.Retention)

	// To modify an existing policy on an Unlocked object, set
	// OverrideUnlockedRetention on the ObjectHandle.
	objectAttrsToUpdate.Retention.RetainUntil = time.Now().Add(time.Hour * 24 * 9)
	attrs, err = o.OverrideUnlockedRetention(true).Update(ctx, objectAttrsToUpdate)
	if err != nil {
		return fmt.Errorf("Object(%q).Update: %w", object, err)
	}
	fmt.Fprintf(w, "Retention policy for object %s was updated to %v\n", object, attrs.Retention)
	return nil
}

// [END storage_set_object_retention_policy]
