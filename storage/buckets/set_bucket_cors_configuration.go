// Copyright 2020 Google LLC
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

package buckets

// [START storage_cors_configuration]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// setBucketCORSConfiguration sets a CORS configuration on a bucket.
func setBucketCORSConfiguration(w io.Writer, bucketName string, maxAge time.Duration, methods, origins, responseHeaders []string) error {
	// bucketName := "bucket-name"
	// maxAge := time.Hour
	// methods := []string{"GET"}
	// origins := []string{"some-origin.com"}
	// responseHeaders := []string{"Content-Type"}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		CORS: []storage.CORS{
			{
				MaxAge:          maxAge,
				Methods:         methods,
				Origins:         origins,
				ResponseHeaders: responseHeaders,
			}},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Bucket %v was updated with a CORS config to allow %v requests from %v sharing %v responses across origins\n", bucketName, methods, origins, responseHeaders)
	return nil
}

// [END storage_cors_configuration]
