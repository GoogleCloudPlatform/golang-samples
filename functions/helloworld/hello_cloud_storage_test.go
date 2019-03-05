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

// [START functions_storage_unit_test]

package helloworld

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestHelloGCS(t *testing.T) {
	name := "hello_gcs.txt"
	tests := []struct {
		resourceState  string
		metageneration string
		want           string
	}{
		{
			resourceState: "not_exists",
			want:          fmt.Sprintf("File %s deleted.\n", name),
		},
		{
			metageneration: "1",
			want:           fmt.Sprintf("File %s created.\n", name),
		},
		{
			want: fmt.Sprintf("File %s metadata updated.\n", name),
		},
	}

	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		originalFlags := log.Flags()
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		e := GCSEvent{
			Name:           name,
			ResourceState:  test.resourceState,
			Metageneration: test.metageneration,
		}
		HelloGCS(context.Background(), e)

		w.Close()
		log.SetOutput(os.Stderr)
		log.SetFlags(originalFlags)

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}

		if got := string(out); got != test.want {
			t.Errorf("HelloGCS(%+v) = %q, want %q", e, got, test.want)
		}
	}
}

// [END functions_storage_unit_test]
