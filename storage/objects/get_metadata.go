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

package objects

// [START storage_get_metadata]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// getMetadata prints all of the object attributes.
func getMetadata(w io.Writer, bucket, object string) (*storage.ObjectAttrs, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).Attrs: %w", object, err)
	}
	fmt.Fprintf(w, "Bucket: %v\n", attrs.Bucket)
	fmt.Fprintf(w, "CacheControl: %v\n", attrs.CacheControl)
	fmt.Fprintf(w, "ContentDisposition: %v\n", attrs.ContentDisposition)
	fmt.Fprintf(w, "ContentEncoding: %v\n", attrs.ContentEncoding)
	fmt.Fprintf(w, "ContentLanguage: %v\n", attrs.ContentLanguage)
	fmt.Fprintf(w, "ContentType: %v\n", attrs.ContentType)
	fmt.Fprintf(w, "Crc32c: %v\n", attrs.CRC32C)
	fmt.Fprintf(w, "Generation: %v\n", attrs.Generation)
	fmt.Fprintf(w, "KmsKeyName: %v\n", attrs.KMSKeyName)
	fmt.Fprintf(w, "Md5Hash: %v\n", attrs.MD5)
	fmt.Fprintf(w, "MediaLink: %v\n", attrs.MediaLink)
	fmt.Fprintf(w, "Metageneration: %v\n", attrs.Metageneration)
	fmt.Fprintf(w, "Name: %v\n", attrs.Name)
	fmt.Fprintf(w, "Size: %v\n", attrs.Size)
	fmt.Fprintf(w, "StorageClass: %v\n", attrs.StorageClass)
	fmt.Fprintf(w, "TimeCreated: %v\n", attrs.Created)
	fmt.Fprintf(w, "Updated: %v\n", attrs.Updated)
	fmt.Fprintf(w, "Event-based hold enabled? %t\n", attrs.EventBasedHold)
	fmt.Fprintf(w, "Temporary hold enabled? %t\n", attrs.TemporaryHold)
	fmt.Fprintf(w, "Retention expiration time %v\n", attrs.RetentionExpirationTime)
	fmt.Fprintf(w, "Custom time %v\n", attrs.CustomTime)
	fmt.Fprintf(w, "Retention: %+v\n", attrs.Retention)
	fmt.Fprintf(w, "\n\nMetadata\n")
	for key, value := range attrs.Metadata {
		fmt.Fprintf(w, "\t%v = %v\n", key, value)
	}
	return attrs, nil
}

// [END storage_get_metadata]
