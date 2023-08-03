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
// limitations under the License.

package snippets

// [START aiplatform_text_embeddings]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// GenerateEmbeddings creates embeddings from text provided.
func GenerateEmbeddings(w io.Writer, prompt, project, location, publisher, model string) error {
	ctx := context.Background()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	log.Printf("apiendpoint: %s", apiEndpoint)

	c, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		log.Printf("unable to create prediction client: %v", err)
		return err
	}
	defer c.Close()

	// PredictRequest requires an Endpoint, Instances, and Parameters
	// Endpoint
	base := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models", project, location, publisher)
	url := fmt.Sprintf("%s/%s", base, model)

	// Instances: the prompt
	promptValue, err := structpb.NewValue(map[string]interface{}{
		"content": prompt,
	})
	if err != nil {
		log.Printf("unable to convert prompt to Value: %v", err)
	}

	// PredictRequest
	req := &aiplatformpb.PredictRequest{
		Endpoint:  url,
		Instances: []*structpb.Value{promptValue},
	}
	log.Printf("PredictRequest.Endpoint:   %v", req.GetEndpoint())
	log.Printf("PredictRequest.Instances:  %v", req.GetInstances())
	log.Printf("PredictRequest.Parameters: %v", req.GetParameters())

	// PredictResponse
	resp, err := c.Predict(ctx, req)
	if err != nil {
		log.Printf("error in prediction: %v", err)
		return err
	}
	jsonData, err := json.MarshalIndent(resp, "", " ")
	log.Printf("%s\n", jsonData)
	fmt.Fprintf(w, "embeddings generated: %v", resp.Predictions[0])
	return nil
}

// [END aiplatform_text_embeddings]
