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
    aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1beta1"
)

func main() {
    // TODO(developer): Replace these variables before running the sample.
    project := "YOUR_PROJECT_ID"
    instance :=
        `{ "prefix": "Write a unit test for this function:\n"
            + "    def is_leap_year(year):\n"
            + "        if year % 4 == 0:\n"
            + "            if year % 100 == 0:\n"
            + "                if year % 400 == 0:\n"
            + "                    return True\n"
            + "                else:\n"
            + "                    return False\n"
            + "            else:\n"
            + "                return True\n"
            + "        else:\n"
            + "            return False\n"
            + "\"}`
    parameters := `{\n"temperature": 0.5,\n"maxOutputTokens": 256\n}`
    location := "us-central1"
    publisher := "google"
    model := "code-bison@001"

    ctx := context.Background()
    client, err := aiplatform.NewPredictionClient(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    endpoint, err := client.GetEndpoint(ctx, &aiplatformpb.GetEndpointRequest{
        Name: fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", project, location, model),
    })
    if err != nil {
        log.Fatal(err)
    }

    req := &aiplatformpb.PredictRequest{
        Endpoint:      endpoint.Name,
        Instances:     []*aiplatformpb.Value{stringToValue(instance)},
        Parameters:    stringToValue(parameters),
        OutputEncoding: aiplatformpb.PredictRequest_JSON,
    }
    resp, err := client.Predict(ctx, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Predict Response")
    fmt.Println(resp)
}

// stringToValue converts a string to a Value.
func stringToValue(s string) *aiplatformpb.Value {
    v := &aiplatformpb.Value{}
    err := json.Unmarshal([]byte(s), v)
    if err != nil {
        log.Fatal(err)
    }
    return v
}
