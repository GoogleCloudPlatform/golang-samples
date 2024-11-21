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

// [START functions_cloudevent_storage]

// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	functions.CloudEvent("HelloStorage", helloStorage)
}

// helloStorage consumes a CloudEvent message and logs details about the changed object.
func helloStorage(ctx context.Context, e event.Event) error {
	log.Printf("Event ID: %s", e.ID())
	log.Printf("Event Type: %s", e.Type())

	var data storagedata.StorageObjectData
	err := protojson.Unmarshal(e.Data(), &data)

	// Unmarshal() returns an unexported error.
	// Parse the error string to determine whether there is an unknown field.
	if strings.Contains(err.Error(), "unknown field") {
		log.Println(err.Error())
		options := protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}
		err = options.Unmarshal(e.Data(), &data)
	}
	if err != nil {
		return fmt.Errorf("protojson.Unmarshal: %w", err)
	}

	log.Printf("Bucket: %s", data.GetBucket())
	log.Printf("File: %s", data.GetName())
	log.Printf("Metageneration: %d", data.GetMetageneration())
	log.Printf("Created: %s", data.GetTimeCreated().AsTime())
	log.Printf("Updated: %s", data.GetUpdated().AsTime())
	return nil
}

// [END functions_cloudevent_storage]
