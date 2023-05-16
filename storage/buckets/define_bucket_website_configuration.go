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

// [START storage_define_bucket_website_configuration]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// setBucketWebsiteInfo sets website configuration on a bucket.
func setBucketWebsiteInfo(w io.Writer, bucketName, indexPage, notFoundPage string) error {
	// bucketName := "www.example.com"
	// indexPage := "index.html"
	// notFoundPage := "404.html"
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
		Website: &storage.BucketWebsite{
			MainPageSuffix: indexPage,
			NotFoundPage:   notFoundPage,
		},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Static website bucket %v is set up to use %v as the index page and %v as the 404 page\n", bucketName, indexPage, notFoundPage)
	return nil
}

// [END storage_define_bucket_website_configuration]
