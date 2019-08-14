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

// [START functions_ocr_process]
import (
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// processImage is executed when a file is uploaded to the Cloud Storage bucket you created for uploading images.
// It runs detectText, which processes the image for text.
func processImage(w io.Writer, projectID string, file storage.ObjectAttrs) error {
	// projectID := "my-project-id"
	if file.Bucket == "" {
		return fmt.Errorf("empty file.Bucket")
	}
	if file.Name == "" {
		return fmt.Errorf("empty file.Name")
	}
	detectText(w, projectID, file.Bucket, file.Name)
	fmt.Fprintf(w, "File %s processed.", file.Name)
	return nil
}

// [END functions_ocr_process]
