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
    "google.golang.org/api/iterator"
    "google.golang.org/api/option"
    aiplatform "google.golang.org/genproto/googleapis/cloud/aiplatform/v1"
)

func TestPredictComment(t *testing.T) {
    ctx := context.Background()
    c, err := aiplatform.NewPredictionServiceClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
    if err != nil {
        t.Fatal(err)
    }
    defer c.Close()

    req := &aiplatform.PredictRequest{
        Instance: `{ "prefix": "
def reverse_string(s):
  return s[::-1]
#This function
"}`,
        Parameters: `{
  "temperature": 0.2,
  "maxOutputTokens": 64
}`,
        Model: "projects/PROJECT_ID/locations/us-central1/models/code-gecko@001",
    }

    resp, err := c.Predict(ctx, req)
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("Predict Response")
    fmt.Println(resp)
}
