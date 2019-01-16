// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_log_stackdriver]

// Package log contains examples for handling Cloud Functions logs.
package log

import (
	"context"
	"log"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ProcessLogEntry processes a Pub/Sub message from Stackdriver.
func ProcessLogEntry(ctx context.Context, m PubSubMessage) error {
	log.Printf("Log entry data: %s", string(m.Data))
	return nil
}

// [END functions_log_stackdriver]
