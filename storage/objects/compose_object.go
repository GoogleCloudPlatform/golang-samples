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

// [START storage_compose_file]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// composeFile composes source objects to create a composite object.
func composeFile(w io.Writer, bucket, object1, object2, toObject string) error {
	// bucket := "bucket-name"
	// object1 := "object-name-1"
	// object2 := "object-name-2"
	// toObject := "object-name-3"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	src1 := client.Bucket(bucket).Object(object1)
	src2 := client.Bucket(bucket).Object(object2)
	dst := client.Bucket(bucket).Object(toObject)

	// ComposerFrom takes varargs, so you can put as many objects here
	// as you want.
	_, err = dst.ComposerFrom(src1, src2).Run(ctx)
	if err != nil {
		return fmt.Errorf("ComposerFrom: %w", err)
	}
	fmt.Fprintf(w, "New composite object %v was created by combining %v and %v\n", toObject, object1, object2)
	return nil
}

// [END storage_compose_file]
