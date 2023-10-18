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
	"os"
	"strings"
	"time"

	"testing"

	"cloud.google.com/go/bigquery"
	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

const (
	filePathToGCSForDeidTest       = "./testdata/dlp_sample.csv"
	tableID                        = "dlp_test_deid_table"
	dataSetID                      = "dlp_test_deid_dataset"
	deidentifyTemplateID           = "deidentified-templat-test-go"
	deidentifyStructuredTemplateID = "deidentified-structured-template-go"
	redactImageTemplate            = "redact-image-template-go"
)

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

func TestDeIdentifyTableWithCryptoHash(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	transientKeyName := "YOUR_TRANSIENT_CRYPTO_KEY_NAME"

	if err := deIdentifyTableWithCryptoHash(&buf, tc.ProjectID, transientKeyName); err != nil {
		t.Fatal(err)
	}
	got := buf.String()

	if want := "Table after de-identification :"; !strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithCryptoHash got %q, want %q", got, want)
	}
	if want := "user3@example.org"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithCryptoHash got %q, want %q", got, want)
	}
	if want := "858-555-0224"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithCryptoHash got %q, want %q", got, want)
	}
	if want := "user2@example.org"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithCryptoHash got %q, want %q", got, want)
	}
	if want := "858-555-0223"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithCryptoHash got %q, want %q", got, want)
	}
	if want := "user1@example.org"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithCryptoHash got %q, want %q", got, want)
	}
	if want := "858-555-0222"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithCryptoHash got %q, want %q", got, want)
	}
}

func TestDeIdentifyTableWithMultipleCryptoHash(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := deIdentifyTableWithMultipleCryptoHash(&buf, tc.ProjectID, "your-transient-crypto-key-name-1", "your-transient-crypto-key-name-2"); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Table after de-identification :"; !strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithMultipleCryptoHash got %q, want %q", got, want)
	}
	if want := "user1@example.org"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithMultipleCryptoHash got %q, want %q", got, want)
	}
	if want := "858-555-0222"; strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithMultipleCryptoHash got %q, want %q", got, want)
	}
	if want := "abbyabernathy1"; !strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTableWithMultipleCryptoHash got %q, want %q", got, want)
	}
}

func TestDeIdentifyTablePrimitiveBucketing(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	if err := deIdentifyTablePrimitiveBucketing(&buf, tc.ProjectID); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Table after de-identification :"; !strings.Contains(got, want) {
		t.Errorf("TestDeIdentifyTablePrimitiveBucketing got %q, want %q", got, want)
	}

	if want := `values:{string_value:"High"}`; !strings.Contains(got, want) {
		t.Errorf("TestDeidentifyDataReplaceWithDictionary got %q, want %q", got, want)
	}
	if want := `values:{string_value:"75"}`; strings.Contains(got, want) {
		t.Errorf("TestDeidentifyDataReplaceWithDictionary got %q, want %q", got, want)
	}
}

func TestDeidentifyDataReplaceWithDictionary(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := deidentifyDataReplaceWithDictionary(&buf, tc.ProjectID, "My name is Alicia Abernathy, and my email address is aabernathy@example.com."); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	want1 := "output: My name is Alicia Abernathy, and my email address is izumi@example.com."
	want2 := "output: My name is Alicia Abernathy, and my email address is alex@example.com."
	if !strings.Contains(got, want1) && !strings.Contains(got, want2) {
		t.Errorf("TestDeidentifyDataReplaceWithDictionary got %q, output does not contains value from dictionary", got)
	}
}

func TestDeidentifyCloudStorage(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	gcsURI := fmt.Sprint("gs://" + bucketForDeidCloudStorageForInput + "/" + filePathToGCSForDeidTest)
	outputBucket := fmt.Sprint("gs://" + bucketForDeidCloudStorageForOutput)

	fullDeidentifyTemplateID := fmt.Sprint("projects/" + tc.ProjectID + "/deidentifyTemplates/" + deidentifyTemplateID)
	fullDeidentifyStructuredTemplateID := fmt.Sprint("projects/" + tc.ProjectID + "/deidentifyTemplates/" + deidentifyStructuredTemplateID)
	fullRedactImageTemplate := fmt.Sprint("projects/" + tc.ProjectID + "/deidentifyTemplates/" + redactImageTemplate)

	if err := deidentifyCloudStorage(&buf, tc.ProjectID, gcsURI, tableID, dataSetID, outputBucket, fullDeidentifyTemplateID, fullDeidentifyStructuredTemplateID, fullRedactImageTemplate); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Job created successfully:"; !strings.Contains(got, want) {
		t.Errorf("TestDeidentifyCloudStorage got %q, want %q", got, want)
	}
}

var (
	u                                  = uuid.New().String()[:8]
	bucketForDeidCloudStorageForInput  = "dlp-test-deid-input-" + u
	bucketForDeidCloudStorageForOutput = "dlp-test-deid-output-" + u
)

func TestMain(m *testing.M) {
	tc := testutil.Context{}
	tc.ProjectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	createRedactImageTemplate(tc.ProjectID, redactImageTemplate)
	createDeidentifiedTemplate(tc.ProjectID, deidentifyTemplateID)
	createStructuredDeidentifiedTemplate(tc.ProjectID, deidentifyStructuredTemplateID)
	v := []string{bucketForDeidCloudStorageForInput, bucketForDeidCloudStorageForOutput}
	for _, v := range v {
		createBucket(tc.ProjectID, v)
	}
	filePathtoGCS(tc.ProjectID)
	createBigQueryDataSetId(tc.ProjectID)
	createTableInsideDataset(tc.ProjectID, dataSetID)
	m.Run()
	deleteBigQueryAssets(tc.ProjectID)
	for _, v := range v {
		deleteBucket(tc.ProjectID, v)
	}
	deleteTemplate(tc.ProjectID)
}

func createDeidentifiedTemplate(projectID, deidentifyTemplateID string) error {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	replaceWithInfoTypeConfig := &dlppb.ReplaceWithInfoTypeConfig{}

	infoTypeTransformations := &dlppb.InfoTypeTransformations{
		Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
			{PrimitiveTransformation: &dlppb.PrimitiveTransformation{
				Transformation: &dlppb.PrimitiveTransformation_ReplaceWithInfoTypeConfig{
					ReplaceWithInfoTypeConfig: replaceWithInfoTypeConfig,
				},
			}},
		},
	}
	deidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: infoTypeTransformations,
		},
	}
	template := &dlppb.DeidentifyTemplate{
		DeidentifyConfig: deidentifyConfig,
	}
	req := &dlppb.CreateDeidentifyTemplateRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyTemplate: template,
		TemplateId:         deidentifyTemplateID,
	}
	resp, err := client.CreateDeidentifyTemplate(ctx, req)
	if err != nil {
		return err
	}
	fmt.Print("\n" + "template " + resp.Name + "is created")
	return nil
}

func createStructuredDeidentifiedTemplate(projectID, deidentifyStructuredTemplateID string) error {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	value := &dlppb.Value_StringValue{
		StringValue: "Hello",
	}
	replaceValueConfig := &dlppb.ReplaceValueConfig{
		NewValue: &dlppb.Value{
			Type: value,
		},
	}
	recordTransformations := &dlppb.RecordTransformations{
		FieldTransformations: []*dlppb.FieldTransformation{
			{
				Transformation: &dlppb.FieldTransformation_PrimitiveTransformation{
					PrimitiveTransformation: &dlppb.PrimitiveTransformation{
						Transformation: &dlppb.PrimitiveTransformation_ReplaceConfig{
							ReplaceConfig: replaceValueConfig,
						},
					},
				},
			},
		},
	}
	deidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_RecordTransformations{
			RecordTransformations: recordTransformations,
		},
	}
	template := &dlppb.DeidentifyTemplate{
		DeidentifyConfig: deidentifyConfig,
	}
	req := &dlppb.CreateDeidentifyTemplateRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyTemplate: template,
		TemplateId:         deidentifyStructuredTemplateID,
	}
	resp, err := client.CreateDeidentifyTemplate(ctx, req)
	if err != nil {
		return err
	}
	fmt.Print("\n" + "template " + resp.Name + "is created")
	return nil
}

func createRedactImageTemplate(projectID, redactImageTemplate string) error {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	imageTransformation := &dlppb.ImageTransformations_ImageTransformation{
		RedactionColor: &dlppb.Color{
			Red:   1,
			Green: 0,
			Blue:  0,
		},
	}
	imageTransformations := &dlppb.ImageTransformations{
		Transforms: []*dlppb.ImageTransformations_ImageTransformation{
			imageTransformation,
		},
	}
	deidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_ImageTransformations{
			ImageTransformations: imageTransformations,
		},
	}
	template := &dlppb.DeidentifyTemplate{
		DeidentifyConfig: deidentifyConfig,
	}
	req := &dlppb.CreateDeidentifyTemplateRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyTemplate: template,
		TemplateId:         redactImageTemplate,
	}
	resp, err := client.CreateDeidentifyTemplate(ctx, req)
	if err != nil {
		return err
	}
	fmt.Print("\n" + "template " + resp.Name + "is created")
	return nil
}

func deleteTemplate(projectID string) error {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	abc := []string{deidentifyTemplateID, deidentifyStructuredTemplateID, redactImageTemplate}
	for _, v := range abc {
		name := fmt.Sprint("projects/" + projectID + "/deidentifyTemplates/" + v)
		req := &dlppb.DeleteDeidentifyTemplateRequest{
			Name: name,
		}
		err := client.DeleteDeidentifyTemplate(ctx, req)
		if err != nil {
			return err
		}
		log.Printf("[info] deleted a template : %s", v)
	}
	return nil
}

func createBucket(projectID, bucketName string) error {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Check if the bucket already exists.
	bucketExists := false
	_, err = client.Bucket(bucketName).Attrs(ctx)
	if err == nil {
		bucketExists = true
	}

	// If the bucket doesn't exist, create it.
	if !bucketExists {
		if err := client.Bucket(bucketName).Create(ctx, projectID, &storage.BucketAttrs{
			StorageClass: "STANDARD",
			Location:     "us-central1",
		}); err != nil {
			log.Fatalf("---Failed to create bucket: %v", err)
			return err
		}
		fmt.Printf("---Bucket '%s' created successfully.\n", bucketName)
	} else {
		fmt.Printf("---Bucket '%s' already exists.\n", bucketName)
	}
	fmt.Println("createbucket function is executed-------")
	return nil
}

func filePathtoGCS(projectID string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Check if the bucket already exists.
	bucketExists := false
	_, err = client.Bucket(bucketForDeidCloudStorageForInput).Attrs(ctx)
	if err == nil {
		bucketExists = true
	}

	// If the bucket doesn't exist, create it.
	if !bucketExists {
		if err := client.Bucket(bucketForDeidCloudStorageForInput).Create(ctx, projectID, &storage.BucketAttrs{
			StorageClass: "STANDARD",
			Location:     "us-central1",
		}); err != nil {
			return err
		}
		fmt.Printf("Bucket '%s' created successfully.\n", bucketForDeidCloudStorageForInput)
	} else {
		fmt.Printf("Bucket '%s' already exists.\n", bucketForDeidCloudStorageForInput)
	}

	// Check if the directory already exists in the bucket.
	dirExists := false
	query := &storage.Query{Prefix: filePathToGCSForDeidTest}
	it := client.Bucket(bucketForDeidCloudStorageForInput).Objects(ctx, query)
	_, err = it.Next()
	if err == nil {
		dirExists = true
	}

	// If the directory doesn't exist, create it.
	if !dirExists {
		obj := client.Bucket(bucketForDeidCloudStorageForInput).Object(filePathToGCSForDeidTest)
		if _, err := obj.NewWriter(ctx).Write([]byte("")); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		fmt.Printf("Directory '%s' created successfully in bucket '%s'.\n", filePathToGCSForDeidTest, bucketForDeidCloudStorageForInput)
	} else {
		fmt.Printf("Directory '%s' already exists in bucket '%s'.\n", filePathToGCSForDeidTest, bucketForDeidCloudStorageForInput)
	}

	// file upload code

	// Open local file.
	file, err := os.ReadFile(filePathToGCSForDeidTest)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return err
	}

	// Get a reference to the bucket
	bucket := client.Bucket(bucketForDeidCloudStorageForInput)

	// Upload the file
	object := bucket.Object(filePathToGCSForDeidTest)
	writer := object.NewWriter(ctx)
	_, err = writer.Write(file)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close writer: %v", err)
		return err
	}
	fmt.Printf("File uploaded successfully: %v\n", filePathToGCSForDeidTest)

	// Check if the file exists in the bucket
	_, err = bucket.Object(filePathToGCSForDeidTest).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			fmt.Printf("File %v does not exist in bucket %v\n", filePathToGCSForDeidTest, bucketForDeidCloudStorageForInput)
		} else {
			log.Fatalf("Failed to check file existence: %v", err)
		}
	} else {
		fmt.Printf("File %v exists in bucket %v\n", filePathToGCSForDeidTest, bucketForDeidCloudStorageForInput)
	}

	fmt.Println("filePathtoGCS function is executed-------")
	return nil
}

func deleteBucket(projectID, bucketName string) error {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	// List all objects in the bucket.
	objs := bucket.Objects(ctx, nil)
	for {
		objAttrs, err := objs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list objects in bucket: %v", err)
		}

		// Delete each object in the bucket.
		if err := bucket.Object(objAttrs.Name).Delete(ctx); err != nil {
			log.Fatalf("Failed to delete object %s: %v", objAttrs.Name, err)
		}
		fmt.Printf("Deleted object: %s\n", objAttrs.Name)
	}
	if err := bucket.Delete(ctx); err != nil {
		log.Fatalf("Failed to delete bucket: %v", err)
	}

	return nil
}

func createBigQueryDataSetId(projectID string) error {

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}

	if err := client.Dataset(dataSetID).Create(ctx, meta); err != nil {
		return err
	}

	return nil
}

func createTableInsideDataset(projectID, dataSetID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "user_id", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
		{Name: "title", Type: bigquery.StringFieldType},
		{Name: "score", Type: bigquery.StringFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema:         sampleSchema,
		ExpirationTime: time.Now().AddDate(1, 0, 0), // Table will be automatically deleted in 1 year.
	}

	tableRef := client.Dataset(dataSetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		log.Printf("[INFO] createBigQueryDataSetId Error while table creation: %v", err)
		return err
	}

	duration := time.Duration(90) * time.Second
	time.Sleep(duration)

	inserter := client.Dataset(dataSetID).Table(tableID).Inserter()
	items := []*BigQueryTableItem{
		// Item implements the ValueSaver interface.
		{UserId: "602-61-8588", Age: 32, Title: "Biostatistician III", Score: "A"},
		{UserId: "618-96-2322", Age: 69, Title: "Programmer I", Score: "C"},
		{UserId: "618-96-2322", Age: 69, Title: "Executive Secretary", Score: "C"},
	}
	if err := inserter.Put(ctx, items); err != nil {
		return err
	}

	return nil
}

type BigQueryTableItem struct {
	UserId string
	Age    int
	Title  string
	Score  string
}

func (i *BigQueryTableItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id": i.UserId,
		"age":     i.Age,
		"title":   i.Title,
		"score":   i.Score,
	}, bigquery.NoDedupeID, nil
}

func deleteBigQueryAssets(projectID string) error {

	log.Printf("[START] deleteBigQueryAssets: projectID %v and ", projectID)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	log.Printf("[INFO] deleteBigQueryAssets: delete dataset err %v", err)

	if err := client.Dataset(dataSetID).DeleteWithContents(ctx); err != nil {
		log.Printf("[INFO] deleteBigQueryAssets: delete dataset err %v", err)
		return err
	}

	duration := time.Duration(30) * time.Second
	time.Sleep(duration)

	log.Printf("[END] deleteBigQueryAssets:")
	return nil
}
