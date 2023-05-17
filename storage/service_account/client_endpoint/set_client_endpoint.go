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
	// opts := []option.ClientOption{}
	ctx := context.Background()

	// Add the custom endpoint option to any other desired options passed to storage.NewClient.
	opts = append(opts, option.WithEndpoint(customEndpoint))
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Use the client as per your custom endpoint, for example, attempt to get a bucket's metadata.
	client.Bucket("bucket-name").Attrs(ctx)
	return nil
}

// [END storage_set_client_endpoint]
