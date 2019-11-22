// Copyright 2019 Google LLC
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

/*
Package cloudrunci is a testing utility that facilitates end-to-end system testing on Cloud Run.

cloudrunci facilitates management of ephemeral (temporary) Cloud Run services on all Cloud Run platforms.

Example Usage

	package main_test

	import (
		"io/ioutil"
		"log"
		"os"
		"strings"
		"testing"

		"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	)

	// TestMyService shows the simplest approach to create e2e tests with cloudrun-ci.
	// It uses Cloud Run (fully managed) in the us-central1 region.
	func TestMyService(t *testing.T) {
		myService := cloudrunci.NewService("my-service", os.Getenv("GOOGLE_CLOUD_PROJECT"))
		if err := myService.Deploy(); err != nil {
			t.Fatalf("could not deploy %s: %v", myService.Name, err)
		}
		defer myService.Clean()

		resp, err := myService.Request("GET", "/")
		if err != nil {
			t.Errorf("Get: %v", err)
			return
		}
		defer resp.Body.Close()

		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("ioutil.ReadAll: %q", err)
		}

		want := "Hello World!"
		if got := string(got); !strings.Contains(got, want) {
			t.Errorf("got\n----\n%s\n----\nWant to contain:\n----\n%s\n", got, shouldContain)
		}
	}

Configure a pre-built image for testing:

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	myService := cloudrunci.NewService("my-service", projectID)
	myService.Image = "gcr.io/" + projectID + "/my-service:v1.0"

Configure the service to deploy to Cloud Run (fully managed) in the us-east1 region:

	myService := &cloudrunci.Service{
		Name:      "my-service",
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Platform:  cloudrunci.ManagedPlatform{"us-east1"},
	}

Configure the service for deploying to a Cloud Run for Anthos on GKE cluster:

	myService := &cloudrunci.Service{
		Name:      "my-service",
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Platform:  cloudrunci.GKEPlatform{Cluster: "my-cluster", ClusterLocation: "us-central1-c"},
	}

Configure the service for deploying to a Cloud Run for Anthos on VMWare cluster:

	myService := &cloudrunci.Service{
		Name:      "my-service",
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Platform:  cloudrunci.KubernetesPlatform{Kubeconfig: "~/.kubeconfig", Context: "my-cluster"},
	}
*/
package cloudrunci
