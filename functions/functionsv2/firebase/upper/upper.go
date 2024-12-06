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

// [START functions_cloudevent_firebase_reactive]

// Package upper contains a Firestore Cloud Function.
package upper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"google.golang.org/protobuf/proto"
)

// set the GOOGLE_CLOUD_PROJECT environment variable when deploying.
var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
var client *firestore.Client

func init() {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	// Register cloud event function
	functions.CloudEvent("MakeUpperCase", MakeUpperCase)
}

// MakeUpperCase is triggered by a change to a Firestore document. It updates
// the `original` value of the document to upper case.
func MakeUpperCase(ctx context.Context, e event.Event) error {
	var data firestoredata.DocumentEventData

	// If you omit `DiscardUnknown`, protojson.Unmarshal returns an error
	// when encountering a new or unknown field.
	options := proto.UnmarshalOptions{
		DiscardUnknown: true,
	}
	err := options.Unmarshal(e.Data(), &data)

	if err != nil {
		return fmt.Errorf("proto.Unmarshal: %w", err)
	}

	if data.GetValue() == nil {
		return errors.New("Invalid message: 'Value' not present")
	}

	fullPath := strings.Split(data.GetValue().GetName(), "/documents/")[1]
	pathParts := strings.Split(fullPath, "/")
	collection := pathParts[0]
	doc := strings.Join(pathParts[1:], "/")

	var originalStringValue string
	if v, ok := data.GetValue().GetFields()["original"]; ok {
		originalStringValue = v.GetStringValue()
	} else {
		return errors.New("Document did not contain field \"original\"")
	}

	newValue := strings.ToUpper(originalStringValue)
	if originalStringValue == newValue {
		log.Printf("%q is already upper case: skipping", originalStringValue)
		return nil
	}
	log.Printf("Replacing value: %q -> %q", originalStringValue, newValue)

	newDocumentEntry := map[string]string{"original": newValue}
	_, err = client.Collection(collection).Doc(doc).Set(ctx, newDocumentEntry)
	if err != nil {
		return fmt.Errorf("Set: %w", err)
	}
	return nil
}

// [END functions_cloudevent_firebase_reactive]
