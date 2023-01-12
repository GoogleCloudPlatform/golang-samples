// Copyright 2022 Google LLC
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

package clientendpoint

// [START storage_set_client_endpoint]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// setClientEndpoint sets the request endpoint.
func setClientEndpoint(w io.Writer, customEndpoint string, opts ...option.ClientOption) error {
	// customEndpoint := "https://my-custom-endpoint.example.com/storage/v1/"
	// opts := []option.ClientOption
	ctx := context.Background()

	// Set a custom request endpoint for this client.
	opts = append(opts, option.WithEndpoint(customEndpoint))
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Perfrom some operations with custom request endpoint set.
	performSomeOperations(client)
	return nil
}

// [END storage_set_client_endpoint]

func performSomeOperations(client *storage.Client) error {
	ctx := context.Background()
	bucket := "myBucket"
	object := "myObject"

	// Get bucket metadata.
	client.Bucket(bucket).Attrs(ctx)

	// Upload an object with storage.Writer.
	o := client.Bucket(bucket).Object(object)
	w := o.NewWriter(ctx)
	if _, err := w.Write([]byte("hello world")); err != nil {
		return fmt.Errorf("writing object: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	// Delete an object.
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	return nil
}
