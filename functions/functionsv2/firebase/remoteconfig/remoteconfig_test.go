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

package remoteconfig

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/firebase/remoteconfigdata"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestHelloRemoteConfig(t *testing.T) {
	tests := []*remoteconfigdata.RemoteConfigEventData{
		{
			UpdateType:    remoteconfigdata.RemoteConfigUpdateType_INCREMENTAL_UPDATE,
			UpdateOrigin:  remoteconfigdata.RemoteConfigUpdateOrigin_CONSOLE,
			VersionNumber: 1,
		},
		{
			UpdateType:    remoteconfigdata.RemoteConfigUpdateType_INCREMENTAL_UPDATE,
			UpdateOrigin:  remoteconfigdata.RemoteConfigUpdateOrigin_CONSOLE,
			VersionNumber: 2,
		},
	}

	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		originalFlags := log.Flags()
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		jsonData, err := protojson.Marshal(test)
		if err != nil {
			t.Fatalf("protojson.Marshal: %v", err)
		}

		e := event.New()
		e.SetDataContentType("application/json")
		e.SetData(e.DataContentType(), jsonData)

		HelloRemoteConfig(context.Background(), e)

		w.Close()
		log.SetOutput(os.Stderr)
		log.SetFlags(originalFlags)

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}

		want := fmt.Sprintf("Update type: INCREMENTAL_UPDATE\nOrigin: CONSOLE\nVersion: %d\n", test.VersionNumber)
		got := string(out)
		if !strings.Contains(got, want) {
			t.Errorf("HelloRemoteConfig(%v) got %q, want to contain UpdateType %q", e, got, want)
		}
	}
}
