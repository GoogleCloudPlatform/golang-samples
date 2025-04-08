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

package snippets

// [START iam_delete_key]
import (
	"context"
	"fmt"
	"io"

	iam "google.golang.org/api/iam/v1"
)

// deleteKey deletes a service account key.
func deleteKey(w io.Writer, fullKeyName string) error {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return fmt.Errorf("iam.NewService: %w", err)
	}

	_, err = service.Projects.ServiceAccounts.Keys.Delete(fullKeyName).Do()
	if err != nil {
		return fmt.Errorf("Projects.ServiceAccounts.Keys.Delete: %w", err)
	}
	fmt.Fprintf(w, "Deleted key: %v", fullKeyName)
	return nil
}

// [END iam_delete_key]
