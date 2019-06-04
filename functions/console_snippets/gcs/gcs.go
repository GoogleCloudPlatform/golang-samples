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

// Package p contains a Google Cloud Storage Cloud Function.
package p

import (
	"context"
	"log"
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// HelloGCS prints a message when a file is changed in a Cloud Storage bucket.
func HelloGCS(ctx context.Context, e GCSEvent) error {
	log.Printf("Processing file: %s", e.Name)
	return nil
}
