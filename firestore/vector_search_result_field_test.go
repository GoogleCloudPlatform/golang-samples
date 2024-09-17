// Copyright 2024 Google LLC
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

package firestore

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestVectorSearchDistanceResultField(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}

	buf := new(bytes.Buffer)
	if err := vectorSearchDistanceResultField(buf, projectID); err != nil {
		t.Errorf("vectorSearchDistanceResultField: %v", err)
	}

	// Compare console outputs
	got := buf.String()
	want := "Sleepy coffee beans, Distance: 0\n" +
		"Kahawa coffee beans, Distance: 2.449489742783178\n" +
		"Owl coffee beans, Distance: 5.744562646538029\n"
	if !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
