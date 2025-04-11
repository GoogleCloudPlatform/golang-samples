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

// [START functions_ocr_save]

package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// SaveResult is executed when a message is published to the Cloud Pub/Sub topic
// specified by the RESULT_TOPIC environment vairable, and saves the data packet
// to a file in GCS.
func SaveResult(ctx context.Context, event PubSubMessage) error {
	if err := setup(ctx); err != nil {
		return fmt.Errorf("ProcessImage: %w", err)
	}
	var message ocrMessage
	if event.Data == nil {
		return fmt.Errorf("Empty data")
	}
	if err := json.Unmarshal(event.Data, &message); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	log.Printf("Received request to save file %q.", message.FileName)

	resultFilename := fmt.Sprintf("%s_%s.txt", message.FileName, message.Lang)
	bucket := storageClient.Bucket(resultBucket)

	log.Printf("Saving result to %q in bucket %q.", resultFilename, resultBucket)

	w := bucket.Object(resultFilename).NewWriter(ctx)
	defer w.Close()
	fmt.Fprint(w, message.Text)

	log.Printf("File saved.")
	return nil
}

// [END functions_ocr_save]
