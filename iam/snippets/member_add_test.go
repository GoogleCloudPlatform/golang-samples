package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/iam/v1"
)

func TestAddMember(t *testing.T) {
	tests := []struct {
		name, role, member, wantOutput string
		wantPolicy                     *iam.Policy
	}{
		{
			name:       "Role found; member added",
			role:       "roles/viewer",
			member:     "user:bar@example.com",
			wantOutput: "Member added.",
			wantPolicy: &iam.Policy{
				Bindings: []*iam.Binding{
					{
						Role: "roles/viewer",
						Members: []string{
							"user:foo@example.com",
							"user:bar@example.com",
						},
					},
				},
			},
		},
		{
			name:       "Role found; member already exists",
			role:       "roles/viewer",
			member:     "user:foo@example.com",
			wantOutput: "Member already exists.",
			wantPolicy: &iam.Policy{
				Bindings: []*iam.Binding{
					{
						Role:    "roles/viewer",
						Members: []string{"user:foo@example.com"},
					},
				},
			},
		},
		{
			name:       "Role not found; member not added",
			role:       "roles/owner",
			member:     "user:bar@example.com",
			wantOutput: "Member not added.",
			wantPolicy: &iam.Policy{
				Bindings: []*iam.Binding{
					{
						Role:    "roles/viewer",
						Members: []string{"user:foo@example.com"},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			policy := &iam.Policy{
				Bindings: []*iam.Binding{
					{
						Role:    "roles/viewer",
						Members: []string{"user:foo@example.com"},
					},
				},
			}
			addMember(buf, policy, tc.role, tc.member)
			if !strings.Contains(buf.String(), tc.wantOutput) {
				t.Errorf("addMember got output %q, want output %q", buf.String(), tc.wantOutput)
			}
			if diff := cmp.Diff(tc.wantPolicy, policy); diff != "" {
				t.Errorf("addMember returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
