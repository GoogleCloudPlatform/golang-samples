// Copyright 2026 Google LLC
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

package kms

import (
	"encoding/base64"
	"fmt"
	"io"

	"filippo.io/hpke"
)

func wrapHPKE(w io.Writer, publicKeyB64 string, targetKeyB64 string) error {
	// publicKeyB64 is the base64-encoded key of the Cloud KMS ImportJob.
	// targetKeyB64 is the key you want to wrap & import, base64-encoded.

	// 1. Decode the ImportJob public key (obtained from Cloud KMS in NIST_PQC format)
	pkBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// 2. Decode the target key material
	targetKey, err := base64.StdEncoding.DecodeString(targetKeyB64)
	if err != nil {
		return fmt.Errorf("failed to decode target key: %w", err)
	}

	// 3. Define the HPKE suite for Cloud KMS quantum-safe import
	// KEM: ML-KEM-768, KDF: HKDF-SHA256, AEAD: AES-256-GCM
	//
	// Note: you can replace MLKEM768 with MLKEM1024 or MLKEM768X25519
	kem := hpke.MLKEM768()
	kdf := hpke.HKDFSHA256()
	aead := hpke.AES256GCM()

	// 4. Load the public key into the KEM
	pub, err := kem.NewPublicKey(pkBytes)
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}

	// 5. Perform One-Shot Seal (HPKE SetupSender + Seal)
	// info is nil as Cloud KMS expect it to be empty.
	// This function automatically returns the concatenation (enc || ciphertext).
	wrappedKey, err := hpke.Seal(pub, kdf, aead, nil, targetKey)
	if err != nil {
		return fmt.Errorf("failed to wrap key: %w", err)
	}

	fmt.Fprintf(w, "Base64 wrappedKey:\n%s\n", base64.StdEncoding.EncodeToString(wrappedKey))
	return nil
}
