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

package metadata

// [START dlp_list_info_types]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// infoTypes returns the info types in the given language and matching the given filter.
func infoTypes(w io.Writer, languageCode, filter string) error {
	// languageCode := "en-US"
	// filter := "supported_by=INSPECT"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	req := &dlppb.ListInfoTypesRequest{
		LanguageCode: languageCode,
		Filter:       filter,
	}
	resp, err := client.ListInfoTypes(ctx, req)
	if err != nil {
		return fmt.Errorf("ListInfoTypes: %w", err)
	}
	for _, it := range resp.GetInfoTypes() {
		fmt.Fprintln(w, it.GetName())
	}
	return nil
}

// [END dlp_list_info_types]
