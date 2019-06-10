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

// Package background contains a Cloud Function to translate text.
// The function listens to Pub/Sub, does the translations, and stores the
// result in Firestore.
package background

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// A Translation contains the original and translated text.
type Translation struct {
	Original         string `json:"original"`
	Translated       string `json:"translated"`
	OriginalLanguage string `json:"original_language"`
	Language         string `json:"language"`
}

// Clients reused between function invocations.
var (
	translateClient *translate.Client
	firestoreClient *firestore.Client
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Translate translates the given message to French.
func Translate(ctx context.Context, m PubSubMessage) error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return fmt.Errorf("GOOGLE_CLOUD_PROJECT must be set")
	}

	if translateClient == nil {
		// Pre-declare err to avoid shadowing translateClient.
		var err error
		// Use context.Background() so the client can be reused.
		translateClient, err = translate.NewClient(context.Background())
		if err != nil {
			return fmt.Errorf("translate.NewClient: %v", err)
		}
	}
	if firestoreClient == nil {
		// Pre-declare err to avoid shadowing firestoreClient.
		var err error
		// Use context.Background() so the client can be reused.
		firestoreClient, err = firestore.NewClient(context.Background(), projectID)
		if err != nil {
			return fmt.Errorf("firestore.NewClient: %v", err)
		}
	}

	t := Translation{}
	if err := json.Unmarshal(m.Data, &t); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	lang, err := language.Parse(t.Language)
	if err != nil {
		return fmt.Errorf("language.Parse: %v", err)
	}

	outs, err := translateClient.Translate(ctx, []string{t.Original}, lang, nil)
	if err != nil {
		return fmt.Errorf("Translate: %v", err)
	}

	if len(outs) < 1 {
		return fmt.Errorf("Translate got %d translations, need at least 1", len(outs))
	}

	t.Translated = outs[0].Text
	t.OriginalLanguage = outs[0].Source.String()

	if _, err := firestoreClient.Collection("translations").NewDoc().Create(ctx, t); err != nil {
		return fmt.Errorf("Create: %v", err)
	}

	return nil
}
