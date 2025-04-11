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

// [START iam_quickstart]

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func main() {
	// TODO: Add your project ID
	projectID := flag.String("project_id", "", "Cloud Project ID")
	// TODO: Add the ID of your principal.
	// For examples, see https://cloud.google.com/iam/docs/principal-identifiers
	member := flag.String("member_id", "", "Your principal ID")
	flag.Parse()

	// The role to be granted
	var role string = "roles/logging.logWriter"

	// Initializes the Cloud Resource Manager service
	ctx := context.Background()
	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Fatalf("cloudresourcemanager.NewService: %v", err)
	}

	// Grants your principal the "Log writer" role for your project
	addBinding(crmService, *projectID, *member, role)

	// Gets the project's policy and prints all principals with the "Log Writer" role
	policy := getPolicy(crmService, *projectID)
	// Find the policy binding for role. Only one binding can have the role.
	var binding *cloudresourcemanager.Binding
	for _, b := range policy.Bindings {
		if b.Role == role {
			binding = b
			break
		}
	}
	fmt.Println("Role: ", binding.Role)
	fmt.Print("Members: ", strings.Join(binding.Members, ", "))

	// Removes member from the "Log writer" role
	removeMember(crmService, *projectID, *member, role)

}

// addBinding adds the principal to the project's IAM policy
func addBinding(crmService *cloudresourcemanager.Service, projectID, member, role string) {

	policy := getPolicy(crmService, projectID)

	// Find the policy binding for role. Only one binding can have the role.
	var binding *cloudresourcemanager.Binding
	for _, b := range policy.Bindings {
		if b.Role == role {
			binding = b
			break
		}
	}

	if binding != nil {
		// If the binding exists, adds the principal to the binding
		binding.Members = append(binding.Members, member)
	} else {
		// If the binding does not exist, adds a new binding to the policy
		binding = &cloudresourcemanager.Binding{
			Role:    role,
			Members: []string{member},
		}
		policy.Bindings = append(policy.Bindings, binding)
	}

	setPolicy(crmService, projectID, policy)

}

// removeMember removes the principal from the project's IAM policy
func removeMember(crmService *cloudresourcemanager.Service, projectID, member, role string) {

	policy := getPolicy(crmService, projectID)

	// Find the policy binding for role. Only one binding can have the role.
	var binding *cloudresourcemanager.Binding
	var bindingIndex int
	for i, b := range policy.Bindings {
		if b.Role == role {
			binding = b
			bindingIndex = i
			break
		}
	}

	// Order doesn't matter for bindings or members, so to remove, move the last item
	// into the removed spot and shrink the slice.
	if len(binding.Members) == 1 {
		// If the principal is the only member in the binding, removes the binding
		last := len(policy.Bindings) - 1
		policy.Bindings[bindingIndex] = policy.Bindings[last]
		policy.Bindings = policy.Bindings[:last]
	} else {
		// If there is more than one member in the binding, removes the principal
		var memberIndex int
		for i, mm := range binding.Members {
			if mm == member {
				memberIndex = i
			}
		}
		last := len(policy.Bindings[bindingIndex].Members) - 1
		binding.Members[memberIndex] = binding.Members[last]
		binding.Members = binding.Members[:last]
	}

	setPolicy(crmService, projectID, policy)

}

// getPolicy gets the project's IAM policy
func getPolicy(crmService *cloudresourcemanager.Service, projectID string) *cloudresourcemanager.Policy {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	request := new(cloudresourcemanager.GetIamPolicyRequest)
	policy, err := crmService.Projects.GetIamPolicy(projectID, request).Do()
	if err != nil {
		log.Fatalf("Projects.GetIamPolicy: %v", err)
	}

	return policy
}

// setPolicy sets the project's IAM policy
func setPolicy(crmService *cloudresourcemanager.Service, projectID string, policy *cloudresourcemanager.Policy) {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	request := new(cloudresourcemanager.SetIamPolicyRequest)
	request.Policy = policy
	policy, err := crmService.Projects.SetIamPolicy(projectID, request).Do()
	if err != nil {
		log.Fatalf("Projects.SetIamPolicy: %v", err)
	}
}

// [END iam_quickstart]
