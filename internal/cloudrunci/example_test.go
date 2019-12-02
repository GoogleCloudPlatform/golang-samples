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

package cloudrunci_test

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
)

func init() {
	os.Setenv("BUILD_ID", "12345")
}

// ExampleService_Build shows how to build once, deploy many.
func ExampleService_Build() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	version := os.Getenv("BUILD_ID")

	myService := cloudrunci.NewService("my-service", projectID)
	myService.Image = fmt.Sprintf("gcr.io/%s/my-service:%s", projectID, version)
	myService.Dir = "example"
	myService.Build()

	myService.Env = map[string]string{
		"CUSTOM_VARIABLE": "42",
	}
	myService.Deploy()

	myService.Env = map[string]string{
		"CUSTOM_VARIABLE": "88",
	}
	myService.Deploy()
}

// ExampleService_Request shows how to send an HTTP request to a service using defaults.
func ExampleService_Request() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	myService := cloudrunci.NewService("my-service", projectID)
	myService.Dir = "example"
	myService.Deploy()

	resp, err := myService.Request("GET", "/")
	if err != nil {
		fmt.Printf("myService.Get: %q", err)
	}

	fmt.Println(resp.StatusCode)
}

// ExampleService_Request shows how to run a customized HTTP request.
func ExampleService_NewRequest() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	myService := cloudrunci.NewService("my-service", projectID)
	myService.Dir = "example"
	if err := myService.Deploy(); err != nil {
		fmt.Printf("service.Deploy: %q", err)
		return
	}

	req, err := myService.NewRequest("GET", "/")
	if err != nil {
		fmt.Printf("service.NewRequest: %q", err)
		return
	}
	req.Header.Set("Custom-Header", "42")

	myClient := &http.Client{}
	resp, err := myClient.Do(req)
	if err != nil {
		fmt.Printf("http.Client: %q", err)
		return
	}

	fmt.Println(resp.StatusCode)
}
