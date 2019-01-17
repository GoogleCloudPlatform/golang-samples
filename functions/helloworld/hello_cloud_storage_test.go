// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
