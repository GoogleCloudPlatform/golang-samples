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

    aiplatform "cloud.google.com/go/aiplatform/apiv1"
    "github.com/golang/protobuf/jsonpb"
    "github.com/google/uuid"
    longrunning "google.golang.org/genproto/googleapis/longrunning"
)

func TestPredictFunction(t *testing.T) {
    ctx := context.Background()
    c, err := aiplatform.NewPredictionServiceClient(ctx)
    if err != nil {
        t.Fatal(err)
    }
    defer c.Close()

    // Generate a unique ID for the request.
    reqID := uuid.New().String()

    // Create the request.
    req := &aiplatform.PredictRequest{
        Endpoint: "projects/" + os.Getenv("PROJECT_ID") + "/locations/us-central1/endpoints/code-bison@001",
        Instance: &aiplatform.Instance{
            Content: `{"prefix": "Write a function that checks if a year is a leap year."}`,
        },
        Parameters: &aiplatform.Parameters{
            Content: `{"temperature": 0.5, "maxOutputTokens": 256}`,
        },
    }

    // Send the request.
    op, err := c.Predict(ctx, req)
    if err != nil {
        t.Fatal(err)
    }

    // Wait for the operation to complete.
    resp, err := op.Wait(ctx)
    if err != nil {
        t.Fatal(err)
    }

    // Print the response.
    if err := jsonpb.Unmarshal(resp.GetPayload(), &resp); err != nil {
        t.Fatal(err)
    }
    fmt.Println(resp)
}
