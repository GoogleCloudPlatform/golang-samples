// Copyright 2024 Google LLC
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

package firestore

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	apiv1 "cloud.google.com/go/firestore/apiv1/admin"
	"cloud.google.com/go/firestore/apiv1/admin/adminpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMultipleInequalities(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	projectID = tc.ProjectID
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		client.Close()
	})

	cities := []City{
		{
			Name:       "San Francisco",
			State:      "CA",
			Country:    "USA",
			Capital:    false,
			Population: 860000,
			Density:    18000,
			Regions:    []string{"west_coast", "norcal"},
		},
		{
			Name:       "Los Angeles",
			State:      "CA",
			Country:    "USA",
			Capital:    false,
			Population: 3900000,
			Density:    8300,
			Regions:    []string{"west_coast", "socal"},
		},
		{
			Name:       "Washington D.C.",
			Country:    "USA",
			Capital:    true,
			Population: 680000,
			Density:    11300,
			Regions:    []string{"east_coast"},
		},
		{
			Name:       "Tokyo",
			Country:    "Japan",
			Capital:    true,
			Population: 9000000,
			Density:    16000,
			Regions:    []string{"kanto", "honshu"},
		},
		{
			Name:       "Beijing",
			Country:    "China",
			Capital:    true,
			Population: 21500000,
			Density:    3500,
			Regions:    []string{"jingjinji", "hebei"},
		},
	}

	// Populate data
	var refs []*firestore.DocumentRef
	colName := "cities"
	bw := client.BulkWriter(ctx)
	for _, city := range cities {
		ref := client.Collection(colName).NewDoc()
		if _, err := bw.Create(ref, city); err != nil {
			t.Fatal(err)
		}
		refs = append(refs, ref)
	}
	bw.End()
	t.Cleanup(func() {
		// New BulkWriter instance
		bw = client.BulkWriter(ctx)

		for _, r := range refs {
			if _, err := bw.Delete(r); err != nil {
				t.Fatal(err)
			}
		}
		bw.End()

	})

	// Create indexes for multiple inequality query
	indexParent := "projects/" + tc.ProjectID + "/databases/(default)/collectionGroups/" + colName
	adminPbIndexFields := []*adminpb.Index_IndexField{
		{
			FieldPath: "density",
			ValueMode: &adminpb.Index_IndexField_Order_{
				Order: adminpb.Index_IndexField_ASCENDING,
			},
		},
		{
			FieldPath: "population",
			ValueMode: &adminpb.Index_IndexField_Order_{
				Order: adminpb.Index_IndexField_ASCENDING,
			},
		},
	}
	req := &adminpb.CreateIndexRequest{
		Parent: indexParent,
		Index: &adminpb.Index{
			QueryScope: adminpb.Index_COLLECTION,
			Fields:     adminPbIndexFields,
		},
	}
	adminClient, err := apiv1.NewFirestoreAdminClient(ctx)
	if err != nil {
		t.Fatalf("NewFirestoreAdminClient: %v", err)
	}
	t.Cleanup(func() {
		adminClient.Close()
	})
	op, createErr := adminClient.CreateIndex(ctx, req)
	if createErr != nil {
		t.Fatalf("CreateIndex: %v", createErr)
	}
	createdIndex, waitErr := op.Wait(ctx)
	if waitErr != nil {
		t.Fatalf("CreateIndexes failed. Wait: %v", waitErr)
	}
	t.Cleanup(func() {
		if err = adminClient.DeleteIndex(ctx, &adminpb.DeleteIndexRequest{
			Name: createdIndex.Name,
		}); err != nil {
			t.Fatalf("Failed to delete index \"%s\": %+v\n", createdIndex.Name, err)
		}
	})

	// Run sample and capture console output
	buf := new(bytes.Buffer)
	if err = multipleInequalities(buf, projectID); err != nil {
		t.Errorf("multipleInequalities: %v", err)
	}

	// Compare console outputs
	got := buf.String()
	want := "map[" +
		"capital:true country:China density:3500 name:Beijing population:21500000 regions:[jingjinji hebei]" +
		"]\nmap[" +
		"country:USA density:8300 name:Los Angeles population:3900000 regions:[west_coast socal] state:CA]\n"
	if !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
