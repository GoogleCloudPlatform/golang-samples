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

//go:build ignore
// +build ignore

// Sample objects creates, list, deletes objects and runs
// other similar operations on them by using the Google Storage API.
// More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.

package objects

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

func main() {
	log.Fatalf("Running main is not supported.")
}

func list(w io.Writer, client *storage.Client, bucket string) error {
	// [START storage_list_files]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files]
	return nil
}

func listByPrefix(w io.Writer, client *storage.Client, bucket, prefix, delim string) error {
	// [START storage_list_files_with_prefix]
	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimiter, the entire  tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, attrs.Name)
	}
	// [END storage_list_files_with_prefix]
	return nil
}

func setEventBasedHold(client *storage.Client, bucket, object string) error {
	// [START storage_set_event_based_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		EventBasedHold: true,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_set_event_based_hold]
	return nil
}

func releaseEventBasedHold(client *storage.Client, bucket, object string) error {
	// [START storage_release_event_based_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		EventBasedHold: false,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_release_event_based_hold]
	return nil
}

func setTemporaryHold(client *storage.Client, bucket, object string) error {
	// [START storage_set_temporary_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		TemporaryHold: true,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_set_temporary_hold]
	return nil
}

func releaseTemporaryHold(client *storage.Client, bucket, object string) error {
	// [START storage_release_temporary_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := client.Bucket(bucket).Object(object)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		TemporaryHold: false,
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_release_temporary_hold]
	return nil
}

// writeEncryptedObject writes an object encrypted with user-provided AES key to a bucket.
func writeEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) error {
	// [START storage_upload_encrypted_file]
	ctx := context.Background()
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Encrypt the object's contents.
	wc := obj.Key(secretKey).NewWriter(ctx)
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_encrypted_file]
	return nil
}

// writeWithKMSKey writes an object encrypted with KMS-provided key to a bucket.
func writeWithKMSKey(client *storage.Client, bucket, object string, keyName string) error {
	// [START storage_upload_with_kms_key]
	ctx := context.Background()
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Encrypt the object's contents
	wc := obj.NewWriter(ctx)
	wc.KMSKeyName = keyName
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_with_kms_key]
	return nil
}

func readEncryptedObject(client *storage.Client, bucket, object string, secretKey []byte) ([]byte, error) {
	// [START storage_download_encrypted_file]
	ctx := context.Background()
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := obj.Key(secretKey).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	// [END storage_download_encrypted_file]
	return data, nil
}

func rotateEncryptionKey(client *storage.Client, bucket, object string, key, newKey []byte) error {
	// [START storage_rotate_encryption_key]
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	// obj is encrypted with key, we are encrypting it with the newKey.
	_, err = obj.Key(newKey).CopierFrom(obj.Key(key)).Run(ctx)
	if err != nil {
		return err
	}
	// [END storage_rotate_encryption_key]
	return nil
}

func downloadUsingRequesterPays(client *storage.Client, object, bucketName, localpath, billingProjectID string) error {
	// [START storage_download_file_requester_pays]
	ctx := context.Background()

	bucket := client.Bucket(bucketName).UserProject(billingProjectID)
	src := bucket.Object(object)

	f, err := os.OpenFile(localpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := src.NewReader(ctx)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, rc); err != nil {
		return err
	}
	if err := rc.Close(); err != nil {
		return err
	}
	fmt.Printf("Downloaded using %v as billing project.\n", billingProjectID)
	// [END storage_download_file_requester_pays]
	return nil
}
