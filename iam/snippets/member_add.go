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
