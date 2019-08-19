// Copyright 2018, Google, LLC.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

// SaveResult is executed when a message is published to the Cloud Pub/Sub topic specified by
// RESULT_TOPIC in config.json file, and saves the data packet to a file in GCS.
func SaveResult(ctx context.Context, event PubSubMessage) error {
	var message ocrMessage
	if event.Data != nil {
		messageData := event.Data
		err := json.Unmarshal(messageData, &message)
		if err != nil {
			return fmt.Errorf("json.Unmarshal: %v", err)
		}
	} else {
		return fmt.Errorf("Empty data")
	}
	text := message.Text
	fileName := message.FileName
	lang := message.Lang

	log.Printf("Received request to save file %q.", fileName)

	bucketName := config.ResultBucket
	resultFilename := fmt.Sprintf("%s_%s.txt", fileName, lang)
	bucket := storageClient.Bucket(bucketName)

	log.Printf("Saving result to %q in bucket %q.", resultFilename, bucketName)

	file := bucket.Object(resultFilename).NewWriter(ctx)
	defer file.Close()
	fmt.Fprint(file, text)

	log.Printf("File saved.")
	return nil
}

// [END functions_ocr_save]
