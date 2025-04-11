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
		return fmt.Errorf("metadata.FromContext: %w", err)
	}
	log.Printf("Function triggered by Google Analytics event: %v", meta.Resource)
	log.Printf("%+v", e)
	return nil
}

// [END functions_firebase_analytics]
