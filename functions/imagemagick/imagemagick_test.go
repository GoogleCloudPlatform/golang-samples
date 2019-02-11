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

package imagemagick

import (
	"context"
	"io/ioutil"
	"log"
	"testing"
)

func TestBlurOffensiveImages(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// TODO: use testutil
	t.Skip("convert is not available in test images")
	e := GCSEvent{
		Bucket: "golang-samples-tests",
		Name:   "functions/zombie.jpg",
	}
	ctx := context.Background()

	outputBlob := storageClient.Bucket(e.Bucket).Object("blurred-" + e.Name)
	outputBlob.Delete(ctx) // Ensure the output file doesn't already exist.

	if err := BlurOffensiveImages(ctx, e); err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error: %v", e, err)
	}

	if _, err := outputBlob.Attrs(ctx); err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error when checking output: %v", e, err)
	}
	outputBlob.Delete(ctx)

	// Check proper handling of already-blurred images.
	e.Name = "blurred-" + e.Name
	if err := BlurOffensiveImages(ctx, e); err != nil {
		t.Fatalf("BlurOffensiveImages(%v) got error on already blurred image: %v", e, err)
	}
}
