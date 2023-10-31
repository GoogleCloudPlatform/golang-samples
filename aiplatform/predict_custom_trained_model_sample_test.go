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
    "encoding/base64"
    "fmt"
    "io/ioutil"
    "os"
    "testing"

    aiplatform "cloud.google.com/go/aiplatform/apiv1"
    "github.com/golang/protobuf/ptypes/wrappers"
    "google.golang.org/api/option"
)

var (
    projectId = os.Getenv("UCAIP_PROJECT_ID")
    endpointId = os.Getenv("PREDICT_CUSTOM_TRAINED_MODEL_ENDPOINT_ID")
)

func TestPredictCustomTrainedModel(t *testing.T) {
    ctx := context.Background()
    client, err := aiplatform.NewPredictionServiceClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
    if err != nil {
        t.Fatal(err)
    }

    encoded, err := base64.StdEncoding.EncodeToString(ioutil.ReadFile("resources/daisy.jpg"))
    if err != nil {
        t.Fatal(err)
    }

    instance := `[{'image_bytes': {'b64': '` + encoded + `'}, 'key':'0'}]`

    resp, err := client.Predict(ctx, &aiplatform.PredictRequest{
        Endpoint: fmt.Sprintf("projects/%s/endpoints/%s", projectId, endpointId),
        Instances: []*wrappers.StringValue{
            {Value: instance},
        },
    })
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("Predict Custom Trained model Response")
    fmt.Println(resp)
}
