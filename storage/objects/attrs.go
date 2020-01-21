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
package objects

// [START get_metadata]
import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

// attrs prints all of the object attributes.
func attrs(bucket, object string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Bucket: %v\n", attrs.Bucket)
	log.Printf("CacheControl: %v\n", attrs.CacheControl)
	log.Printf("ContentDisposition: %v\n", attrs.ContentDisposition)
	log.Printf("ContentEncoding: %v\n", attrs.ContentEncoding)
	log.Printf("ContentLanguage: %v\n", attrs.ContentLanguage)
	log.Printf("ContentType: %v\n", attrs.ContentType)
	log.Printf("Crc32c: %v\n", attrs.CRC32C)
	log.Printf("Generation: %v\n", attrs.Generation)
	log.Printf("KmsKeyName: %v\n", attrs.KMSKeyName)
	log.Printf("Md5Hash: %v\n", attrs.MD5)
	log.Printf("MediaLink: %v\n", attrs.MediaLink)
	log.Printf("Metageneration: %v\n", attrs.Metageneration)
	log.Printf("Name: %v\n", attrs.Name)
	log.Printf("Size: %v\n", attrs.Size)
	log.Printf("StorageClass: %v\n", attrs.StorageClass)
	log.Printf("TimeCreated: %v\n", attrs.Created)
	log.Printf("Updated: %v\n", attrs.Updated)
	log.Printf("Event-based hold enabled? %t\n", attrs.EventBasedHold)
	log.Printf("Temporary hold enabled? %t\n", attrs.TemporaryHold)
	log.Printf("Retention expiration time %v\n", attrs.RetentionExpirationTime)
	log.Print("\n\nMetadata\n")
	for key, value := range attrs.Metadata {
		log.Printf("\t%v = %v\n", key, value)
	}

	return attrs, nil
}

// [END get_metadata]
