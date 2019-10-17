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

// Tests for asymmetric keys feature of Cloud Key Management Service: https://cloud.google.com/kms/
package kms

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetPublicKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	v := getTestVariables(tc.ProjectID)

	rsaDecryptPub, err := getAsymmetricPublicKey(v.rsaDecryptPath)
	if err != nil {
		t.Fatalf("getAsymmetricPiblicKey(%s): %v", v.rsaDecryptPath, err)
	}
	_, ok := rsaDecryptPub.(*rsa.PublicKey)
	if ok != true {
		t.Errorf("expected *rsa.PublicKey type")
	}

	rsaSignPub, err := getAsymmetricPublicKey(v.rsaSignPath)
	if err != nil {
		t.Fatalf("getAsymmetricPiblicKey(%s): %v", v.rsaSignPath, err)
	}
	_, ok = rsaSignPub.(*rsa.PublicKey)
	if ok != true {
		t.Errorf("expected *rsa.PublicKey type")
	}
	ecPub, err := getAsymmetricPublicKey(v.ecSignPath)
	if err != nil {
		t.Fatalf("getAsymmetricPiblicKey(%s): %v", v.ecSignPath, err)
	}
	_, ok = ecPub.(*ecdsa.PublicKey)
	if ok != true {
		t.Errorf("expected *ecdsa.PublicKey type")
	}
}

func TestRSAEncryptDecrypt(t *testing.T) {
	tc := testutil.SystemTest(t)
	v := getTestVariables(tc.ProjectID)

	ciphertext, err := encryptRSA(v.rsaDecryptPath, []byte(v.message))
	if err != nil {
		t.Fatalf("encryptRSA(%s, %s): %v", v.rsaDecryptPath, []byte(v.message), err)
	}
	if len(ciphertext) != 256 {
		t.Fatalf("len(ciphertext) = %d; want: %d", len(ciphertext), 256)
	}
	plainBytes, err := decryptRSA(v.rsaDecryptPath, ciphertext)
	if err != nil {
		t.Fatalf("decryptRSA(%s, %s): %v", ciphertext, v.rsaDecryptPath, err)
	}
	if !bytes.Equal(plainBytes, []byte(v.message)) {
		t.Fatalf("decrypted plaintext does not match input message: want %s, got %s", []byte(v.message), plainBytes)
	}
	if bytes.Equal(ciphertext, []byte(v.message)) {
		t.Fatalf("ciphertext and plaintext bytes are identical: %s", ciphertext)
	}
	plaintext := string(plainBytes)
	if plaintext != v.message {
		t.Fatalf("failed to decypt expected plaintext: want %s, got %s", v.message, plaintext)
	}
}

func TestRSASignVerify(t *testing.T) {
	tc := testutil.SystemTest(t)
	v := getTestVariables(tc.ProjectID)

	sig, err := signAsymmetric(v.rsaSignPath, []byte(v.message))
	if err != nil {
		t.Fatalf("signAsymmetric(%s, %s): %v", v.message, v.rsaSignPath, err)
	}
	if len(sig) != 256 {
		t.Errorf("sig length = %d; want: %d", len(sig), 256)
	}
	if err = verifySignatureRSA(v.rsaSignPath, sig, []byte(v.message)); err != nil {
		t.Fatalf("verifySignatureRSA(%s, %s, %s): %v", sig, v.message, v.rsaSignPath, err)
	}
	changed := v.message + "."
	if err = verifySignatureRSA(v.rsaSignPath, sig, []byte(changed)); err == nil {
		t.Errorf("verification for modified message should fail")
	}
}

func TestECSignVerify(t *testing.T) {
	tc := testutil.SystemTest(t)
	v := getTestVariables(tc.ProjectID)

	sig, err := signAsymmetric(v.ecSignPath, []byte(v.message))
	if err != nil {
		t.Fatalf("signAsymmetric(%s, %s): %v", v.message, v.ecSignPath, err)
	}
	if len(sig) < 50 || len(sig) > 300 {
		t.Errorf("Length = %d; want between 50-300", len(sig))
	}

	if err = verifySignatureEC(v.ecSignPath, sig, []byte(v.message)); err != nil {
		t.Fatalf("verifySignatureEC(%s, %s, %s): %v", sig, v.message, v.ecSignPath, err)
	}
	changed := v.message + "."
	if err = verifySignatureEC(v.ecSignPath, sig, []byte(changed)); err == nil {
		t.Errorf("verification for modified message should fail")
	}
}
