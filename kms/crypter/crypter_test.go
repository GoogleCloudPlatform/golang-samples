// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCrypter(t *testing.T) {
	tc := testutil.SystemTest(t)

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")
	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	plaintext := []byte("Hello KMS")
	ciphertext, err := encrypt(tc.ProjectID, "global", keyRingID, cryptoKeyID, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Errorf("Ciphertext is the same as plaintext")
	}

	gotPlaintext, err := decrypt(tc.ProjectID, "global", keyRingID, cryptoKeyID, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotPlaintext, plaintext) {
		t.Errorf("decrypt: got %q; want %q", string(gotPlaintext), string(plaintext))
	}
}
