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

// [START monitoring_irm_create_incident]
import (
	"context"
	"fmt"
	"io"

	irm "cloud.google.com/go/irm/apiv1alpha2"
	irmpb "google.golang.org/genproto/googleapis/cloud/irm/v1alpha2"
)

// createIncident creates an incident.
func createIncident(w io.Writer, projectID string) (*irmpb.Incident, error) {
	ctx := context.Background()

	client, err := irm.NewIncidentClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("irm.NewIncidentClient: %v", err)
	}
	defer client.Close()

	req := &irmpb.CreateIncidentRequest{
		Parent: "projects/" + projectID,
		Incident: &irmpb.Incident{
			Title: "Somebody pushed the red button!",
			Synopsis: &irmpb.Synopsis{
				Author: &irmpb.User{
					User: &irmpb.User_Email{
						Email: "janedoe@example.com",
					},
				},
				Content:     "Nobody should ever push the red button.",
				ContentType: "text/plain",
				// UpdateTime:  time.Now(),
			},
			Severity: irmpb.Incident_SEVERITY_MAJOR,
			Stage:    irmpb.Incident_STAGE_UNSPECIFIED,
		},
	}

	incident, err := client.CreateIncident(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateIncident: %v", err)
	}

	fmt.Fprintf(w, "Created incident: %q", incident.Name)

	return incident, nil
}

// [END monitoring_irm_create_incident]
