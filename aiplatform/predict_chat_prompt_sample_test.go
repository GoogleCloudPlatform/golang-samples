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
    "io"
    "io/ioutil"
    "os"
    "testing"

    dialog "cloud.google.com/dialogflow/cx/v3"
    dialogpb "google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
)

var (
    projectId  = os.Getenv("UCAIP_PROJECT_ID")
    instanceId = "my-instance"
    publisher  = "google"
    model      = "chat-bison@001"
)

func TestPredictChatPrompt(t *testing.T) {
    ctx := context.Background()
    client, err := dialog.NewAgentsClient(ctx)
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    req := &dialogpb.PredictRequest{
        Session: fmt.Sprintf("projects/%s/locations/global/agents/%s/sessions/123", projectId, instanceId),
        QueryInput: &dialogpb.QueryInput{
            Text: &dialogpb.TextInput{
                Text: "What is the capital of France?",
            },
        },
    }

    resp, err := client.Predict(ctx, req)
    if err != nil {
        t.Fatalf("Failed to predict: %v", err)
    }

    fmt.Println("Predict Response:")
    fmt.Println(resp)
}

func TestMain(m *testing.M) {
    if projectId == "" {
        fmt.Println("UCAIP_PROJECT_ID environment variable must be set")
        os.Exit(1)
    }

    exitCode := m.Run()

    if exitCode == 0 {
        fmt.Println("All tests passed!")
    } else {
        fmt.Println("Some tests failed!")
    }

    os.Exit(exitCode)
}
