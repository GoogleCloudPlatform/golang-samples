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

package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"

    aiplatform "cloud.google.com/go/aiplatform/apiv1"
    aiplatform_pb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1"
)

// [START aiplatform_predict_custom_trained_model_sample]

func main() {
    ctx := context.Background()

    // TODO(developer): Replace these variables before running the sample.
    projectID := "YOUR_PROJECT_ID"
    endpointID := "YOUR_ENDPOINT_ID"
    instance := `[{"feature_column_a": "value", "feature_column_b": "value"}]`

    client, err := aiplatform.NewPredictionClient(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    req := &aiplatform_pb.PredictRequest{
        Endpoint: fmt.Sprintf("projects/%s/locations/us-central1/endpoints/%s", projectID, endpointID),
        Instances: []*aiplatform_pb.Value{
            {
                ListValue: &aiplatform_pb.ListValue{
                    Values: []*aiplatform_pb.Value{
                        {
                            StringValue: instance,
                        },
                    },
                },
            },
        },
    }

    resp, err := client.Predict(ctx, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Predict Custom Trained model Response")
    fmt.Printf("\tDeployed Model Id: %s\n", resp.GetDeployedModelId())
    fmt.Println("Predictions")
    for _, prediction := range resp.GetPredictions() {
        fmt.Printf("\tPrediction: %s\n", prediction)
    }
}

// [END aiplatform_predict_custom_trained_model_sample]
