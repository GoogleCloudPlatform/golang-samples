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
    aipb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1beta1"
)

func main() {
    ctx := context.Background()

    // TODO(developer): Replace this variable before running the sample.
    project := "YOUR_PROJECT_ID"

    // Learn more about creating prompts to work with a code chat model at:
    // https://cloud.google.com/vertex-ai/docs/generative-ai/code/code-chat-prompts
    instance :=
        "{ \"messages\": [\n"
            + "{\n"
            + "  \"author\": \"user\",\n"
            + "  \"content\": \"Hi, how are you?\"\n"
            + "},\n"
            + "{\n"
            + "  \"author\": \"system\",\n"
            + "  \"content\": \"I am doing good. What can I help you in the coding world?\"\n"
            + " },\n"
            + "{\n"
            + "  \"author\": \"user\",\n"
            + "  \"content\":\n"
            + "     \"Please help write a function to calculate the min of two numbers.\"\n"
            + "}\n"
            + "]}";
    parameters := "{\n" + "  \"temperature\": 0.5,\n" + "  \"maxOutputTokens\": 1024\n" + "}";
    location := "us-central1"
    publisher = "google"
    model := "codechat-bison@001"

    predictCodeChat(ctx, instance, parameters, project, location, publisher, model)
}

// Use a code chat model to generate a code function
func predictCodeChat(
    ctx context.Context,
    instance string,
    parameters string,
    project string,
    location string,
    publisher string,
    model string,
) {
    endpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
    client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(endpoint))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    endpointName := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", project, location, model)
    instanceValue, err := stringToValue(instance)
    if err != nil {
        log.Fatal(err)
    }
    instances := []*aipb.Value{instanceValue}

    parameterValue, err := stringToValue(parameters)
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Predict(ctx, &aipb.PredictRequest{
        Endpoint: endpointName,
        Instances: instances,
        Parameters: parameterValue,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Predict Response")
    fmt.Println(resp)
}

// Convert a Json string to a protobuf.Value
func stringToValue(value string) (*aipb.Value, error) {
    valueBuilder := &aipb.Value{}
    err := jsonpb.UnmarshalString(value, valueBuilder)
    if err != nil {
        return nil, err
    }
    return valueBuilder, nil
}
