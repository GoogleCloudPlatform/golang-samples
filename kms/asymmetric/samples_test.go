// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Tests for asymmetric keys feature of Cloud Key Management Service: https://cloud.google.com/kms/
package samples

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudkms/v1"
)

type TestVariables struct {
	client         *cloudkms.Service
	ctx            context.Context
	message        string
	rsaDecryptPath string
	rsaSignPath    string
	ecSignPath     string
	rsaDecryptId   string
	rsaSignId      string
	ecSignId       string
	keyRing        string
}

func getTestVariables(projectID string) (TestVariables, error) {
	var v TestVariables
	parent := "projects/" + projectID + "/locations/global"
	keyRing := "kms-asymmetric-sample"

	rsaDecryptId := "rsa-decrypt"
	rsaSignId := "rsa-sign"
	ecSignId := "ec-sign"

	rsaDecrypt := parent + "/keyRings/" + keyRing + "/cryptoKeys/" + rsaDecryptId + "/cryptoKeyVersions/1"
	rsaSign := parent + "/keyRings/" + keyRing + "/cryptoKeys/" + rsaSignId + "/cryptoKeyVersions/1"
	ecSign := parent + "/keyRings/" + keyRing + "/cryptoKeys/" + ecSignId + "/cryptoKeyVersions/1"

	message := "test message 123"

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return v, err
	}
	kmsClient, err := cloudkms.New(client)
	if err != nil {
		return v, err
	}

	v = TestVariables{kmsClient, ctx, message, rsaDecrypt, rsaSign, ecSign, rsaDecryptId, rsaSignId, ecSignId, keyRing}
	return v, nil
}

func createKeyHelper(v TestVariables, keyId, keyPath, purpose, algorithm, parent string) bool {
	if _, err := getAsymmetricPublicKey(v.ctx, v.client, keyPath); err != nil {
		v.client.Projects.Locations.KeyRings.Create(parent, &cloudkms.KeyRing{}).KeyRingId(v.keyRing).Do()
		v.client.Projects.Locations.KeyRings.CryptoKeys.Create(
			parent+"/keyRings/"+v.keyRing, &cloudkms.CryptoKey{
				Purpose: purpose,
				VersionTemplate: &cloudkms.CryptoKeyVersionTemplate{
					Algorithm: algorithm,
				},
			},
		).CryptoKeyId(keyId).Do()
		return true
	}
	return false
}

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	v, err := getTestVariables(tc.ProjectID)
	parent := "projects/" + tc.ProjectID + "/locations/global"
	if ok && err == nil {
		//Create cryptokeys in the test project if needed.
		s1 := createKeyHelper(v, v.rsaDecryptId, v.rsaDecryptPath, "ASYMMETRIC_DECRYPT", "RSA_DECRYPT_OAEP_2048_SHA256", parent)
		s2 := createKeyHelper(v, v.rsaSignId, v.rsaSignPath, "ASYMMETRIC_SIGN", "RSA_SIGN_PSS_2048_SHA256", parent)
		s3 := createKeyHelper(v, v.ecSignId, v.ecSignPath, "ASYMMETRIC_SIGN", "EC_SIGN_P256_SHA256", parent)
		if s1 || s2 || s3 {
			//Leave time for keys to initialize.
			time.Sleep(20 * time.Second)
		}
	}
	//Run tests.
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetPublicKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	v, err := getTestVariables(tc.ProjectID)
	if err != nil {
		t.Fatalf("intial variable setup failed: %v", err)
	}

	rsaDecryptPub, err := getAsymmetricPublicKey(v.ctx, v.client, v.rsaDecryptPath)
	if err != nil {
		t.Fatalf("getAsymmetricPiblicKey(%s): %v", v.rsaDecryptPath, err)
	}
	_, ok := rsaDecryptPub.(*rsa.PublicKey)
	if ok != true {
		t.Errorf("expected *rsa.PublicKey type")
	}

	rsaSignPub, err := getAsymmetricPublicKey(v.ctx, v.client, v.rsaSignPath)
	if err != nil {
		t.Fatalf("getAsymmetricPiblicKey(%s): %v", v.rsaSignPath, err)
	}
	_, ok = rsaSignPub.(*rsa.PublicKey)
	if ok != true {
		t.Errorf("expected *rsa.PublicKey type")
	}
	ecPub, err := getAsymmetricPublicKey(v.ctx, v.client, v.ecSignPath)
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
	v, err := getTestVariables(tc.ProjectID)
	if err != nil {
		t.Fatalf("intial variable setup failed: %v", err)
	}

	cipherBytes, err := encryptRSA(v.ctx, v.client, v.rsaDecryptPath, []byte(v.message))
	ciphertext := base64.StdEncoding.EncodeToString(cipherBytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(ciphertext) != 344 {
		t.Errorf("ciphertext length = %d; want: %d", len(ciphertext), 344)
	}
	if ciphertext[len(ciphertext)-2:] != "==" {
		t.Errorf("ciphertet ending: %s; want: %s", ciphertext[len(ciphertext)-2:], "==")
	}
	plainBytes, err := decryptRSA(v.ctx, v.client, v.rsaDecryptPath, cipherBytes)
	if err != nil {
		t.Fatalf("decryptRSA(%s, %s): %v", ciphertext, v.rsaDecryptPath, err)
	}
	if !bytes.Equal(plainBytes, []byte(v.message)) {
		t.Fatalf("decrypted plaintext does not match input message: want %s, got %s", []byte(v.message), plainBytes)
	}
	plaintext := string(plainBytes)
	if plaintext != v.message {
		t.Fatalf("failed to decypt expected plaintext: want %s, got %s", v.message, plaintext)
	}
}

func TestRSASignVerify(t *testing.T) {
	tc := testutil.SystemTest(t)
	v, err := getTestVariables(tc.ProjectID)
	if err != nil {
		t.Fatalf("intial variable setup failed: %v", err)
	}

	sig, err := signAsymmetric(v.ctx, v.client, v.rsaSignPath, []byte(v.message))
	if err != nil {
		t.Fatalf("signAsymmetric(%s, %s): %v", v.message, v.rsaSignPath, err)
	}
	if len(sig) != 344 {
		t.Errorf("sig length = %d; want: %d", len(sig), 344)
	}
	if sig[len(sig)-2:] != "==" {
		t.Errorf("sig ending: %s; want: %s", sig[len(sig)-2:], "==")
	}
	if err = verifySignatureRSA(v.ctx, v.client, sig, v.rsaSignPath, []byte(v.message)); err != nil {
		t.Fatalf("verifySignatureRSA(%s, %s, %s): %v", sig, v.message, v.rsaSignPath, err)
	}
	changed := v.message + "."
	if err = verifySignatureRSA(v.ctx, v.client, sig, v.rsaSignPath, []byte(changed)); err == nil {
		t.Errorf("verification for modified message should fail")
	}
}

func TestECSignVerify(t *testing.T) {
	tc := testutil.SystemTest(t)
	v, err := getTestVariables(tc.ProjectID)
	if err != nil {
		t.Fatalf("intial variable setup failed: %v", err)
	}

	sig, err := signAsymmetric(v.ctx, v.client, v.ecSignPath, []byte(v.message))
	if err != nil {
		t.Fatalf("signAsymmetric(%s, %s): %v", v.message, v.ecSignPath, err)
	}
	if len(sig) < 50 || len(sig) > 300 {
		t.Errorf("Length = %d; want between 50-300", len(sig))
	}

	if err = verifySignatureEC(v.ctx, v.client, sig, v.ecSignPath, []byte(v.message)); err != nil {
		t.Fatalf("verifySignatureEC(%s, %s, %s): %v", sig, v.message, v.ecSignPath, err)
	}
	changed := v.message + "."
	if err = verifySignatureEC(v.ctx, v.client, sig, v.ecSignPath, []byte(changed)); err == nil {
		t.Errorf("verification for modified message should fail")
	}
}
