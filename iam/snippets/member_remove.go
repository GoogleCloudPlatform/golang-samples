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

// [START iam_modify_policy_remove_member]
import (
	"fmt"
	"io"

	"google.golang.org/api/iam/v1"
)

// removeMember removes a member from a role binding.
func removeMember(w io.Writer, policy *iam.Policy, role, member string) {
	bindings := policy.Bindings
	bindingIndex, memberIndex := -1, -1
	for bIdx := range bindings {
		if bindings[bIdx].Role != role {
			continue
		}
		bindingIndex = bIdx
		for mIdx := range bindings[bindingIndex].Members {
			if bindings[bindingIndex].Members[mIdx] != member {
				continue
			}
			memberIndex = mIdx
			break
		}
	}
	if bindingIndex == -1 {
		fmt.Fprintf(w, "Role %q not found. Member not removed.\n", role)
		return
	}
	if memberIndex == -1 {
		fmt.Fprintf(w, "Role %q found. Member not found.\n", role)
		return
	}

	members := removeIdx(bindings[bindingIndex].Members, memberIndex)
	bindings[bindingIndex].Members = members
	if len(members) == 0 {
		bindings = removeIdx(bindings, bindingIndex)
		policy.Bindings = bindings
	}
	fmt.Fprintf(w, "Role %q found. Member removed.\n", role)
}

// removeIdx removes arr[idx] from arr.
func removeIdx[T any](arr []T, idx int) []T {
	return append(arr[:idx], arr[idx+1:]...)
}

// [END iam_modify_policy_remove_member]
