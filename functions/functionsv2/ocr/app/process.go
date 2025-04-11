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

// [START functions_ocr_process]

package ocr

import (
	"context"
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	functions.CloudEvent("process-image", ProcessImage)
}

// ProcessImage is executed when a file is uploaded to the Cloud Storage bucket you
// created for uploading images. It runs detectText, which processes the image for text.
func ProcessImage(ctx context.Context, cloudevent event.Event) error {
	if err := setup(ctx); err != nil {
		return fmt.Errorf("ProcessImage: %w", err)
	}

	var data storagedata.StorageObjectData

	// If you omit `DiscardUnknown`, then protojson.Unmarshal returns an error
	// when encountering a new or unknown field.
	options := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	err := options.Unmarshal(cloudevent.Data(), &data)
	if err != nil {
		return fmt.Errorf("protojson.Unmarshal: Failed to parse CloudEvent data: %w", err)
	}
	if data.GetBucket() == "" {
		return fmt.Errorf("empty file.Bucket")
	}
	if data.GetName() == "" {
		return fmt.Errorf("empty file.Name")
	}
	if err := detectText(ctx, data.GetBucket(), data.GetName()); err != nil {
		return fmt.Errorf("detectText: %w", err)
	}
	log.Printf("File %s processed.", data.GetName())
	return nil
}

// [END functions_ocr_process]
