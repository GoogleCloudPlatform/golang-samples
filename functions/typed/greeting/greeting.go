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

// [START functions_typed_greeting]

// Package greeting provides a set of Cloud Functions samples.
package greeting

import (
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.Typed("Greeting", greeting)
}

// greeting is a Typed Cloud Function.
func greeting(request *GreetingRequest) (*GreetingResponse, error) {
	return &GreetingResponse{
		Message: fmt.Sprintf("Hello %v %v!", request.FirstName, request.LastName),
	}, nil
}

type GreetingRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GreetingResponse struct {
	Message string `json:"message"`
}

// [END functions_typed_greeting]
