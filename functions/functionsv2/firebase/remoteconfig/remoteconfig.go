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

// [START functions_cloudevent_firebase_remote_config]

// Package remote_config provides cloud event sample for firebase remote config updates.
package remoteconfig

import (
	"context"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/firebase/remoteconfigdata"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	functions.CloudEvent("HelloRemoteConfig", HelloRemoteConfig)
}

// HelloRemoteConfig handles Firebase Remote Config events.
func HelloRemoteConfig(ctx context.Context, e event.Event) error {
	unmarshalOptions := protojson.UnmarshalOptions{DiscardUnknown: true}

	var data remoteconfigdata.RemoteConfigEventData
	if err := unmarshalOptions.Unmarshal(e.Data(), &data); err != nil {
		return fmt.Errorf("UnmarshalTo: %w", err)
	}

	log.Printf("Update type: %+v", data.GetUpdateType())
	log.Printf("Origin: %+v", data.GetUpdateOrigin())
	log.Printf("Version: %+v", data.GetVersionNumber())
	return nil
}

// [END functions_cloudevent_firebase_remote_config]
