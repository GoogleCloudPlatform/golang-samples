// Copyright 2021 Google LLC
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

// [START functions_cloudevent_storage_unit_test]

package helloworld

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
)

func TestHelloStorage(t *testing.T) {
	r, w, _ := os.Pipe()
	log.SetOutput(w)
	originalFlags := log.Flags()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	d := StorageObjectData{
		Name:        "hello_gcs.txt",
		TimeCreated: time.Now(),
	}
	e := event.New()
	e.SetDataContentType("application/json")
	e.SetData(e.DataContentType(), d)

	HelloStorage(context.Background(), e)

	w.Close()
	log.SetOutput(os.Stderr)
	log.SetFlags(originalFlags)

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	got := string(out)
	if want := d.Name; strings.Contains(want, got) {
		t.Errorf("HelloStorage = %q, want to contain %q", got, want)
	}
	if want := d.TimeCreated.String(); strings.Contains(want, got) {
		t.Errorf("HelloStorage = %q, want to contain %q", got, want)
	}
}

// [END functions_cloudevent_storage_unit_test]
