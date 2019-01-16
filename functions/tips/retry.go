// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_tips_retry]

// Package tips contains tips for writing Cloud Functions in Go.
package tips

import (
	"context"
	"errors"
	"log"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// RetryPubSub demonstrates how to toggle using retries.
func RetryPubSub(ctx context.Context, m PubSubMessage) error {
	name := string(m.Data)
	if name == "" {
		name = "World"
	}

	// A misconfigured client will stay broken until the function is redeployed.
	client, err := MisconfiguredDataClient()
	if err != nil {
		log.Printf("MisconfiguredDataClient (retry denied):  %v", err)
		// A nil return indicates that the function does not need a retry.
		return nil
	}

	// Runtime error might be resolved with a new attempt.
	if err = FailedWriteOperation(client, name); err != nil {
		log.Printf("FailedWriteOperation (retry expected): %v", err)
		// A non-nil return indicates that a retry is needed.
		return err
	}

	return nil
}

// [END functions_tips_retry]

var misconfigured = true

// MisconfiguredDataClient simulates failure to retrieve a data storage client.
func MisconfiguredDataClient() (interface{}, error) {
	if misconfigured {
		return nil, errors.New("access denied to data service")
	}
	return nil, nil
}

// FailedWriteOperation simulates failure to write data via the provided client
func FailedWriteOperation(client interface{}, data string) error {
	return errors.New("runtime error")
}
