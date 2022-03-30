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

// [START storage_set_lifecycle_abort_multipart_upload]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// setLifecycleToAbortMultipartUploads adds an abort multipart upload lifecycle
// rule with the condition that the multipart upload was created 7 days ago.
func setLifecycleToAbortMultipartUploads(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()

	// Initialize the Storage client that will be used to send requests. This
	// client only needs to be created once, and should be reused for multiple
	// requests. After completing all of your requests, call the Close method on
	// the client to safely clean up any remaining background resources.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)

	// Append an abort incomplete multipart upload rule to the existing lifecycle
	// rules on the bucket
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	newLifecycleRules := append(attrs.Lifecycle.Rules, storage.LifecycleRule{
		Action: storage.LifecycleAction{Type: storage.AbortIncompleteMPUAction},
		Condition: storage.LifecycleCondition{
			AgeInDays: 7,
		},
	})

	// Update the bucket with the new lifecycle rule added
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		Lifecycle: &storage.Lifecycle{
			Rules: newLifecycleRules,
		},
	}
	uattrs, err := bucket.Update(ctx, bucketAttrsToUpdate)
	if err != nil {
		return fmt.Errorf("Bucket(%q).Update: %v", bucketName, err)
	}
	fmt.Fprintf(w, "Added a new lifecycle rule on bucket %v.\n The rules are now:\n", bucketName)
	for _, rule := range uattrs.Lifecycle.Rules {
		fmt.Fprintf(w, "Action: %v\n", rule.Action)
		fmt.Fprintf(w, "Condition: %v\n", rule.Condition)
	}

	return nil
}

// [END storage_set_lifecycle_abort_multipart_upload]
