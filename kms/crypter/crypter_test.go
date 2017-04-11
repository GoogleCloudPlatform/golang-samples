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

	keyRing := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKey := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")
	if keyRing == "" || cryptoKey == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	plainText := []byte("Hello KMS")
	cipherText, err := encrypt(tc.ProjectID, keyRing, cryptoKey, plainText)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(cipherText, plainText) {
		t.Errorf("Ciphertext is the same as plaintext")
	}

	decryptedText, err := decrypt(tc.ProjectID, keyRing, cryptoKey, cipherText)
	if !bytes.Equal(decryptedText, plainText) {
		t.Errorf("decrypt: got %q; want %q", string(decryptedText), string(plainText))
	}
}
