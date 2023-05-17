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

// [START getting_started_background_translate_setup]

// Package background contains a Cloud Function to translate text.
// The function listens to Pub/Sub, does the translations, and stores the
// result in Firestore.
package background

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// PubSubMessage is the payload of a Pub/Sub event.
// See https://cloud.google.com/functions/docs/calling/pubsub.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// [END getting_started_background_translate_setup]

// [START getting_started_background_translate_init]

// initializeClients creates translateClient and firestoreClient if they haven't
// been created yet.
func initializeClients() error {
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
			return fmt.Errorf("translate.NewClient: %w", err)
		}
	}
	if firestoreClient == nil {
		// Pre-declare err to avoid shadowing firestoreClient.
		var err error
		// Use context.Background() so the client can be reused.
		firestoreClient, err = firestore.NewClient(context.Background(), projectID)
		if err != nil {
			return fmt.Errorf("firestore.NewClient: %w", err)
		}
	}
	return nil
}

// [END getting_started_background_translate_init]

// [START getting_started_background_translate_string]

// translateString translates text to lang, returning:
// * the translated text,
// * the automatically detected source language, and
// * an error.
func translateString(ctx context.Context, text string, lang string) (translated string, originalLang string, err error) {
	l, err := language.Parse(lang)
	if err != nil {
		return "", "", fmt.Errorf("language.Parse: %w", err)
	}

	outs, err := translateClient.Translate(ctx, []string{text}, l, nil)
	if err != nil {
		return "", "", fmt.Errorf("Translate: %w", err)
	}

	if len(outs) < 1 {
		return "", "", fmt.Errorf("Translate got %d translations, need at least 1", len(outs))
	}

	return outs[0].Text, outs[0].Source.String(), nil
}

// [END getting_started_background_translate_string]

// [START getting_started_background_translate]

// Translate translates the given message and stores the result in Firestore.
func Translate(ctx context.Context, m PubSubMessage) error {
	initializeClients()

	t := Translation{}
	if err := json.Unmarshal(m.Data, &t); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	// Use a unique document name to prevent duplicate translations.
	key := fmt.Sprintf("%s/%s", t.Language, t.Original)
	sum := sha512.Sum512([]byte(key))
	// Base64 encode the sum to make a nice string. The [:] converts the byte
	// array to a byte slice.
	docName := base64.StdEncoding.EncodeToString(sum[:])
	// The document name cannot contain "/".
	docName = strings.Replace(docName, "/", "-", -1)
	ref := firestoreClient.Collection("translations").Doc(docName)

	// Run in a transation to prevent concurrent duplicate translations.
	err := firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil && status.Code(err) != codes.NotFound {
			return fmt.Errorf("Get: %w", err)
		}
		// Do nothing if the document already exists.
		if doc.Exists() {
			return nil
		}

		translated, originalLang, err := translateString(ctx, t.Original, t.Language)
		if err != nil {
			return fmt.Errorf("translateString: %w", err)
		}
		t.Translated = translated
		t.OriginalLanguage = originalLang

		if err := tx.Set(ref, t); err != nil {
			return fmt.Errorf("Set: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("RunTransaction: %w", err)
	}
	return nil
}

// [END getting_started_background_translate]
