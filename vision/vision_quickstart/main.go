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

// [START vision_quickstart]

// Sample vision-quickstart uses the Google Cloud Vision API to label an image.
package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	//	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()

	// Creates a cloud storage client.
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create storageClient: %v", err)
	}
	defer storageClient.Close()

	// Creates a cloud vision client.
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create visionClient: %v", err)
	}
	defer visionClient.Close()

	// Declare a source bucket
	sourceBucketName := "david-chang-bucket-sandbox"
	sourceBucket := storageClient.Bucket(sourceBucketName)

	// Declare a destination bucket
	destinationBucketName := "david-chang-bucket-processed"

	// Generate a list of all objects in source bucket
	query := &storage.Query{Prefix: ""}

	// Declare a dictionary of ContentTypes we want to include
	contentTypesToInclude := map[string]bool {
		"image/CR2": true,
		"image/gif": true,
		"image/jpeg": true,
		"image/png": true,
		"image/x-icon": true,
		// "video/3gpp": true,
		// "video/avi": true,
		// "video/mp4": true,
		// "video/mpeg": true,
		// "video/quicktime": true,
		// "video/x-ms-wmv": true,
	}

	// Declare a SafeSearch values to include
	safeSearchValuesToInclude := map[string]bool {
		"UNKOWN": true,
		"POSSIBLE": true,
		"LIKELY": true,
		"VERY_LIKELY": true,
	}

	var names []string
	it := sourceBucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if contentTypesToInclude[attrs.ContentType] {
			names = append(names, attrs.Name)
		}
	}

	// Loop through all objects in source bucket and run through SafeSearch
	// Copy objects that get flagged in SafeSearch
	for _, name := range names {
		sourceObjURI := fmt.Sprintf("gs://%s/%s", sourceBucketName, name)
		fmt.Println(sourceObjURI)

		// Call vision API to detect SafeSearch attributes
		image := vision.NewImageFromURI(sourceObjURI)
		props, err := visionClient.DetectSafeSearch(ctx, image, nil)
		
		// if err := detectSafeSearchURI(os.Stdout, sourceObjURI); err != nil {
		// 	fmt.Println("detectSafeSearchURI: %v", err)
		// }

		if err != nil {
			fmt.Println("DetectSafeSearch error: %v", err)
		}

		// props values can be one of:
		// UNKOWN, VERY_UNLIKELY, UNLIKELY, POSSIBLE, LIKELY, VERY_LIKELY
		fmt.Fprintln(os.Stdout, "Safe Search properties:")
		fmt.Fprintln(os.Stdout, "Adult:", props.Adult)
		fmt.Fprintln(os.Stdout, "Medical:", props.Medical)
		fmt.Fprintln(os.Stdout, "Racy:", props.Racy)
		fmt.Fprintln(os.Stdout, "Spoofed:", props.Spoof)
		fmt.Fprintln(os.Stdout, "Violence:", props.Violence)
		adult := fmt.Sprintf("%s", props.Adult)
		racy := fmt.Sprintf("%s", props.Racy)
		if (safeSearchValuesToInclude[adult]) {
			if err := copyFile(ioutil.Discard, destinationBucketName, sourceBucketName, name, "adult", adult); err != nil {
				fmt.Println("copyFile error: %v", err)
			}
		}
		if (safeSearchValuesToInclude[racy]) {
			if err := copyFile(ioutil.Discard, destinationBucketName, sourceBucketName, name, "racy", racy); err != nil {
				fmt.Println("copyFile error: %v", err)
			}
		}
	}
}

// [END vision_quickstart]

// [START storage_copy_file]

// copyFile copies an object into specified bucket.
func copyFile(w io.Writer, dstBucket, srcBucket, srcObject string, categorization string, confidence string) error {
	// dstBucket := "bucket-1"
	// srcBucket := "bucket-2"
	// srcObject := "object"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dstObject := categorization + "/" + confidence + "/" + srcObject
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstObject, srcObject, err)
	}
	fmt.Fprintf(w, "Blob %v in bucket %v copied to blob %v in bucket %v.\n", srcObject, srcBucket, dstObject, dstBucket)
	return nil
}

// [END storage_copy_file]
