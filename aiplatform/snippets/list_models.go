// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package snippets

// [START aiplatform_list_models]

import (
	"context"
	"fmt"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// listAllModels lists all of the Vertex AI models in a project for
// a specific region.
func listAllModels(projectID string, region string) error {

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", region)
	clientOption := option.WithEndpoint(apiEndpoint)

	ctx := context.Background()
	aiplatformService, err := aiplatform.NewModelClient(ctx, clientOption)
	if err != nil {
		return err
	}
	defer aiplatformService.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)

	req := &aiplatformpb.ListModelsRequest{
		Parent: parent,
	}

	it := aiplatformService.ListModels(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}
		fmt.Printf("Model: %s\n", resp.GetName())
	}
	return nil
}

// [END aiplatform_list_models]
