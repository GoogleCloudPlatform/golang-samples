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

package snippets

// [START monitoring_irm_create_signal]
import (
	"context"
	"fmt"
	"io"

	irm "cloud.google.com/go/irm/apiv1alpha2"
	irmpb "google.golang.org/genproto/googleapis/cloud/irm/v1alpha2"
)

// createSignal creates a new signal.
func createSignal(w io.Writer, projectID string) error {
	ctx := context.Background()

	client, err := irm.NewIncidentClient(ctx)
	if err != nil {
		return fmt.Errorf("irm.NewIncidentClient: %v", err)
	}
	defer client.Close()

	req := &irmpb.CreateSignalRequest{
		Parent: "projects/" + projectID,
		Signal: &irmpb.Signal{
			Title:       "Red button pushed.",
			Content:     "Somebody pushed the red button!",
			ContentType: "text/plain",
		},
	}

	signal, err := client.CreateSignal(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateSignal: %v", err)
	}

	fmt.Fprintf(w, "Created signal: %q", signal.Name)

	return nil
}

// [END monitoring_irm_create_signal]
