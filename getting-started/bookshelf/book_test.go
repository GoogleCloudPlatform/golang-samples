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

	if got, want := b.CreatedByDisplayName(), b.CreatedBy; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAnonymous(t *testing.T) {
	b := &Book{
		CreatedByID: "homer",
		CreatedBy:   "Homer Simpson",
	}
	b.SetCreatorAnonymous()

	if got, want := b.CreatedByDisplayName(), "Anonymous"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
