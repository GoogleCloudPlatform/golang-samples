// Copyright 2023 Google LLC
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

// Package kms contains samples for asymmetric keys feature of Cloud Key Management Service
// https://cloud.google.com/kms/
package kms

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var fixture *kmsFixture

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Print("skipping - unset GOLANG_SAMPLES_PROJECT_ID?")
		return
	}

	var err error
	fixture, err = NewKMSFixture(tc.ProjectID)
	if err != nil {
		log.Fatalf("failed to create fixture: %s", err)
	}

	exitCode := m.Run()

	if err := fixture.Cleanup(); err != nil {
		log.Fatalf("failed to cleanup resources: %s", err)
	}

	os.Exit(exitCode)
}

func TestCreateKeyAsymmetricDecrypt(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyAsymmetricDecrypt(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeyAsymmetricDecrypt: expected %q to contain %q", got, want)
	}
}

func TestCreateKeyAsymmetricSign(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyAsymmetricSign(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeyAsymmetricSign: expected %q to contain %q", got, want)
	}
}

func TestCreateKeyHSM(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyHSM(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeyHSM: expected %q to contain %q", got, want)
	}
}

func TestCreateKeyLabels(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyLabels(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeyLabels: expected %q to contain %q", got, want)
	}
}

func TestCreateKeyMAC(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyMac(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeyMac: expected %q to contain %q", got, want)
	}
}

func TestCreateKeySymmetricEncryptDecrypt(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeySymmetricEncryptDecrypt(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeySymmetricEncryptDecrypt: expected %q to contain %q", got, want)
	}
}

func TestCreateKeyRotationSchedule(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyRotationSchedule(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeyRotationSchedule: expected %q to contain %q", got, want)
	}
}

func TestCreateKeyVersion(t *testing.T) {
	testutil.SystemTest(t)

	parent := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := createKeyVersion(&b, parent); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key version:"; !strings.Contains(got, want) {
		t.Errorf("createKeyVersion: expected %q to contain %q", got, want)
	}
}

func TestCreateKeyRing(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.LocationName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyRing(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key ring:"; !strings.Contains(got, want) {
		t.Errorf("createKeyRing: expected %q to contain %q", got, want)
	}
}

func TestDecryptAsymmetric(t *testing.T) {
	testutil.SystemTest(t)

	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.AsymmetricDecryptKeyName)

	// Encrypt some data to decrypt.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{
		Name: name,
	})
	if err != nil {
		t.Fatal(err)
	}

	block, _ := pem.Decode([]byte(response.Pem))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("public key is not rsa")
	}

	// Encrypt data using the RSA public key.
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, []byte("fruitloops"), nil)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := decryptAsymmetric(&b, name, ciphertext); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Decrypted plaintext:"; !strings.Contains(got, want) {
		t.Errorf("decryptAsymmetric: expected %q to contain %q", got, want)
	}
}

func TestDecryptSymmetric(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	// Encrypt some data to decrypt.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Encrypt(ctx, &kmspb.EncryptRequest{
		Name:      name,
		Plaintext: []byte("fruitloops"),
	})
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := decryptSymmetric(&b, name, result.Ciphertext); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Decrypted plaintext:"; !strings.Contains(got, want) {
		t.Errorf("decryptSymmetric: expected %q to contain %q", got, want)
	}
}

func TestDestroyRestoreKeyVersion(t *testing.T) {
	testutil.SystemTest(t)

	parent, err := fixture.CreateSymmetricKey(fixture.KeyRingName)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("%s/cryptoKeyVersions/1", parent)

	t.Run("destroy", func(t *testing.T) {
		var b bytes.Buffer
		if err := destroyKeyVersion(&b, name); err != nil {
			t.Fatal(err)
		}

		if got, want := b.String(), "Destroyed key version:"; !strings.Contains(got, want) {
			t.Errorf("destroyKeyVersion: expected %q to contain %q", got, want)
		}

		t.Run("restore", func(t *testing.T) {
			var b bytes.Buffer
			if err := restoreKeyVersion(&b, name); err != nil {
				t.Fatal(err)
			}

			if got, want := b.String(), "Restored key version:"; !strings.Contains(got, want) {
				t.Errorf("restoreKeyVersion: expected %q to contain %q", got, want)
			}
		})
	})
}

func TestDisableEnableKeyVersion(t *testing.T) {
	testutil.SystemTest(t)

	parent, err := fixture.CreateSymmetricKey(fixture.KeyRingName)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("%s/cryptoKeyVersions/1", parent)

	t.Run("disable", func(t *testing.T) {
		var b bytes.Buffer
		if err := disableKeyVersion(&b, name); err != nil {
			t.Fatal(err)
		}

		if got, want := b.String(), "Disabled key version:"; !strings.Contains(got, want) {
			t.Errorf("disableKeyVersion: expected %q to contain %q", got, want)
		}
	})

	t.Run("enable", func(t *testing.T) {
		var b bytes.Buffer
		if err := enableKeyVersion(&b, name); err != nil {
			t.Fatal(err)
		}

		if got, want := b.String(), "Enabled key version:"; !strings.Contains(got, want) {
			t.Errorf("enableKeyVersion: expected %q to contain %q", got, want)
		}
	})
}

func TestEncryptAsymmetric(t *testing.T) {
	testutil.SystemTest(t)

	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.AsymmetricDecryptKeyName)

	var b bytes.Buffer
	if err := encryptAsymmetric(&b, name, "fruitloops"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Encrypted ciphertext:"; !strings.Contains(got, want) {
		t.Errorf("encryptAsymmetric: expected %q to contain %q", got, want)
	}
}

func TestEncryptSymmetric(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := encryptSymmetric(&b, name, "fruitloops"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Encrypted ciphertext:"; !strings.Contains(got, want) {
		t.Errorf("encryptSymmetric: expected %q to contain %q", got, want)
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.LocationName

	var b bytes.Buffer
	if err := generateRandomBytes(&b, name, 256); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Random bytes:"; !strings.Contains(got, want) {
		t.Errorf("generateRandomBytes: expected %q to contain %q", got, want)
	}
}

func TestGetKeyVersionAttestation(t *testing.T) {
	testutil.SystemTest(t)

	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.HSMKeyName)

	var b bytes.Buffer
	if err := getKeyVersionAttestation(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "CAVIUM"; !strings.Contains(got, want) {
		t.Errorf("getKeyVersionAttestation: expected %q to contain %q", got, want)
	}
}

func TestGetKeyLabels(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := getKeyLabels(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "foo=bar"; !strings.Contains(got, want) {
		t.Errorf("getKeyLabels: expected %q to contain %q", got, want)
	}
}

func TestGetPublicKey(t *testing.T) {
	testutil.SystemTest(t)

	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.AsymmetricDecryptKeyName)

	var b bytes.Buffer
	if err := getPublicKey(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Retrieved public key:"; !strings.Contains(got, want) {
		t.Errorf("getPublicKey: expected %q to contain %q", got, want)
	}
}

func TestGetPublicKeyJwk(t *testing.T) {
	testutil.SystemTest(t)

	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.AsymmetricDecryptKeyName)

	var b bytes.Buffer
	if err := getPublicKeyJwk(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "kty"; !strings.Contains(got, want) {
		t.Errorf("getPublicKeyJwk: expected %q to contain %q", got, want)
	}
}

func TestIAMAddMember(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := iamAddMember(&b, name, "group:test@google.com"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated IAM"; !strings.Contains(got, want) {
		t.Errorf("IAMAddMember: expected %q to contain %q", got, want)
	}
}

func TestIAMGetPolicy(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	handle := client.ResourceIAM(name)

	policy, err := handle.Policy(ctx)
	if err != nil {
		t.Fatal(err)
	}
	policy.Add("group:test@google.com", "roles/cloudkms.cryptoKeyEncrypterDecrypter")

	if err := handle.SetPolicy(ctx, policy); err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := iamGetPolicy(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "test@google.com"; !strings.Contains(got, want) {
		t.Errorf("iamGetPolicy: expected %q to contain %q", got, want)
	}
}

func TestIAMRemoveMember(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := iamRemoveMember(&b, name, "group:test@google.com"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated IAM"; !strings.Contains(got, want) {
		t.Errorf("iamRemoveMember: expected %q to contain %q", got, want)
	}
}

func TestImportEndToEnd(t *testing.T) {
	testutil.SystemTest(t)
	var b bytes.Buffer

	// Create key for import.
	cryptoKeyID := fixture.RandomID()
	if err := createKeyForImport(&b, fixture.KeyRingName, cryptoKeyID); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Created key"; !strings.Contains(got, want) {
		t.Fatalf("createKeyForImport: expected %q to contain %q", got, want)
	}
	cryptoKeyName := fmt.Sprintf("%s/cryptoKeys/%s", fixture.KeyRingName, cryptoKeyID)

	// Create import job.
	b.Reset()
	importJobID := fixture.RandomID()
	if err := createImportJob(&b, fixture.KeyRingName, importJobID); err != nil {
		t.Fatal(err)
	}
	if got, want := b.String(), "Created import job"; !strings.Contains(got, want) {
		t.Fatalf("createImportJob: expected %q to contain %q", got, want)
	}
	importJobName := fmt.Sprintf("%s/importJobs/%s", fixture.KeyRingName, importJobID)

	// Check import job state (wait for ACTIVE).
	b.Reset()
	for !strings.Contains(b.String(), "ACTIVE") {
		if err := checkStateImportJob(&b, importJobName); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second * 2)
	}

	// Import the key.
	b.Reset()
	if err := importManuallyWrappedKey(&b, importJobName, cryptoKeyName); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created crypto key version"; !strings.Contains(got, want) {
		t.Fatalf("checkStateImportedKey: expected %q to contain %q", got, want)
	}
	cryptoKeyVersionName := fmt.Sprintf("%s/cryptoKeyVersions/1", cryptoKeyName)

	// Wait for the key to finish importing.
	importInProgressStatus := kmspb.CryptoKeyVersion_CryptoKeyVersionState_name[int32(kmspb.CryptoKeyVersion_PENDING_IMPORT)]
	for {
		b.Reset()
		if err := checkStateImportedKey(&b, cryptoKeyVersionName); err != nil {
			t.Fatal(err)
		}

		got := b.String()
		if want := "Current state"; !strings.Contains(got, want) {
			t.Errorf("checkStateImportedKey: expected %q to contain %q", got, want)
		}

		if !strings.Contains(got, importInProgressStatus) {
			break
		}

		time.Sleep(time.Second * 2)
	}
}

func TestSignAsymmetric(t *testing.T) {
	testutil.SystemTest(t)

	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.AsymmetricSignRSAKeyName)

	var b bytes.Buffer
	if err := signAsymmetric(&b, name, "applejacks"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Signed digest:"; !strings.Contains(got, want) {
		t.Errorf("signAsymmetric: expected %q to contain %q", got, want)
	}
}

func TestSignMac(t *testing.T) {
	testutil.SystemTest(t)

	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.HMACKeyName)

	var b bytes.Buffer
	if err := signMac(&b, name, "fruitloops"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Signature:"; !strings.Contains(got, want) {
		t.Errorf("signMac: expected %q to contain %q", got, want)
	}
}

func TestUpdateKeyUpdateLabels(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := updateKeyUpdateLabels(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "new_label=new_value"; !strings.Contains(got, want) {
		t.Errorf("updateKeyUpdateLabels: expected %q to contain %q", got, want)
	}
}

func TestUpdateKeyAddRotation(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := addRotationSchedule(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated key:"; !strings.Contains(got, want) {
		t.Errorf("addRotationSchedule: expected %q to contain %q", got, want)
	}
}

func TestUpdateKeyRemoveLabels(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := updateKeyRemoveLabels(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated key:"; !strings.Contains(got, want) {
		t.Errorf("updateKeyRemoveLabels: expected %q to contain %q", got, want)
	}
}

func TestUpdateKeyRemoveRotation(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := removeRotationSchedule(&b, name); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated key:"; !strings.Contains(got, want) {
		t.Errorf("removeRotationSchedule: expected %q to contain %q", got, want)
	}
}

func TestUpdateKeySetPrimary(t *testing.T) {
	testutil.SystemTest(t)

	name := fixture.SymmetricKeyName

	var b bytes.Buffer
	if err := updateKeySetPrimary(&b, name, "1"); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Updated key primary:"; !strings.Contains(got, want) {
		t.Errorf("updateKeySetPrimary: expected %q to contain %q", got, want)
	}
}

func TestVerifyAsymmetricEC(t *testing.T) {
	testutil.SystemTest(t)

	message := []byte("applejacks")
	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.AsymmetricSignECKeyName)

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	digest := sha256.Sum256(message)
	result, err := client.AsymmetricSign(ctx, &kmspb.AsymmetricSignRequest{
		Name: name,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest[:],
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := verifyAsymmetricSignatureEC(&b, name, message, result.Signature); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Verified signature"; !strings.Contains(got, want) {
		t.Errorf("verifyAsymmetricEC: expected %q to contain %q", got, want)
	}
}

func TestVerifyAsymmetricRSA(t *testing.T) {
	testutil.SystemTest(t)

	message := []byte("applejacks")
	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.AsymmetricSignRSAKeyName)

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	digest := sha256.Sum256(message)
	result, err := client.AsymmetricSign(ctx, &kmspb.AsymmetricSignRequest{
		Name: name,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest[:],
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := verifyAsymmetricSignatureRSA(&b, name, message, result.Signature); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Verified signature"; !strings.Contains(got, want) {
		t.Errorf("verifyAsymmetricRSA: expected %q to contain %q", got, want)
	}
}

func TestVerifyMac(t *testing.T) {
	testutil.SystemTest(t)

	message := []byte("fruitloops")
	name := fmt.Sprintf("%s/cryptoKeyVersions/1", fixture.HMACKeyName)

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.MacSign(ctx, &kmspb.MacSignRequest{
		Name: name,
		Data: message,
	})
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	if err := verifyMac(&b, name, message, result.Mac); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Verified: true"; !strings.Contains(got, want) {
		t.Errorf("verifyAsymmetricRSA: expected %q to contain %q", got, want)
	}
}
