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
    aiplatform_v1beta1 "cloud.google.com/go/aiplatform/apiv1beta1"
    "google.golang.org/api/iterator"
    "google.golang.org/api/option"
    "google.golang.org/protobuf/encoding/protojson"
    "google.golang.org/protobuf/types/known/structpb"
)

func main() {
    ctx := context.Background()
    // TODO(developer): Replace this variable before running the sample.
    project := "YOUR_PROJECT_ID"

    // Learn how to create prompts to work with a code model to create code completion suggestions:
    // https://cloud.google.com/vertex-ai/docs/generative-ai/code/code-completion-prompts
    instance :=
        "{ \"prefix\": \""
            + "def reverse_string(s):\n"
            + "  return s[::-1]\n"
            + "def test_empty_input_string()"
            + "}"
    parameters := "{\n" + "  \"temperature\": 0.2,\n" + "  \"maxOutputTokens\": 64,\n" + "}"
    location := "us-central1"
    publisher := "google"
    model := "code-gecko@001"

    predictTestFunction(ctx, instance, parameters, project, location, publisher, model)
}

// Use Codey for Code Completion to complete a test function
func predictTestFunction(
    ctx context.Context,
    instance string,
    parameters string,
    project string,
    location string,
    publisher string,
    model string,
) {
    endpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
    opts := option.WithEndpoint(endpoint)
    client, err := aiplatform.NewPredictionClient(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }

    endpointName := aiplatform_v1beta1.EndpointName(project, location, publisher, model)

    instanceValue, err := structpb.NewValue(instance)
    if err != nil {
        log.Fatal(err)
    }
    instances := []*structpb.Value{instanceValue}

    parameterValue, err := structpb.NewValue(parameters)
    if err != nil {
        log.Fatal(err)
    }

    predictResponse, err := client.Predict(ctx, endpointName, instances, parameterValue)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Predict Response")
    fmt.Println(predictResponse)
}
