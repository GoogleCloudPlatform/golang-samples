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

// [START object_list_acl_filtered]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// objectListACLFiltered lists ACL of the specified object with filter.
func objectListACLFiltered(w io.Writer, bucket, object string, entity storage.ACLEntity) error {
	// bucket := "bucket-name"
	// object := "object-name"
	// entity := storage.AllAuthenticatedUsers
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	rules, err := client.Bucket(bucket).ACL().List(ctx)
	if err != nil {
		return err
	}
	for _, r := range rules {
		if r.Entity == entity {
			fmt.Fprintf(w, "ACL rule role: %v\n", r.Role)
		}
	}
	return nil
}

// [END object_list_acl_filtered]
