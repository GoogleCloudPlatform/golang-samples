// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
