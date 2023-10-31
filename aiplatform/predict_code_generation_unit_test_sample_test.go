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

func TestPredictUnitTest(t *testing.T) {
    ctx := context.Background()
    c, err := aiplatform.NewPredictionServiceClient(ctx)
    if err != nil {
        t.Fatal(err)
    }
    defer c.Close()

    req := &aiplatform.PredictRequest{
        Instance: INSTANCE,
        Parameters: PARAMETERS,
    }

    op, err := c.Predict(ctx, req)
    if err != nil {
        t.Fatal(err)
    }

    err = op.Wait(ctx)
    if err != nil {
        t.Fatal(err)
    }

    resp, err := op.GetResult()
    if err != nil {
        t.Fatal(err)
    }

    fmt.Println("Predict Response:")
    if err := jsonpb.Marshal(resp, os.Stdout); err != nil {
        t.Fatal(err)
    }
}

var (
    PROJECT = os.Getenv("UCAIP_PROJECT_ID")
    INSTANCE = `{ "prefix": "Write a unit test for this function:\n"
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
    PARAMETERS = `{"temperature": 0.5, "maxOutputTokens": 256}`
    PUBLISHER = "google"
    LOCATION = "us-central1"
    MODEL = "code-bison@001"
)

func requireEnvVar(t *testing.T, varName string) {
    if val := os.Getenv(varName); val == "" {
        t.Fatalf("Environment variable '%s' is required to perform these tests.", varName)
    }
}

func TestMain(m *testing.M) {
    requireEnvVar(m, "GOOGLE_APPLICATION_CREDENTIALS")
    requireEnvVar(m, "UCAIP_PROJECT_ID")
    os.Exit(m.Run())
}
