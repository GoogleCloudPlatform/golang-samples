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

    // TODO(developer): Replace these variables before running the sample.
    project := "YOUR_PROJECT_ID"
    location := "us-central1"
    publisher = "google"
    model = "chat-bison@001"
    instance := `
        {
           "context":  "My name is Ned. You are my personal assistant. My favorite movies"
            + " are Lord of the Rings and Hobbit.",
           "examples": [ { 
               "input": {"content": "Who do you work for?"},
               "output": {"content": "I work for Ned."}
            },
            { 
               "input": {"content": "What do I like?"},
               "output": {"content": "Ned likes watching movies."}
            }],
           "messages": [\n"
            + "    { 
               "author": "user",
               "content": "Are my favorite movies based on a book series?"
            + "    }]\n"
            + "}`
    parameters := `
        {
          "temperature": 0.3,
          "maxDecodeSteps": 200,
          "topP": 0.8,
          "topK": 40
        }`

    // Create the client.
    client, err := aiplatform.NewPredictionServiceClient(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create the request.
    endpointName := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", project, location, model)
    req := &aipb.PredictRequest{
        Endpoint: endpointName,
        Instances: []*aipb.Value{
            {
                Value: &aipb.Value_StringValue{
                    StringValue: instance,
                },
            },
        },
        Parameters: parameters,
    }

    // Send the request.
    resp, err := client.Predict(ctx, req)
    if err != nil {
        log.Fatal(err)
    }

    // Print the response.
    fmt.Println(resp)
}
