package main

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudkms/v1"
	"testing"
)

type TestVariables struct {
	client         *cloudkms.Service
	ctx            context.Context
	message        string
	rsaDecryptPath string
	rsaSignPath    string
	ecSignPath     string
}

func setup(t *testing.T) TestVariables {
	tc := testutil.SystemTest(t)
	projectId := tc.ProjectID
	parent := "projects/" + projectId + "/locations/global"
	keyRing := "kms-asymmetric-samples"

	rsaDecrypt := parent + "/keyRings/" + keyRing + "/cryptoKeys/rsa-decrypt/cryptoKeyVersions/1"
	rsaSign := parent + "/keyRings/" + keyRing + "/cryptoKeys/rsa-sign/cryptoKeyVersions/1"
	ecSign := parent + "/keyRings/" + keyRing + "/cryptoKeys/ec-sign/cryptoKeyVersions/1"

	message := "test message 123"

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		t.Fatal(err)
	}
	kmsClient, err := cloudkms.New(client)
	if err != nil {
		t.Fatal(err)
	}
	//create cryptokeys
	kmsClient.Projects.Locations.KeyRings.Create(parent, &cloudkms.KeyRing{}).KeyRingId(keyRing).Do()
	kmsClient.Projects.Locations.KeyRings.CryptoKeys.Create(
		parent+"/keyRings/"+keyRing, &cloudkms.CryptoKey{
			Purpose: "ASYMMETRIC_DECRYPT",
			VersionTemplate: &cloudkms.CryptoKeyVersionTemplate{
				Algorithm: "RSA_DECRYPT_OAEP_2048_SHA256",
			},
		}).CryptoKeyId("rsa-decrypt").Do()
	kmsClient.Projects.Locations.KeyRings.CryptoKeys.Create(
		parent+"/keyRings/"+keyRing, &cloudkms.CryptoKey{
			Purpose: "ASYMMETRIC_SIGN",
			VersionTemplate: &cloudkms.CryptoKeyVersionTemplate{
				Algorithm: "RSA_SIGN_PSS_2048_SHA256",
			},
		}).CryptoKeyId("rsa-sign").Do()
	kmsClient.Projects.Locations.KeyRings.CryptoKeys.Create(
		parent+"/keyRings/"+keyRing, &cloudkms.CryptoKey{
			Purpose: "ASYMMETRIC_SIGN",
			VersionTemplate: &cloudkms.CryptoKeyVersionTemplate{
				Algorithm: "EC_SIGN_P224_SHA256",
			},
		}).CryptoKeyId("ec-sign").Do()

	v := TestVariables{kmsClient, ctx, message, rsaDecrypt, rsaSign, ecSign}
	return v
}

//test equality between two values
func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%s != %s", a, b)
	}
}

func TestGetPublicKey(t *testing.T) {
	v := setup(t)

	rsaDecryptPub, err := getAsymmetricPublicKey(v.client, v.ctx, v.rsaDecryptPath)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := rsaDecryptPub.(*rsa.PublicKey)
	assertEqual(t, ok, true)

	rsaSignPub, err := getAsymmetricPublicKey(v.client, v.ctx, v.rsaSignPath)
	if err != nil {
		t.Fatal(err)
	}
	_, ok = rsaSignPub.(*rsa.PublicKey)
	assertEqual(t, ok, true)

	ecPub, err := getAsymmetricPublicKey(v.client, v.ctx, v.ecSignPath)
	if err != nil {
		t.Fatal(err)
	}
	_, ok = ecPub.(*ecdsa.PublicKey)
	assertEqual(t, ok, true)
}

func TestRSAEncryptDecrypt(t *testing.T) {
	v := setup(t)

	cipherText, err := encryptRSA(v.client, v.ctx, v.message, v.rsaDecryptPath)
	if err != nil {
		t.Fatal(err)
	}
	//cipher text should be 344 characters with base64 and RSA 2048
	assertEqual(t, len(cipherText), 344)
	assertEqual(t, cipherText[len(cipherText)-2:], "==")
	plainText, err := decryptRSA(v.client, v.ctx, cipherText, v.rsaDecryptPath)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, plainText, v.message)
}

func TestRSASignVerify(t *testing.T) {
	v := setup(t)

	sig, err := signAsymmetric(v.client, v.ctx, v.message, v.rsaSignPath)
	if err != nil {
		t.Fatal(err)
	}
	//cipher text should be 344 characters with base64 and RSA 2048
	assertEqual(t, len(sig), 344)
	assertEqual(t, sig[len(sig)-2:], "==")

	err = verifySignatureRSA(v.client, v.ctx, sig, v.message, v.rsaSignPath)
	if err != nil {
		t.Fatal(err)
	}
	err = verifySignatureRSA(v.client, v.ctx, sig, v.message+".", v.rsaSignPath)
	if err == nil {
		t.Errorf("verification for modified message should fail")
	}
}

func TestECSignVerify(t *testing.T) {
	v := setup(t)

	sig, err := signAsymmetric(v.client, v.ctx, v.message, v.ecSignPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(sig) < 50 || len(sig) > 300 {
		t.Errorf("signature length outside expected range. Length: %d", len(sig))
	}

	err = verifySignatureEC(v.client, v.ctx, sig, v.message, v.ecSignPath)
	if err != nil {
		t.Fatal(err)
	}
	err = verifySignatureEC(v.client, v.ctx, sig, v.message+".", v.ecSignPath)
	if err == nil {
		t.Errorf("verification for modified message should fail")
	}
}
