// Copyright 2023 Google LLC
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

// [START iam_modify_policy_add_member]
import (
	"fmt"
	"io"

	"google.golang.org/api/iam/v1"
)

// addMember adds a member to a role binding.
func addMember(w io.Writer, policy *iam.Policy, role, member string) {
	for _, binding := range policy.Bindings {
		if binding.Role != role {
			continue
		}
		for _, m := range binding.Members {
			if m != member {
				continue
			}
			fmt.Fprintf(w, "Role %q found. Member already exists.\n", role)
			return
		}
		binding.Members = append(binding.Members, member)
		fmt.Fprintf(w, "Role %q found. Member added.\n", role)
		return
	}
	fmt.Fprintf(w, "Role %q not found. Member not added.\n", role)
}

// [END iam_modify_policy_add_member]
