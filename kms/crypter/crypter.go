// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command crypter encrypts and decrypts a file.
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1beta1"
)

func main() {
	if len(os.Args) < 7 {
		log.Fatal("usage: go run crypter.go {encrypt,decrypt} PROJECTID KEYRING CRYPTOKEY INFILE OUTFILE")
	}
	var (
		command   = os.Args[1]
		projectID = os.Args[2]
		keyRing   = os.Args[3]
		cryptoKey = os.Args[4]
		inPath    = os.Args[5]
		outPath   = os.Args[6]
	)

	input, err := ioutil.ReadFile(inPath)
	if err != nil {
		log.Fatalf("Error reading file %q: %v", inPath, err)
	}

	var output []byte
	switch command {
	case "encrypt":
		output, err = encrypt(projectID, keyRing, cryptoKey, input)
		if err != nil {
			log.Fatalf("Error while encrypting: %v", err)
		}
	case "decrypt":
		output, err = decrypt(projectID, keyRing, cryptoKey, input)
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

func encrypt(projectID, keyRing, cryptoKey string, plainText []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		log.Fatal(err)
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, "global", keyRing, cryptoKey)

	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.
		Encrypt(parentName, &cloudkms.EncryptRequest{
			Plaintext: base64.StdEncoding.EncodeToString(plainText),
		}).Do()
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resp.Ciphertext)
}

func decrypt(projectID, keyRing, cryptoKey string, cipherText []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		log.Fatal(err)
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, "global", keyRing, cryptoKey)

	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.
		Decrypt(parentName, &cloudkms.DecryptRequest{
			Ciphertext: base64.StdEncoding.EncodeToString(cipherText),
		}).Do()
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(resp.Plaintext)
}
