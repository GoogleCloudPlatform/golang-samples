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
    "os"
    "testing"

    aiplatform "cloud.google.com/go/aiplatform/apiv1"
    "github.com/golang/protobuf/jsonpb"
    "google.golang.org/api/option"
    aiplatform_v1 "google.golang.org/genproto/googleapis/cloud/aiplatform/v1"
)

func TestPredictCodeChat(t *testing.T) {
    ctx := context.Background()
    c, err := aiplatform.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
    if err != nil {
        t.Fatal(err)
    }
    defer c.Close()

    project := os.Getenv("UCAIP_PROJECT_ID")
    location := "us-central1"
    publisher = "google"
    model = "codechat-bison@001"

    instance, err := ioutil.ReadFile("instance.json")
    if err != nil {
        t.Fatal(err)
    }

    parameters, err := ioutil.ReadFile("parameters.json")
    if err != nil {
        t.Fatal(err)
    }

    req := &aiplatform_v1.PredictRequest{
        Instance:       instance,
        Parameters:     parameters,
        Model:          fmt.Sprintf("projects/%s/locations/%s/models/%s", project, location, publisher+"/"+model),
    }

    resp, err := c.Predict(ctx, req)
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("Predict Response:")
    fmt.Println(resp)

    if resp.Predictions[0].GetError() != nil {
        t.Fatal(resp.Predictions[0].GetError())
    }

    unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
    var prediction aiplatform_v1.TextSegment
    if err := unmarshaler.Unmarshal(resp.Predictions[0].GetOutput(), &prediction); err != nil {
        t.Fatal(err)
    }

    fmt.Println("Prediction:")
    fmt.Println(prediction)
}
