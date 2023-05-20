package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/iam/v1"
)

func TestRemoveMember(t *testing.T) {
	tests := []struct {
		name, role, member, wantOutput string
		wantPolicy                     *iam.Policy
	}{
		{
			name:       "Role found; member removed",
			role:       "roles/viewer",
			member:     "user:foo@example.com",
			wantOutput: "Member removed.",
			wantPolicy: &iam.Policy{
				Bindings: []*iam.Binding{},
			},
		},
		{
			name:       "Role found; member not found",
			role:       "roles/viewer",
			member:     "user:bar@example.com",
			wantOutput: "Member not found.",
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
			name:       "Role not found; member not removed",
			role:       "roles/owner",
			member:     "user:bar@example.com",
			wantOutput: "Member not removed.",
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
			removeMember(buf, policy, tc.role, tc.member)
			if !strings.Contains(buf.String(), tc.wantOutput) {
				t.Errorf("removeMember got output %q, want output %q", buf.String(), tc.wantOutput)
			}
			if diff := cmp.Diff(tc.wantPolicy, policy); diff != "" {
				t.Errorf("removeMember returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
