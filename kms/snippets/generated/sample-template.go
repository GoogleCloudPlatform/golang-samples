// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

//+build ignore

// This file is used as the basis for generating snippet_test.go
// To re-generate, run:
//   go generate
// Boilerplate client code is inserted in the sections marked
//   "Boilerplate is inserted by gen.go"

package kms_snippets

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"

	cloudkms "google.golang.org/api/cloudkms/v1"
)

func init() {
	// For imports
	_ = context.Background
	_ = google.DefaultClient
}

func createKeyring(project, keyRing string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s", project, location)

	_, err = client.Projects.Locations.KeyRings.Create(
		parent, &cloudkms.KeyRing{}).KeyRingId(keyRing).Do()
	if err != nil {
		return err
	}
	log.Print("Created key ring.")

	return nil
}

func createCryptoKey(project, keyRing, key string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", project, location, keyRing)
	purpose := "ENCRYPT_DECRYPT"

	_, err = client.Projects.Locations.KeyRings.CryptoKeys.Create(
		parent, &cloudkms.CryptoKey{
			Purpose: purpose,
		}).CryptoKeyId(key).Do()
	if err != nil {
		return err
	}
	log.Print("Created crypto key.")

	return nil
}

func disableCryptoKeyVersion(project, keyRing, key, version string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeyVersions/%s",
		project, location, keyRing, version)

	_, err = client.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.Patch(
		parent, &cloudkms.CryptoKeyVersion{
			State: "DISABLED",
		}).UpdateMask("state").Do()
	if err != nil {
		return err
	}
	log.Print("Disabled crypto key version.")

	return nil
}

func enableCryptoKeyVersion(project, keyRing, key, version string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeyVersions/%s",
		project, location, keyRing, version)

	_, err = client.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.Patch(
		parent, &cloudkms.CryptoKeyVersion{
			State: "ENABLED",
		}).UpdateMask("state").Do()
	if err != nil {
		return err
	}
	log.Print("Enabled crypto key version.")

	return nil
}

func destroyCryptoKeyVersion(project, keyRing, key, version string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeyVersions/%s",
		project, location, keyRing, version)

	_, err = client.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.Destroy(
		parent, &cloudkms.DestroyCryptoKeyVersionRequest{}).Do()
	if err != nil {
		return err
	}
	log.Print("Destroyed crypto key version.")

	return nil
}

func restoreCryptoKeyVersion(project, keyRing, key, version string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeyVersions/%s",
		project, location, keyRing, version)

	_, err = client.Projects.Locations.KeyRings.CryptoKeys.CryptoKeyVersions.Restore(
		parent, &cloudkms.RestoreCryptoKeyVersionRequest{}).Do()
	if err != nil {
		return err
	}
	log.Print("Restored crypto key version.")

	return nil
}

func getRingPolicy(project, keyRing string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s",
		project, location, keyRing)

	policy, err := client.Projects.Locations.KeyRings.GetIamPolicy(parent).Do()
	if err != nil {
		return err
	}
	for _, binding := range policy.Bindings {
		log.Printf("Role: %s\n", binding.Role)
		log.Printf("Members: %v\n", binding.Members)
	}

	return nil
}

func addMemberRingPolicy(project, location, keyRing, role, member string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go

	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s",
		project, location, keyRing)

	policy, err := client.Projects.Locations.KeyRings.GetIamPolicy(parent).Do()
	if err != nil {
		return err
	}
	policy.Bindings = append(policy.Bindings, &cloudkms.Binding{
		Role:    role,
		Members: []string{member},
	})
	if err != nil {
		return err
	}

	_, err = client.Projects.Locations.KeyRings.SetIamPolicy(
		parent, &cloudkms.SetIamPolicyRequest{
			Policy: policy,
		}).Do()
	if err != nil {
		return err
	}

	return nil
}

func getCryptoKeyPolicy(project, keyRing, key string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeyVersions/%s",
		project, location, keyRing, key)

	policy, err := client.Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(parent).Do()
	if err != nil {
		return err
	}
	for _, binding := range policy.Bindings {
		log.Printf("Role: %s\n", binding.Role)
		log.Printf("Members: %v\n", binding.Members)
	}

	return nil
}

func addMemberCryptoKeyPolicy(project, keyRing, key, role, member string) error {
	var client *cloudkms.Service // Boilerplate is inserted by gen.go
	location := "global"
	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeyVersions/%s",
		project, location, keyRing, key)

	policy, err := client.Projects.Locations.KeyRings.CryptoKeys.GetIamPolicy(parent).Do()
	if err != nil {
		return err
	}
	policy.Bindings = append(policy.Bindings, &cloudkms.Binding{
		Role:    role,
		Members: []string{member},
	})
	if err != nil {
		return err
	}

	_, err = client.Projects.Locations.KeyRings.CryptoKeys.SetIamPolicy(
		parent, &cloudkms.SetIamPolicyRequest{
			Policy: policy,
		}).Do()
	if err != nil {
		return err
	}

	return nil
}
