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

// This file might not be used for Go
// [START functions_ocr_process]
import (
	"context"
	"fmt"
	"io"
)

// processImage is executed when a file is uploaded to the Cloud Storage bucket you created for uploading images.
// It runs detectText, which processes the image for text.
func processImage(ctx context.Context, w io.Writer, projectID, bucket, name string) {
	// bucket, err := validateMessage(file, "bucket")
	// name, err := validateMessage(file, "name")

	detectText(w, projectID, bucket, name)

	fmt.Fprintf(w, "File %s processed.", name)
}

// [END functions_ocr_process]
