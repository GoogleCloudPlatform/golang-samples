// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_tips_infinite_retries]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/functions/metadata"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// FiniteRetryPubSub demonstrates how to avoid inifinite retries.
func FiniteRetryPubSub(ctx context.Context, m PubSubMessage) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		// Assume an error on the function invoker and try again.
		return fmt.Errorf("metadata.FromContext: %v", err)
	}

	// Ignore events that are too old.
	expiration := meta.Timestamp.Add(10 * time.Second)
	if time.Now().After(expiration) {
		log.Printf("event timeout: halting retries for expired event '%q'", meta.EventID)
		return nil
	}

	// Add your message processing logic.
	return processTheMessage(m)
}

// [END functions_tips_infinite_retries]

// processTheMessage is a stand-in for the domain logic of your function.
func processTheMessage(m PubSubMessage) error {
	return fmt.Errorf("processTheMessage: runtime error")
}
