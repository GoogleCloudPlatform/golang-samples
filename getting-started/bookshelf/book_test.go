// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package bookshelf

import "testing"

func TestCreatedByName(t *testing.T) {
	b := &Book{
		CreatedByID: "homer",
		CreatedBy:   "Homer Simpson",
	}

	if want, got := b.CreatedBy, b.CreatedByDisplayName(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestAnonymous(t *testing.T) {
	b := &Book{
		CreatedByID: "homer",
		CreatedBy:   "Homer Simpson",
	}
	b.SetCreatorAnonymous()

	if want, got := "Anonymous", b.CreatedByDisplayName(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
