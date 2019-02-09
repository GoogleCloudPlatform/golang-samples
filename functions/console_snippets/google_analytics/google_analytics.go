// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_firebase_analytics]

// Package p contains a Google Analytics for Firebase Cloud Function.
package p

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/functions/metadata"
)

// AnalyticsEvent is the payload of an Analytics log event.
// Please refer to the docs for additional information
// regarding Analytics events.
type AnalyticsEvent struct {
	EventDimensions []EventDimensions `json:"eventDim"`
	UserDimensions  interface{}       `json:"userDim"`
}

// EventDimensions holds Analytics event dimensions.
type EventDimensions struct {
	Name                    string      `json:"name"`
	Date                    string      `json:"date"`
	TimestampMicros         string      `json:"timestampMicros"`
	PreviousTimestampMicros string      `json:"previousTimestampMicros"`
	Params                  interface{} `json:"params"`
}

// HelloAnalytics handles Firebase Mobile Analytics log events.
func HelloAnalytics(ctx context.Context, e AnalyticsEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by Google Analytics event: %v", meta.Resource)
	log.Printf("%+v", e)
	return nil
}

// [END functions_firebase_analytics]
