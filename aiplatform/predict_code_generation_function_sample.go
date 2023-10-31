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

    aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
    aiplatform_resources "cloud.google.com/go/aiplatform/apiv1beta1/resources"
    "google.golang.org/api/iterator"
    "google.golang.org/api/option"
)

func main() {
    ctx := context.Background()
    client, err := aiplatform.NewPredictionServiceClient(ctx, option.WithCredentialsFile("/path/to/key.json"))
    if err != nil {
        log.Fatal(err)
    }

    // TODO(developer): Replace this variable before running the sample.
    project := "YOUR_PROJECT_ID"

    // Learn how to create prompts to work with a code model to generate code:
    // https://cloud.google.com/vertex-ai/docs/generative-ai/code/code-generation-prompts
    instance := "{ \"prefix\": \"Write a function that checks if a year is a leap year.\"}"
    parameters := "{\n" + "  \"temperature\": 0.5,\n" + "  \"maxOutputTokens\": 256,\n" + "}"
    location := "us-central1"
    publisher = "google"
    model = "code-bison@001"

    endpoint, err := aiplatform_resources.NewEndpointName(project, location, publisher, model)
    if err != nil {
        log.Fatal(err)
    }

    req := &aiplatform.PredictRequest{
        Endpoint: endpoint.String(),
        Instances: []*aiplatform.Instance{
            {
                Data: map[string]*aiplatform.Value{
                    "instance": {
                        StringValue: instance,
                    },
                },
            },
        },
        Parameters: parameters,
    }

    resp, err := client.Predict(ctx, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Predict Response")
    fmt.Println(resp)
}
