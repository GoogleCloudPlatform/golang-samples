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

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestEntityManagement(t *testing.T) {
	t.Skip("flaky")
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID
	agent := fmt.Sprintf("projects/%s/agent", projectID)

	// Create and defer closure of entity types client used in this test
	ctx := context.Background()
	entityTypesClient, err := dialogflow.NewEntityTypesClient(ctx)
	if err != nil {
		t.Error(err)
	}
	defer entityTypesClient.Close()

	// Create and defer deletion of entity type used in this test
	displayName := fmt.Sprintf("test-entity-type-%v", time.Now().Unix())
	entityType := dialogflowpb.EntityType{DisplayName: displayName, Kind: dialogflowpb.EntityType_KIND_MAP}
	creationRequest := dialogflowpb.CreateEntityTypeRequest{Parent: agent, EntityType: &entityType}
	response, err := entityTypesClient.CreateEntityType(ctx, &creationRequest)
	if err != nil {
		t.Error(err)
	}
	entityName := response.GetName()
	deletionRequest := dialogflowpb.DeleteEntityTypeRequest{Name: entityName}
	defer entityTypesClient.DeleteEntityType(ctx, &deletionRequest)

	route := strings.Split(entityName, "/")
	entityTypeID := route[len(route)-1]

	parent := fmt.Sprintf("%s/entityTypes/%s", agent, entityTypeID)

	entityValues := [...]string{fmt.Sprintf("entityValue-%v-1", time.Now().Unix()), fmt.Sprintf("entityValue-%v-2", time.Now().Unix())}
	synonyms := [][]string{{"synonym-1-1", "synonym-1-2"}, {"synonym-2-1"}}

	initialEntities, err := ListEntities(projectID, entityTypeID)

	if err != nil {
		t.Error("Unsuccessful initial ListEntities")
	}

	for i, entityValue := range entityValues {
		err := CreateEntity(projectID, entityTypeID, entityValue, synonyms[i])
		if err != nil {
			t.Errorf("Unsuccessful entity creation: %s", entityValue)
		}
	}

	intermediateEntities, err := ListEntities(projectID, entityTypeID)

	if err != nil {
		t.Error("Unsuccessful intermediate ListEntities")
	}

	if len(intermediateEntities) != len(initialEntities)+len(entityValues) {
		t.Errorf("len(intermediateEntities) = %d; want %d", len(intermediateEntities), len(initialEntities)+len(entityValues))
	}

	for _, entityValue := range entityValues {
		err = DeleteEntity(projectID, entityTypeID, entityValue)
		if err != nil {
			t.Errorf("Unsuccessful entity deletion %s under %s", entityValue, parent)
		}
	}

	finalEntities, err := ListEntities(projectID, entityTypeID)

	if err != nil {
		t.Error("Unsuccessful final ListEntities")
	}

	if len(finalEntities) != len(initialEntities) {
		t.Errorf("Actual len(finalEntities) = %d; want %d", len(finalEntities), len(initialEntities))
	}
}
