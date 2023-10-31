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

    // Learn how to create prompts to work with a code model to create code completion suggestions:
    // https://cloud.google.com/vertex-ai/docs/generative-ai/code/code-completion-prompts
    instance :=
        "{ \"prefix\": \""
            + "def reverse_string(s):\n"
            + "  return s[::-1]\n"
            + "#This function"
            + "\"}"
    parameters := "{\n" + "  \"temperature\": 0.2,\n" + "  \"maxOutputTokens\": 64,\n" + "}"
    location := "us-central1"
    publisher = "google"
    model := "code-gecko@001"

    predictComment(ctx, client, instance, parameters, project, location, publisher, model)
}

// Use Codey for Code Completion to complete a code comment
func predictComment(
    ctx context.Context,
    client *aiplatform.PredictionServiceClient,
    instance string,
    parameters string,
    project string,
    location string,
    publisher string,
    model string) {
    endpointName := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", project, location, model)

    instanceValue, err := aiplatform_resources.ParseValue(instance)
    if err != nil {
        log.Fatal(err)
    }
    instances := []*aiplatform_resources.Value{instanceValue}

    parameterValue, err := aiplatform_resources.ParseValue(parameters)
    if err != nil {
        log.Fatal(err)
    }

    predictRequest := &aiplatform.PredictRequest{
        Endpoint: endpointName,
        Instances: instances,
        Parameters: parameterValue,
    }

    predictResponse, err := client.Predict(ctx, predictRequest)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Predict Response")
    fmt.Println(predictResponse)
}
