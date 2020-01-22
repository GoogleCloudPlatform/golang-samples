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

// [START object_list_acl]
import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// objectListACL lists ACL of the specified object.
func objectListACL(bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	rules, err := client.Bucket(bucket).Object(object).ACL().List(ctx)
	if err != nil {
		return err
	}
	for _, rule := range rules {
		fmt.Printf("ACL rule: %v\n", rule)
	}
	return nil
}

// [END object_list_acl]
