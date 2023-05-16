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

package acl

// [START storage_print_file_acl]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// printFileACL lists ACL of the specified object.
func printFileACL(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	rules, err := client.Bucket(bucket).Object(object).ACL().List(ctx)
	if err != nil {
		return fmt.Errorf("ACLHandle.List: %w", err)
	}
	for _, rule := range rules {
		fmt.Fprintf(w, "ACL rule: %v\n", rule)
	}
	return nil
}

// [END storage_print_file_acl]
