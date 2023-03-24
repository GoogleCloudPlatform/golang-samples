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

package hellofirestore

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/firestoredata"
	"google.golang.org/protobuf/proto"
)

func TestHelloFirestore(t *testing.T) {
	ctx := context.Background()

	// create a Firestore document event
	data := &firestoredata.DocumentEventData{
		OldValue: &firestoredata.Document{
			Name: "oldDocumentName",
		},
		Value: &firestoredata.Document{
			Name: "newDocumentName",
		},
	}
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		t.Fatalf("proto.Marshal: %v", err)
	}

	// create a CloudEvent with the Firestore document event data
	ce := event.New("1.0")
	ce.SetSource("test-source")
	ce.SetType("google.cloud.firestore.document.v1.created")
	ce.SetDataContentType("application/protobuf")
	ce.SetData(ce.DataContentType(), dataBytes)

	// capture the output of fmt.Println
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	// call HelloFirestore with the CloudEvent
	if err := HelloFirestore(ctx, ce); err != nil {
		t.Fatalf("HelloFirestore: %v", err)
	}

	// restore stdout
	w.Close()
	os.Stdout = oldStdout

	// read the output of fmt.Println
	buf := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}

	// verify the log output
	expectedOutput := "Function triggered by change to: test-source\n" +
		"Old value: name:\"oldDocumentName\"\n" +
		"New value: name:\"newDocumentName\"\n"
	if gotOutput := buf.String(); gotOutput != expectedOutput {
		t.Errorf("unexpected log output: got %q, want %q", gotOutput, expectedOutput)
	}
}
