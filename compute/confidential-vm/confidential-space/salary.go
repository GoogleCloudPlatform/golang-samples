// Copyright 2022 Google LLC
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

// [START confidential_space_salary]
// IMPORTANT: Before compiling, customize the details in the collaborators
// variable.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type collaborator struct {
	name         string
	wipname      string
	sa           string
	keyName      string
	inputBucket  string
	inputFile    string
	outputBucket string
	outputFile   string
}

// For simplicity we only have two people sharing data, but this can scale
// to more. You need to customize this section, replacing the <VARIABLE>
// values with your own.
var collaborators = [2]collaborator{
	{
		"Alex", // The name of the collaborator
		"projects/<ALEX_PROJECT_NUMBER>/locations/global/workloadIdentityPools/<ALEX_POOL_NAME>/providers/attestation-verifier", // Alex's workload identity pool. You can get Alex's project number with echo $(gcloud projects list --filter='<ALEX_PROJECT_ID>' --format='value(PROJECT_NUMBER)')
		"<ALEX_SERVICE_ACCOUNT_NAME>@<ALEX_PROJECT_ID>.iam.gserviceaccount.com",                                                 // Alex's service account that can decrypt their salary
		"projects/<ALEX_PROJECT_ID>/locations/global/keyRings/<ALEX_KEYRING_NAME>/cryptoKeys/<ALEX_KEY_NAME>",                   // Alex's KMS key
		"<ALEX_INPUT_BUCKET_NAME>",     // The storage bucket name that contains Alex's encrypted salary file
		"<ALEX_ENCRYPTED_SALARY_FILE>", // The name of Alex's encrypted salary file
		"<ALEX_RESULTS_BUCKET_NAME>",   // The name of Alex's results bucket
		"<ALEX_RESULTS_OUTPUT_FILE>",   // The name of Alex's output file that contains the results
	},
	{
		"Bola", // The name of the collaborator
		"projects/<BOLA_PROJECT_NUMBER>/locations/global/workloadIdentityPools/<BOLA_POOL_NAME>/providers/attestation-verifier", // Bola's workload identity pool. You can get Bola's project number with echo $(gcloud projects list --filter='<BOLA_PROJECT_ID>' --format='value(PROJECT_NUMBER)')
		"<BOLA_SERVICE_ACCOUNT_NAME>@<BOLA_PROJECT_ID>.iam.gserviceaccount.com",                                                 // Bola's service account that can decrypt their salary
		"projects/<BOLA_PROJECT_ID>/locations/global/keyRings/<BOLA_KEYRING_NAME>/cryptoKeys/<BOLA_KEY_NAME>",                   // Bola's KMS key
		"<BOLA_INPUT_BUCKET_NAME>",          // The storage bucket name that contains Bola's encrypted salary file
		"<BOLA_ENCRYPTED_BOLA_SALARY_FILE>", // The name of Bola's encrypted salary file
		"<BOLA_RESULTS_BUCKET_NAME>",        // The name of Bola's results bucket
		"<BOLA_RESULTS_OUTPUT_FILE>",        // The name of Bola's output file that contains the results
	},
}

const credentialConfig = `{
		"type": "external_account",
		"audience": "//iam.googleapis.com/%s",
		"subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
		"token_url": "https://sts.googleapis.com/v1/token",
		"credential_source": {
			"file": "/run/container_launcher/attestation_verifier_claims_token"
		},
		"service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/%s:generateAccessToken"
		}`

func main() {
	fmt.Println("workload started")
	ctx := context.Background()

	storageClient, err := storage.NewClient(ctx) // using the default credential on the Compute Engine VM
	if err != nil {
		fmt.Println(err)
		return
	}

	// get and decrypt
	s0, err := getSalary(ctx, storageClient, collaborators[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	s1, err := getSalary(ctx, storageClient, collaborators[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	res := ""
	if s0 > s1 {
		res = fmt.Sprintf("%s earns more!\n", collaborators[0].name)
	} else if s1 < s0 {
		res = fmt.Sprintf("%s earns more!\n", collaborators[1].name)
	} else {
		res = "earns same\n"
	}

	now := time.Now()
	for _, cw := range collaborators {
		outputWriter := storageClient.Bucket(cw.outputBucket).Object(fmt.Sprintf("%s-%d", cw.outputFile, now.Unix())).NewWriter(ctx)

		_, err = outputWriter.Write([]byte(res))
		if err != nil {
			fmt.Printf("Could not write: %v", err)
			return
		}
		if err = outputWriter.Close(); err != nil {
			fmt.Printf("Could not close: %v", err)
			return
		}
	}
}

func getSalary(ctx context.Context, storageClient *storage.Client, cw collaborator) (float64, error) {
	encryptedBytes, err := getFile(ctx, storageClient, cw.inputBucket, cw.inputFile)
	if err != nil {
		return 0.0, err
	}
	decryptedByte, err := decryptByte(ctx, cw.keyName, cw.sa, cw.wipname, encryptedBytes)
	if err != nil {
		return 0.0, err
	}
	decryptedNumber := strings.TrimSpace(string(decryptedByte))
	num, err := strconv.ParseFloat(decryptedNumber, 64)
	if err != nil {
		return 0.0, err
	}
	return num, nil
}

func decryptByte(ctx context.Context, keyName, trustedServiceAccountEmail, wippro string, encryptedData []byte) ([]byte, error) {
	cc := fmt.Sprintf(credentialConfig, wippro, trustedServiceAccountEmail)
	kmsClient, err := kms.NewKeyManagementClient(ctx, option.WithCredentialsJSON([]byte(cc)))
	if err != nil {
		return nil, fmt.Errorf("creating a new KMS client with federated credentials: %w", err)
	}

	decryptRequest := &kmspb.DecryptRequest{
		Name:       keyName,
		Ciphertext: encryptedData,
	}
	decryptResponse, err := kmsClient.Decrypt(ctx, decryptRequest)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt ciphertext: %w", err)
	}

	return decryptResponse.Plaintext, nil
}

func getFile(ctx context.Context, c *storage.Client, bucketName string, objPath string) ([]byte, error) {
	bucketHandle := c.Bucket(bucketName)
	objectHandle := bucketHandle.Object(objPath)

	objectReader, err := objectHandle.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer objectReader.Close()

	s, err := ioutil.ReadAll(objectReader)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// [END confidential_space_salary]
