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

// [START iam_get_role]
import (
	"context"
	"fmt"
	"io"

	iam "google.golang.org/api/iam/v1"
)

// getRole gets role metadata.
func getRole(w io.Writer, name string) (*iam.Role, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %w", err)
	}

	role, err := service.Roles.Get(name).Do()
	if err != nil {
		return nil, fmt.Errorf("Roles.Get: %w", err)
	}
	fmt.Fprintf(w, "Got role: %v\n", role.Name)
	for _, permission := range role.IncludedPermissions {
		fmt.Fprintf(w, "Got permission: %v\n", permission)
	}
	return role, nil
}

// [END iam_get_role]
