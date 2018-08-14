package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/cloudkms/v1"
	"log"
	"math/big"

	"golang.org/x/oauth2/google"
)

// [START kms_get_asymmetric_public]
// Retrieve a public key from a saved asymmetric key pair on KMS
func getAsymmetricPublicKey(client *cloudkms.Service, ctx context.Context, keyPath string) (interface{}, error) {
	response, err := client.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.
		GetPublicKey(keyPath).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("no response from GetPublicKey")
	}
	keyBytes := []byte(response.Pem)
	block, _ := pem.Decode(keyBytes)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// [END kms_get_asymmetric_public]

// [START kms_decrypt_rsa]
// Attempt to decrypt ciphertext with saved RSA key
func decryptRSA(client *cloudkms.Service, ctx context.Context, ciphertext, keyPath string) (string, error) {
	decryptRequest := &cloudkms.AsymmetricDecryptRequest{
		Ciphertext: ciphertext,
	}
	response, err := client.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.
		AsymmetricDecrypt(keyPath, decryptRequest).Context(ctx).Do()
	if err != nil {
		return "", err
	}
	messageArr, err := base64.StdEncoding.DecodeString(response.Plaintext)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("%s", messageArr)
	return message, nil
}

// [END kms_decrypt_rsa]

// [START kms_encrypt_rsa]
// Encrypt plaintext message using saved RSA public key
func encryptRSA(client *cloudkms.Service, ctx context.Context, message, keyPath string) (string, error) {
	abstractKey, err := getAsymmetricPublicKey(client, ctx, keyPath)
	if err != nil {
		return "", err
	}

	//perform type assertion to get rsa key
	rsaKey := abstractKey.(*rsa.PublicKey)

	ciphertextBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, []byte(message), nil)
	if err != nil {
		return "", err
	}
	ciphertextStr := base64.StdEncoding.EncodeToString(ciphertextBytes)
	return ciphertextStr, nil
}

// [END kms_encrypt_rsa]

// [START kms_sign_asymmetric]
// Sign a message using saved asymmetric private key
func signAsymmetric(client *cloudkms.Service, ctx context.Context, message, keyPath string) (string, error) {
	//find hash of message
	digest := sha256.New()
	digest.Write([]byte(message))
	digestStr := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	asymmetricSignRequest := &cloudkms.AsymmetricSignRequest{
		Digest: &cloudkms.Digest{
			Sha256: digestStr,
		},
	}

	response, err := client.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.
		AsymmetricSign(keyPath, asymmetricSignRequest).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// [END kms_sign_asymmetric]

// [START kms_verify_signature_rsa]
// Verify the cryptographic signature for a message signed with an RSA private key
func verifySignatureRSA(client *cloudkms.Service, ctx context.Context, signature, message, keyPath string) error {
	abstractKey, err := getAsymmetricPublicKey(client, ctx, keyPath)
	if err != nil {
		return err
	}
	//perform type assertion to get rsa key
	rsaKey := abstractKey.(*rsa.PublicKey)
	decodedSignature, _ := base64.StdEncoding.DecodeString(signature)

	digest := sha256.New()
	digest.Write([]byte(message))
	hash := digest.Sum(nil)

	pssOptions := rsa.PSSOptions{SaltLength: len(hash), Hash: crypto.SHA256}
	err = rsa.VerifyPSS(rsaKey, crypto.SHA256, hash, decodedSignature, &pssOptions)
	if err != nil {
		return fmt.Errorf("verification failed")
	}
	return nil
}

// [END kms_verify_signature_rsa]

// [START kms_verify_signature_ec]
// Verify the cryptographic signature for a message signed with an Elliptic Curve private key
func verifySignatureEC(client *cloudkms.Service, ctx context.Context, signature, message, keyPath string) error {
	abstractKey, err := getAsymmetricPublicKey(client, ctx, keyPath)
	if err != nil {
		return err
	}
	//perform type assertion to get elliptic curve key
	ecKey := abstractKey.(*ecdsa.PublicKey)
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)

	var parsedSig struct{ R, S *big.Int }
	rest, err := asn1.Unmarshal(decodedSignature, &parsedSig)
	if err != nil || len(rest) != 0 {
		return err
	}

	digest := sha256.New()
	digest.Write([]byte(message))
	hash := digest.Sum(nil)

	if !ecdsa.Verify(ecKey, hash, parsedSig.R, parsedSig.S) {
		return fmt.Errorf("verification failed")
	}
	return nil
}

// [END kms_verify_signature_ec]

func main() {
	project_id := "sanche-testing-project"
	location := "global"
	keyName := "ec-sign"
	keyring := "test-ring"
	keyVersion := "1"

	keyPath := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s", project_id, location, keyring, keyName, keyVersion)

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	kmsClient, err := cloudkms.New(client)
	if err != nil {
		log.Fatal(err)
	}

	message := "test msg"
	sig, _ := signAsymmetric(kmsClient, ctx, message, keyPath)
	fmt.Println(sig)
	success := verifySignatureEC(kmsClient, ctx, sig, message, keyPath)
	fmt.Println(success)
}
