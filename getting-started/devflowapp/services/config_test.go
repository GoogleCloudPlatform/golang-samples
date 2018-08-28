// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Unit test for config.go

package services

import (
	"os"
	"testing"
)

func TestHandleCheckMessages(t *testing.T) {
	os.Setenv("MESSAGE_SERVICE", "mock")
	newMessageService() // Fatal error if not successful
}
