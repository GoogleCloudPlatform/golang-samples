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

// [START functions_cloudevent_firebase_firestore]

// Package hellofirestore contains a Cloud Event Function triggered by a Cloud Firestore event.
package hellofirestore

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"google.golang.org/protobuf/proto"
)

func init() {
	functions.CloudEvent("helloFirestore", HelloFirestore)
}

// HelloFirestore is triggered by a change to a Firestore document.
func HelloFirestore(ctx context.Context, event event.Event) error {
	var data firestoredata.DocumentEventData

	// If you omit `DiscardUnknown`, protojson.Unmarshal returns an error
	// when encountering a new or unknown field.
	options := proto.UnmarshalOptions{
		DiscardUnknown: true,
	}
	err := options.Unmarshal(event.Data(), &data)

	if err != nil {
		return fmt.Errorf("proto.Unmarshal: %w", err)
	}

	fmt.Printf("Function triggered by change to: %v\n", event.Source())
	fmt.Printf("Old value: %+v\n", data.GetOldValue())
	fmt.Printf("New value: %+v\n", data.GetValue())
	return nil
}

// [END functions_cloudevent_firebase_firestore]
