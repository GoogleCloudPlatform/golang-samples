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

// [START functions_firebase_rtdb]

// Package p contains a Cloud Function triggered by a Firebase Realtime Database
// event.
package p

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/functions/metadata"
)

// RTDBEvent is the payload of a RTDB event.
type RTDBEvent struct {
	Data  interface{} `json:"data"`
	Delta interface{} `json:"delta"`
}

// HelloRTDB handles changes to a Firebase RTDB.
func HelloRTDB(ctx context.Context, e RTDBEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %w", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("%+v", e)
	return nil
}

// [END functions_firebase_rtdb]
