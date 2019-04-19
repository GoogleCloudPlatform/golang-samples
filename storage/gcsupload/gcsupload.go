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

// gcsupload is a CLI to upload a file to Google Cloud Storage.
// Invoke -h to see its flags:
// gcsupload -h
//  Usage of gcsupload:
//    -bucket string
//	    the bucket to upload content to
//    -name string
//	    the name of the file to be stored on GCS
//    -project string
//	    the ID of the GCP project to use
//    -public
//	    whether the item should be available publicly (default true)
//    -source string
//	    the path to the source
//
// For example:
//  gcsupload --project gcs-samples --source ~/Desktop/birthdayPic.jpg --bucket gcs-cli-test
//  URL: https://storage.googleapis.com/gcs-cli-test/birthdayPic.jpg
//  Size: 865096
//  MD5: 5b6c7b4aed837e8ed0f9950564a10b32
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

func main() {
	// Prevent log from printing out time information
	log.SetFlags(0)

	var projectID, bucket, source, name string
	var public bool

	flag.StringVar(&bucket, "bucket", "", "the bucket to upload content to")
	flag.StringVar(&projectID, "project", "", "the ID of the GCP project to use")
	flag.StringVar(&source, "source", "", "the path to the source")
	flag.StringVar(&name, "name", "", "the name of the file to be stored on GCS")
	flag.BoolVar(&public, "public", true, "whether the item should be available publicly")
	flag.Parse()

	// If they haven't set the bucket or projectID nor specified
	// in the environment, then fail if missing.
	bucket = mustGetEnv("GOLANG_SAMPLES_BUCKET", bucket)
	projectID = mustGetEnv("GOLANG_SAMPLES_PROJECT_ID", projectID)

	var r io.Reader
	if source == "" {
		r = os.Stdin
		log.Printf("Reading from stdin...")
	} else {
		f, err := os.Open(source)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r = f
	}

	if name == "" {
		if source != "" {
			name = filepath.Base(source)
		} else {
			name = "test-sample"
		}
	}

	ctx := context.Background()
	_, objAttrs, err := upload(ctx, r, projectID, bucket, name, public)
	if err != nil {
		switch err {
		case storage.ErrBucketNotExist:
			log.Fatal("Please create the bucket first e.g. with `gsutil mb`")
		default:
			log.Fatal(err)
		}
	}

	log.Printf("URL: %s", objectURL(objAttrs))
	log.Printf("Size: %d", objAttrs.Size)
	log.Printf("MD5: %x", objAttrs.MD5)
	log.Printf("objAttrs: %+v", objAttrs)
}

func objectURL(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}

func upload(ctx context.Context, r io.Reader, projectID, bucket, name string, public bool) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	bh := client.Bucket(bucket)
	// Next check if the bucket exists
	if _, err = bh.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bh.Object(name)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		return nil, nil, err
	}
	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if public {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return nil, nil, err
		}
	}

	attrs, err := obj.Attrs(ctx)
	return obj, attrs, err
}

func mustGetEnv(envKey, defaultValue string) string {
	val := os.Getenv(envKey)
	if val == "" {
		val = defaultValue
	}
	if val == "" {
		log.Fatalf("%q should be set", envKey)
	}
	return val
}
