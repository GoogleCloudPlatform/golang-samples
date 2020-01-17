// Copyright 2019 Google LLC
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

// [START storage_get_bucket_policy]
import (
	"context"
	"log"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
)

func getPolicy(c *storage.Client, bucketName string) (*iam.Policy, error) {
	ctx := context.Background()

	policy, err := c.Bucket(bucketName).IAM().Policy(ctx)
	if err != nil {
		return nil, err
	}
	for _, role := range policy.Roles() {
		log.Printf("%q: %q", role, policy.Members(role))
	}
	return policy, nil
}

// [END storage_get_bucket_policy]
