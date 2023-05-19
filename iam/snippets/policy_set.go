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

// [START iam_set_policy]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/cloudresourcemanager/v1"
)

// setPolicy sets IAM policy for a project.
func setPolicy(w io.Writer, projectID string, policy *cloudresourcemanager.Policy) error {
	ctx := context.Background()

	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return err
	}

	req := &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	}
	resp, err := crmService.Projects.SetIamPolicy(projectID, req).Do()
	if err != nil {
		return err
	}

	policyJson, err := resp.MarshalJSON()
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Policy set: %s\n", policyJson)
	return nil
}

// [END iam_set_policy]
