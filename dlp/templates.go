/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	"google.golang.org/api/iterator"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_create_template]
func createInspectTemplate(w io.Writer, client *dlp.Client, project string, minLikelihood dlppb.Likelihood, maxFindings int32, templateID, displayName, description string, infoTypes []string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	req := &dlppb.CreateInspectTemplateRequest{
		Parent:     "projects/" + project,
		TemplateId: templateID,
		InspectTemplate: &dlppb.InspectTemplate{
			DisplayName: displayName,
			Description: description,
			InspectConfig: &dlppb.InspectConfig{
				InfoTypes:     i,
				MinLikelihood: minLikelihood,
				Limits: &dlppb.InspectConfig_FindingLimits{
					MaxFindingsPerRequest: maxFindings,
				},
			},
		},
	}
	r, err := client.CreateInspectTemplate(context.Background(), req)
	if err != nil {
		log.Fatalf("error creating inspect template: %v", err)
	}
	fmt.Fprintf(w, "Successfully created inspect template: %v", r.GetName())
}

// [END dlp_create_template]

// [START dlp_list_templates]
func listInspectTemplates(w io.Writer, client *dlp.Client, project string) {
	req := &dlppb.ListInspectTemplatesRequest{
		Parent: "projects/" + project,
	}
	it := client.ListInspectTemplates(context.Background(), req)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error getting inspect templates: %v", err)
		}
		c := t.GetCreateTime()
		u := t.GetUpdateTime()
		fmt.Fprintf(w, "Inspect template %v\n", t.GetName())
		fmt.Fprintf(w, "  Created: %v\n", time.Unix(c.GetSeconds(), int64(c.GetNanos())).Format(time.RFC1123))
		fmt.Fprintf(w, "  Updated: %v\n", time.Unix(u.GetSeconds(), int64(u.GetNanos())).Format(time.RFC1123))
		fmt.Fprintf(w, "  Display Name: %q\n", t.GetDisplayName())
		fmt.Fprintf(w, "  Description: %q\n", t.GetDescription())
	}
}

// [END dlp_list_templates]

// [START dlp_delete_template]
func deleteInspectTemplate(w io.Writer, client *dlp.Client, templateID string) {
	req := &dlppb.DeleteInspectTemplateRequest{
		Name: templateID,
	}
	err := client.DeleteInspectTemplate(context.Background(), req)
	if err != nil {
		log.Fatalf("error deleting inspect template: %v", err)
	}
	fmt.Fprintf(w, "Successfully deleted inspect template %v", templateID)
}

// [END dlp_delete_template]
