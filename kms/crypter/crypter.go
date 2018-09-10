// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command crypter encrypts and decrypts a file.
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func main() {
	if len(os.Args) < 8 {
		log.Fatal("usage: go run crypter.go {encrypt,decrypt} PROJECT LOCATION KEYRING CRYPTOKEY INFILE OUTFILE")
	}
	var (
		command     = os.Args[1]
		projectID   = os.Args[2]
		locationID  = os.Args[3]
		keyRingID   = os.Args[4]
		cryptoKeyID = os.Args[5]
		inPath      = os.Args[6]
		outPath     = os.Args[7]
	)

	input, err := ioutil.ReadFile(inPath)
	if err != nil {
		log.Fatalf("Error reading file %q: %v", inPath, err)
	}

	var output []byte
	switch command {
	case "encrypt":
		output, err = encrypt(projectID, locationID, keyRingID, cryptoKeyID, input)
		if err != nil {
			log.Fatalf("Error while encrypting: %v", err)
		}
	case "decrypt":
		output, err = decrypt(projectID, locationID, keyRingID, cryptoKeyID, input)
		if err != nil {
			log.Fatalf("Error while decrypting: %v", err)
		}
	default:
		log.Fatalf("Invalid command: %s. Must be 'encrypt' or 'decrypt'.", command)
	}

	if err := ioutil.WriteFile(outPath, output, 0600); err != nil {
		log.Fatalf("Error writing to file %q: %v", outPath, err)
	}
}

func encrypt(projectID, locationID, keyRingID, cryptoKeyID string, plaintext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, locationID, keyRingID, cryptoKeyID)

	req := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString(plaintext),
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(parentName, req).Do()
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resp.Ciphertext)
}

func decrypt(projectID, locationID, keyRingID, cryptoKeyID string, ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, locationID, keyRingID, cryptoKeyID)

	req := &cloudkms.DecryptRequest{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parentName, req).Do()
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(resp.Plaintext)
}
