// Copyright 2019 Google LLC
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

// Package deid contains example snippets using the DLP deidentification API.
package deid

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"testing"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestMask(t *testing.T) {
	tc := testutil.SystemTest(t)
	tests := []struct {
		input            string
		maskingCharacter string
		numberToMask     int32
		want             string
	}{
		{
			input:            "My SSN is 111222333",
			maskingCharacter: "+",
			want:             "My SSN is +++++++++",
		},
		{
			input: "My SSN is 111222333",
			want:  "My SSN is *********",
		},
		{
			input:            "My SSN is 111222333",
			maskingCharacter: "+",
			numberToMask:     6,
			want:             "My SSN is ++++++333",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			test := test
			t.Parallel()
			buf := new(bytes.Buffer)
			err := mask(buf, tc.ProjectID, test.input, []string{"US_SOCIAL_SECURITY_NUMBER"}, test.maskingCharacter, test.numberToMask)
			if err != nil {
				t.Errorf("mask(%q, %s, %v) = error %q, want %q", test.input, test.maskingCharacter, test.numberToMask, err, test.want)
			}
			if got := buf.String(); got != test.want {
				t.Errorf("mask(%q, %s, %v) = %q, want %q", test.input, test.maskingCharacter, test.numberToMask, got, test.want)
			}
		})
	}
}

func TestDeidentifyDateShift(t *testing.T) {
	tc := testutil.SystemTest(t)
	tests := []struct {
		input      string
		want       string
		lowerBound int32
		upperBound int32
	}{
		{
			input:      "2016-01-10",
			lowerBound: 1,
			upperBound: 1,
			want:       "2016-01-11",
		},
		{
			input:      "2016-01-10",
			lowerBound: -1,
			upperBound: -1,
			want:       "2016-01-09",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			test := test
			t.Parallel()
			buf := new(bytes.Buffer)
			err := deidentifyDateShift(buf, tc.ProjectID, test.lowerBound, test.upperBound, test.input)
			if err != nil {
				t.Errorf("deidentifyDateShift(%v, %v, %q) = error '%q', want %q", test.lowerBound, test.upperBound, err, test.input, test.want)
			}
			if got := buf.String(); got != test.want {
				t.Errorf("deidentifyDateShift(%v, %v, %q) = %q, want %q", test.lowerBound, test.upperBound, got, test.input, test.want)
			}
		})
	}
}

func TestDeidentifyTableRowSuppress(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := deidentifyTableRowSuppress(&buf, tc.ProjectID); err != nil {
		t.Errorf("deidentifyTableRowSuppress: %v", err)
	}
	got := buf.String()
	if want := "Table after de-identification"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableRowSuppress got %q, want %q", got, want)
	}
	if want := "values:{string_value:\"Charles Dickens\"} "; strings.Contains(got, want) {
		t.Errorf("deidentifyTableRowSuppress got %q, want %q", got, want)
	}
}

func TestDeidentifyTableInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer

	if err := deidentifyTableInfotypes(&buf, tc.ProjectID); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Table after de-identification"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableInfotypes got %q, want %q", got, want)
	}

	if want := "[PERSON_NAME]"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableInfotypes got %q, want %q", got, want)
	}

	if want := "Charles Dickens"; strings.Contains(got, want) {
		t.Errorf("deidentifyTableInfotypes got %q, want %q", got, want)
	}
	if want := "Mark Twain"; strings.Contains(got, want) {
		t.Errorf("deidentifyTableInfotypes got %q, want %q", got, want)
	}
	if want := "Jane Austen"; strings.Contains(got, want) {
		t.Errorf("deidentifyTableInfotypes got %q, want %q", got, want)
	}

}

func TestDeIdentifyWithRedact(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "My name is Alicia Abernathy, and my email address is aabernathy@example.com."
	infoTypeNames := []string{"EMAIL_ADDRESS"}
	want := "output: My name is Alicia Abernathy, and my email address is ."

	var buf bytes.Buffer

	if err := deidentifyWithRedact(&buf, tc.ProjectID, input, infoTypeNames); err != nil {
		t.Errorf("deidentifyWithRedact(%q) = error '%q', want %q", err, input, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyWithRedact(%q) = %q, want %q", got, input, want)
	}
}

func TestDeidentifyExceptionList(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "jack@example.org accessed customer record of user5@example.com"
	want := "output : jack@example.org accessed customer record of [EMAIL_ADDRESS]"

	var buf bytes.Buffer

	if err := deidentifyExceptionList(&buf, tc.ProjectID, input); err != nil {
		t.Errorf("deidentifyExceptionList(%q) = error '%q', want %q", input, err, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyExceptionList(%q) = %q, want %q", input, got, want)
	}

}

func TestDeIdentifyWithReplacement(t *testing.T) {
	tc := testutil.SystemTest(t)
	input := "My name is Alicia Abernathy, and my email address is aabernathy@example.com."
	infoType := []string{"EMAIL_ADDRESS"}
	replaceVal := "[email-address]"
	want := "output : My name is Alicia Abernathy, and my email address is [email-address]."

	var buf bytes.Buffer
	err := deidentifyWithReplacement(&buf, tc.ProjectID, input, infoType, replaceVal)
	if err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyWithReplacement(%q) = %q, want %q", input, got, want)
	}
}

func TestDeidentifyTableBucketing(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := deIdentifyTableBucketing(&buf, tc.ProjectID); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "values:{string_value:\"70:80\"}}"; !strings.Contains(got, want) {
		t.Errorf("deIdentifyTableBucketing got %q, want %q", got, want)
	}
	if want := "values:{string_value:\"75\"}}"; strings.Contains(got, want) {
		t.Errorf("deIdentifyTableBucketing got %q, want %q", got, want)
	}

}

func TestDeidentifyTableMaskingCondition(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	if err := deidentifyTableMaskingCondition(&buf, tc.ProjectID); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Table after de-identification :"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableMaskingCondition got (%q) =%q ", got, want)
	}
	if want := "values:{string_value:\"**\"}"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableMaskingCondition got (%q) =%q ", got, want)
	}
}

func TestDeidentifyTableConditionInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := deidentifyTableConditionInfoTypes(&buf, tc.ProjectID, []string{"PATIENT", "FACTOID"}); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Table after de-identification"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableConditionInfoTypes got %q, want %q", got, want)
	}
	if want := "values:{string_value:\"[PERSON_NAME] name was a curse invented by [PERSON_NAME].\"}}"; !strings.Contains(got, want) {
		t.Errorf("deidentifyTableConditionInfoTypes got %q, want %q", got, want)
	}
}

func TestDeIdentifyWithWordList(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	input := "Patient was seen in RM-YELLOW then transferred to rm green."
	infoType := "CUSTOM_ROOM_ID"
	wordList := []string{"RM-GREEN", "RM-YELLOW", "RM-ORANGE"}
	want := "output : Patient was seen in [CUSTOM_ROOM_ID] then transferred to [CUSTOM_ROOM_ID]."

	if err := deidentifyWithWordList(&buf, tc.ProjectID, input, infoType, wordList); err != nil {
		t.Errorf("deidentifyWithWordList(%q) = error '%q', want %q", input, err, want)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyWithWordList(%q) = %q, want %q", input, got, want)
	}
}

func TestDeIdentifyWithInfotype(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "My email is test@example.com"
	infoType := []string{"EMAIL_ADDRESS"}
	want := "output : My email is [EMAIL_ADDRESS]"

	var buf bytes.Buffer

	if err := deidentifyWithInfotype(&buf, tc.ProjectID, input, infoType); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); got != want {
		t.Errorf("deidentifyFreeTextWithFPEUsingSurrogate(%q) = %q, want %q", input, got, want)
	}

}

func TestDeidentifyTableFPE(t *testing.T) {
	tc := testutil.SystemTest(t)

	keyRingName, err := createKeyRing(t, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	kmsKeyName, wrappedAesKey, keyVersion, err := createKey(t, tc.ProjectID, keyRingName)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyKey(t, tc.ProjectID, keyVersion)

	contains := "De-identify Table after format-preserving encryption"

	var buf bytes.Buffer

	if err := deidentifyTableFPE(&buf, tc.ProjectID, kmsKeyName, wrappedAesKey); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); !strings.Contains(got, contains) {
		t.Errorf("deidentifyTableFPE() = %q,%q ", got, contains)
	}
}
func TestDeIdentifyDeterministic(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "Jack's phone number is 5555551212"
	infoTypeNames := []string{"PHONE_NUMBER"}
	keyRingName, err := createKeyRing(t, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	keyFileName, cryptoKeyName, keyVersion, err := createKey(t, tc.ProjectID, keyRingName)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyKey(t, tc.ProjectID, keyVersion)

	surrogateInfoType := "PHONE_TOKEN"
	want := "output : Jack's phone number is PHONE_TOKEN(36):"

	var buf bytes.Buffer

	if err := deIdentifyDeterministicEncryption(&buf, tc.ProjectID, input, infoTypeNames, keyFileName, cryptoKeyName, surrogateInfoType); err != nil {
		t.Errorf("deIdentifyDeterministicEncryption(%q) = error '%q', want %q", err, input, want)
	}

	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("deIdentifyDeterministicEncryption(%q) = %q, want %q", input, got, want)
	}

}

func createKeyRing(t *testing.T, projectID string) (string, error) {
	t.Helper()

	u := uuid.New().String()[:8]
	parent := fmt.Sprintf("projects/%v/locations/global", projectID)
	id := "test-dlp-go-lang-key-id-1" + u

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: id,
	}

	// Call the API.
	result, err := client.CreateKeyRing(ctx, req)
	if err != nil {
		return "", err
	}

	return result.Name, nil
}

func createKey(t *testing.T, projectID, keyFileName string) (string, string, string, error) {
	t.Helper()
	u := uuid.New().String()[:8]
	id := "go-lang-dlp-test-wrapped-aes-256" + u
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      keyFileName,
		CryptoKeyId: id,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				ProtectionLevel: kmspb.ProtectionLevel_HSM,
				Algorithm:       kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
		},
	}

	// Call the API.
	result, err := client.CreateCryptoKey(ctx, req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create key: %w", err)
	}

	response, err := client.Encrypt(ctx, &kmspb.EncryptRequest{
		Name:      result.Name,
		Plaintext: []byte("5u8x/A?D(G+KbPeShVmYq3t6w9y$B&E)"),
	})

	if err != nil {
		log.Fatalf("Failed to wrap key: %v", err)
	}

	wrappedKey := response.Ciphertext

	wrappedKeyString := base64.StdEncoding.EncodeToString(wrappedKey)
	return result.Name, wrappedKeyString, response.Name, nil
}

func destroyKey(t *testing.T, projectID, key string) error {
	t.Helper()

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &kmspb.DestroyCryptoKeyVersionRequest{
		Name: key,
	}

	_, err = client.DestroyCryptoKeyVersion(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func TestDeIdentifyFreeTextWithFPEUsingSurrogate(t *testing.T) {
	tc := testutil.SystemTest(t)

	input := "My phone number is 5555551212"
	infoType := "PHONE_NUMBER"
	surrogateType := "PHONE_TOKEN"
	unWrappedKey, err := getUnwrappedKey(t)
	if err != nil {
		t.Fatal(err)
	}
	want := "output: My phone number is PHONE_TOKEN(10):"

	var buf bytes.Buffer
	if err := deidentifyFreeTextWithFPEUsingSurrogate(&buf, tc.ProjectID, input, infoType, surrogateType, unWrappedKey); err != nil {
		t.Fatal(err)
	}
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("deidentifyFreeTextWithFPEUsingSurrogate(%q) = %q, want %q", input, got, want)
	}
}

func getUnwrappedKey(t *testing.T) (string, error) {
	t.Helper()
	key := make([]byte, 32) // 32 bytes for AES-256
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	// Encode the key to base64
	encodedKey := base64.StdEncoding.EncodeToString(key)
	return string(encodedKey), nil

}
