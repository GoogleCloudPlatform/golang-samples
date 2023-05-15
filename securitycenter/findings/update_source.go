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

package findings

// [START securitycenter_update_source]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"google.golang.org/genproto/protobuf/field_mask"
)

// updateSource changes a sources display name to "New Display Name" for a
// specific source. sourceName is the full resource name of the source to be
// updated.
func updateSource(w io.Writer, sourceName string) error {
	// sourceName := "organizations/111122222444/sources/1234"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.UpdateSourceRequest{
		Source: &securitycenterpb.Source{
			Name:        sourceName,
			DisplayName: "New Display Name",
		},
		// Only update the display name field (if not set all mutable
		// fields of the source will be updated.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"display_name"},
		},
	}
	source, err := client.UpdateSource(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateSource: %w", err)
	}
	fmt.Fprintf(w, "Source Name: %s, ", source.Name)
	fmt.Fprintf(w, "Display name: %s, ", source.DisplayName)
	fmt.Fprintf(w, "Description: %s\n", source.Description)

	return nil
}

// [END securitycenter_update_source]
