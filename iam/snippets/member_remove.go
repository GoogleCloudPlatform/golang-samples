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
