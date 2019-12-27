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

package main

// [START vision_set_endpoint]
import (
	"context"
	"fmt"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
)

// setEndpoint changes your endpoint.
func setEndpoint(endpoint string) error {
	// endpoint := "eu-vision.googleapis.com:443"

	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("NewImageAnnotatorClient: %v", err)
	}
	defer client.Close()

	return nil
}

// [END vision_set_endpoint]
