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

func TestPredictTestFunction(t *testing.T) {
    ctx := context.Background()
    c, err := aiplatform.NewPredictionServiceClient(ctx)
    if err != nil {
        t.Fatal(err)
    }
    defer c.Close()

    req := &aiplatform.PredictRequest{
        Endpoint: fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", os.Getenv("PROJECT_ID"), "us-central1", "code-gecko"),
        Instance: &aiplatform.Instance{
            Content: `{"prefix": "def reverse_string(s):\n  return s[::-1]\n\ndef test_empty_input_string()"}`,
        },
        Parameters: &aiplatform.PredictRequest_Parameters{
            Temperature: 0.2,
            MaxOutputTokens: 64,
        },
    }

    op, err := c.Predict(ctx, req)
    if err != nil {
        t.Fatal(err)
    }

    resp, err := op.Wait(ctx)
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("Predict Response:")
    if err := jsonpb.Marshal(resp, os.Stdout); err != nil {
        t.Fatal(err)
    }
}
