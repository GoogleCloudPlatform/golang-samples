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

// [START functions_cloudevent_firebase_rtdb]

// Package rtdb contains a Cloud Function triggered by a Firebase Realtime Database
// event.
package rtdb

import (
	"context"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/firebase/databasedata"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	functions.CloudEvent("HelloRTDB", HelloRTDB)
}

// HelloRTDB handles changes to a Firebase RTDB.
func HelloRTDB(ctx context.Context, e event.Event) error {
	unmarshalOptions := protojson.UnmarshalOptions{DiscardUnknown: true}

	var data databasedata.ReferenceEventData
	if err := unmarshalOptions.Unmarshal(e.Data(), &data); err != nil {
		return fmt.Errorf("Unmarshal: %w", err)
	}
	log.Printf("Function triggered by change to: %v", e.Source())
	log.Printf("Data: %+v", data.GetData())
	log.Printf("Delta: %+v", data.GetDelta())
	return nil
}

// [END functions_cloudevent_firebase_rtdb]
