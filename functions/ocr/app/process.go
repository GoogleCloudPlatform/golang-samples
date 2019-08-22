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

// [START functions_ocr_process]

package ocr

import (
	"context"
	"fmt"
	"log"
)

// ProcessImage is executed when a file is uploaded to the Cloud Storage bucket you
// created for uploading images. It runs detectText, which processes the image for text.
func ProcessImage(ctx context.Context, event GCSEvent) error {
	if err := setup(ctx); err != nil {
		return fmt.Errorf("ProcessImage: %v", err)
	}
	if event.Bucket == "" {
		return fmt.Errorf("empty file.Bucket")
	}
	if event.Name == "" {
		return fmt.Errorf("empty file.Name")
	}
	if err := detectText(ctx, event.Bucket, event.Name); err != nil {
		return fmt.Errorf("detectText: %v", err)
	}
	log.Printf("File %s processed.", event.Name)
	return nil
}

// [END functions_ocr_process]
