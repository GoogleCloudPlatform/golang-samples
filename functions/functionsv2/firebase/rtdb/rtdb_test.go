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

package rtdb

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/firebase/databasedata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestHelloRTDB(t *testing.T) {
	r, w, _ := os.Pipe()
	log.SetOutput(w)
	originalFlags := log.Flags()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	data := databasedata.ReferenceEventData{
		Data:  structpb.NewStringValue("test"),
		Delta: structpb.NewNullValue(),
	}

	jsonData, err := protojson.Marshal(&data)
	if err != nil {
		t.Fatalf("protojson.Marshal: %v", err)
	}

	e := event.New()
	e.SetSource("foo")
	e.SetDataContentType("application/json")
	e.SetData(e.DataContentType(), jsonData)

	HelloRTDB(context.Background(), e)

	w.Close()
	log.SetOutput(os.Stderr)
	log.SetFlags(originalFlags)

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	want := "Function triggered by change to: foo\nData: string_value:\"test\"\nDelta: null_value:NULL_VALUE\n"
	got := string(out)
	if got != want {
		t.Errorf("HelloRTDB(%v) got %q, want %q", e, got, want)
	}
}
