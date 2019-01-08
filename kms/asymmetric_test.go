// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Tests for asymmetric keys feature of Cloud Key Management Service: https://cloud.google.com/kms/
package kms

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
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

	cipherBytes, err := encryptRSA(v.rsaDecryptPath, []byte(v.message))
	ciphertext := base64.StdEncoding.EncodeToString(cipherBytes)
	if err != nil {
		t.Fatalf("encryptRSA(%s, %s): %v", v.rsaDecryptPath, []byte(v.message), err)
	}
	if len(cipherBytes) != 256 {
		t.Fatalf("ciphertext length = %d; want: %d", len(ciphertext), 256)
	}
	plainBytes, err := decryptRSA(v.rsaDecryptPath, cipherBytes)
	if err != nil {
		t.Fatalf("decryptRSA(%s, %s): %v", ciphertext, v.rsaDecryptPath, err)
	}
	if !bytes.Equal(plainBytes, []byte(v.message)) {
		t.Fatalf("decrypted plaintext does not match input message: want %s, got %s", []byte(v.message), plainBytes)
	}
	if bytes.Equal(cipherBytes, []byte(v.message)) {
		t.Fatalf("ciphertext and plaintext bytes are identical: %s", cipherBytes)
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
