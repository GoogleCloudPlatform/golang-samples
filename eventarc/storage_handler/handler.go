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

// [START eventarc_storage_cloudevent_handler]

package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

// HelloStorage receives and processes a CloudEvent containing StorageObjectData
func HelloStorage(w http.ResponseWriter, r *http.Request) {
	ce, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		log.Printf("failed to parse CloudEvent: %v", err)
		http.Error(w, "Bad Request: expected CloudEvent", http.StatusBadRequest)
		return
	}

	unmarshalOptions := protojson.UnmarshalOptions{DiscardUnknown: true}
	var so storagedata.StorageObjectData
	err = unmarshalOptions.Unmarshal(ce.Data(), &so)
	if err != nil {
		log.Printf("failed to unmarshal: %v", err)
		http.Error(w, "Bad Request: expected Cloud Storage event", http.StatusBadRequest)
		return
	}

	s := fmt.Sprintf("Cloud Storage object changed: %s updated at %s",
		path.Join(so.GetBucket(), so.GetName()),
		so.Updated.AsTime().UTC())
	fmt.Fprintln(w, s)
}

// [END eventarc_storage_cloudevent_handler]
