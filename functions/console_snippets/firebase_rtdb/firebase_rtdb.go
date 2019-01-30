// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_firebase_rtdb]

// Package p contains a Firebase Realtime Database Cloud Function.
package p

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/functions/metadata"
)

// RTDBEvent is the payload of a RTDB event.
// Please refer to the docs for additional information
// regarding Firestore events.
type RTDBEvent struct {
	Data  interface{} `json:"data"`
	Delta interface{} `json:"delta"`
}

// HelloRTDB handles changes to a Firebase RTDB.
func HelloRTDB(ctx context.Context, e RTDBEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("%+v", e)
	return nil
}

// [END functions_firebase_rtdb]
