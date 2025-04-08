// Copyright 2020 Google LLC
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

// This file contains shared test helpers, primarily for creating and destroying
// resources.

package kms

import (
	"context"
	"fmt"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/gofrs/uuid"
	"google.golang.org/api/iterator"
	fieldmask "google.golang.org/genproto/protobuf/field_mask"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type kmsFixture struct {
	client *kms.KeyManagementClient

	ProjectID                string
	LocationName             string
	KeyRingName              string
	AsymmetricDecryptKeyName string
	AsymmetricSignECKeyName  string
	AsymmetricSignRSAKeyName string
	HSMKeyName               string
	SymmetricKeyName         string
	HMACKeyName              string
}

func NewKMSFixture(projectID string) (*kmsFixture, error) {
	var k kmsFixture
	var err error

	k.client, err = kms.NewKeyManagementClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %w", err)
	}

	k.ProjectID = projectID

	k.LocationName = fmt.Sprintf("projects/%s/locations/us-east1", k.ProjectID)

	k.KeyRingName, err = k.CreateKeyRing(k.LocationName)
	if err != nil {
		return nil, fmt.Errorf("failed to create key ring: %w", err)
	}

	k.AsymmetricDecryptKeyName, err = k.CreateAsymmetricDecryptKey(k.KeyRingName)
	if err != nil {
		return nil, fmt.Errorf("failed to create asymmetric decrypt key: %w", err)
	}

	k.AsymmetricSignECKeyName, err = k.CreateAsymmetricSignECKey(k.KeyRingName)
	if err != nil {
		return nil, fmt.Errorf("failed to create asymmetric sign ec key: %w", err)
	}

	k.AsymmetricSignRSAKeyName, err = k.CreateAsymmetricSignRSAKey(k.KeyRingName)
	if err != nil {
		return nil, fmt.Errorf("failed to create asymmetric sign rsa key: %w", err)
	}

	k.HSMKeyName, err = k.CreateHSMKey(k.KeyRingName)
	if err != nil {
		return nil, fmt.Errorf("failed to create hsm key: %w", err)
	}

	k.SymmetricKeyName, err = k.CreateSymmetricKey(k.KeyRingName)
	if err != nil {
		return nil, fmt.Errorf("failed to create symmetric key: %w", err)
	}

	k.HMACKeyName, err = k.CreateHMACKey(k.KeyRingName)
	if err != nil {
		return nil, fmt.Errorf("failed to create hmac key: %w", err)
	}

	return &k, nil
}

// Cleanup deletes any resources
func (k *kmsFixture) Cleanup() error {
	ctx := context.Background()

	// Iterate over all keys and clean them up
	keyIt := k.client.ListCryptoKeys(ctx, &kmspb.ListCryptoKeysRequest{
		Parent: k.KeyRingName,
	})
	for {
		key, err := keyIt.Next()
		if err == iterator.Done {
			break
		}
		if grpcstatus.Code(err) == grpccodes.NotFound {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list keys for %s: %w", k.KeyRingName, err)
		}

		// Remove any rotation schedules
		if key.RotationSchedule != nil || key.NextRotationTime != nil {
			if _, err := k.client.UpdateCryptoKey(ctx, &kmspb.UpdateCryptoKeyRequest{
				CryptoKey: &kmspb.CryptoKey{
					Name:             key.Name,
					RotationSchedule: nil,
					NextRotationTime: nil,
				},
				UpdateMask: &fieldmask.FieldMask{
					Paths: []string{"rotation_period", "next_rotation_time"},
				},
			}); err != nil {
				return fmt.Errorf("failed to remove rotation schedule for %s: %w", key.Name, err)
			}
		}

		// Destroy all key versions
		versionIt := k.client.ListCryptoKeyVersions(ctx, &kmspb.ListCryptoKeyVersionsRequest{
			Parent: key.Name,
			Filter: "state != DESTROYED AND state != DESTROY_SCHEDULED",
		})

		for {
			version, err := versionIt.Next()
			if err == iterator.Done {
				break
			}
			if grpcstatus.Code(err) == grpccodes.NotFound {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to list versions for %s: %w", key.Name, err)
			}

			if _, err := k.client.DestroyCryptoKeyVersion(ctx, &kmspb.DestroyCryptoKeyVersionRequest{
				Name: version.Name,
			}); err != nil {
				return fmt.Errorf("failed to destroy version %s: %w", version.Name, err)
			}
		}
	}

	return nil
}

// RandomID returns a random UUID, useful for testing when values need to be
// unique.
func (k *kmsFixture) RandomID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Sprintf("failed to generate uuid: %v", err))
	}
	return u.String()
}

// CreateKeyRing creates a new key ring and returns the full resource name.
func (k *kmsFixture) CreateKeyRing(parent string) (string, error) {
	ctx := context.Background()
	keyRing, err := k.client.CreateKeyRing(ctx, &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: k.RandomID(),
	})
	if err != nil {
		return "", err
	}

	return keyRing.Name, nil
}

// CreateAsymmetricDecryptKey creates a new asymmetric key.
func (k *kmsFixture) CreateAsymmetricDecryptKey(parent string) (string, error) {
	ctx := context.Background()
	key, err := k.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: k.RandomID(),
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ASYMMETRIC_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_2048_SHA256,
			},
			Labels: map[string]string{
				"foo": "bar",
				"zip": "zap",
			},
		},
	})
	if err != nil {
		return "", err
	}

	return key.Name, nil
}

// CreateAsymmetricSignRSAKey creates a new asymmetric RSA key.
func (k *kmsFixture) CreateAsymmetricSignRSAKey(parent string) (string, error) {
	ctx := context.Background()
	key, err := k.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: k.RandomID(),
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ASYMMETRIC_SIGN,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
			},
			Labels: map[string]string{
				"foo": "bar",
				"zip": "zap",
			},
		},
	})
	if err != nil {
		return "", err
	}

	return key.Name, nil
}

// CreateAsymmetricSignECKey creates a new asymmetric EC key.
func (k *kmsFixture) CreateAsymmetricSignECKey(parent string) (string, error) {
	ctx := context.Background()
	key, err := k.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: k.RandomID(),
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ASYMMETRIC_SIGN,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256,
			},
			Labels: map[string]string{
				"foo": "bar",
				"zip": "zap",
			},
		},
	})
	if err != nil {
		return "", err
	}

	return key.Name, nil
}

// CreateHSMKey creates a new key backed by an HSM.
func (k *kmsFixture) CreateHSMKey(parent string) (string, error) {
	ctx := context.Background()
	key, err := k.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: k.RandomID(),
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm:       kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
				ProtectionLevel: kmspb.ProtectionLevel_HSM,
			},
			Labels: map[string]string{
				"foo": "bar",
				"zip": "zap",
			},
		},
	})
	if err != nil {
		return "", err
	}

	return key.Name, nil
}

// CreateSymmetricKey creates a new symmetric key.
func (k *kmsFixture) CreateSymmetricKey(parent string) (string, error) {
	ctx := context.Background()
	key, err := k.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: k.RandomID(),
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
			Labels: map[string]string{
				"foo": "bar",
				"zip": "zap",
			},
		},
	})
	if err != nil {
		return "", err
	}

	return key.Name, nil
}

// CreateHMACKey creates a new symmetric key.
func (k *kmsFixture) CreateHMACKey(parent string) (string, error) {
	ctx := context.Background()
	key, err := k.client.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: k.RandomID(),
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_MAC,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_HMAC_SHA256,
			},
			Labels: map[string]string{
				"foo": "bar",
				"zip": "zap",
			},
		},
	})
	if err != nil {
		return "", err
	}

	return key.Name, nil
}

// WaitForKeyVersionReady waits for the given key version to no longer be in the
// pending_generation state.
func (k *kmsFixture) WaitForKeyVersionReady(name string) error {
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		result, err := k.client.GetCryptoKeyVersion(ctx, &kmspb.GetCryptoKeyVersionRequest{
			Name: name,
		})
		if err != nil {
			return fmt.Errorf("waiting for %s ready: %w", name, err)
		}

		if result.State != kmspb.CryptoKeyVersion_PENDING_GENERATION {
			return nil
		}

		time.Sleep(time.Duration(i*250) * time.Millisecond)
	}

	return fmt.Errorf("key %s still not ready", name)
}
