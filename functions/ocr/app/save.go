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

package ocr

// [START functions_ocr_save]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

// saveResult is executed when a message is published to the Cloud Pub/Sub topic specified by
// RESULT_TOPIC in config.json file, and saves the data packet to a file in GCS.
func saveResult(w io.Writer, event pubsubpb.PubsubMessage) error {
	ctx := context.Background()
	var message ocrmessage
	if event.Data != nil {
		messageData := event.Data
		err := json.Unmarshal(messageData, &message)
		if err != nil {
			return fmt.Errorf("json.Unmarshal: %v", err)
		}
	} else {
		return fmt.Errorf("Empty data")
	}
	text := message.text
	filename := message.filename
	lang := message.lang

	fmt.Fprintf(w, "Received request to save file %q.", filename)

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	config := &config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	bucketName := config.resultBucket
	resultFilename := fmt.Sprintf("%s_%s.txt", filename, lang)
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	bucket := storageClient.Bucket(bucketName)

	fmt.Fprintf(w, "Saving result to %q in bucket %q.", resultFilename, bucketName)

	file := bucket.Object(resultFilename).NewWriter(ctx)
	defer file.Close()
	fmt.Fprint(file, text)

	print("File saved.")
	return nil
}

// [END functions_ocr_save]
