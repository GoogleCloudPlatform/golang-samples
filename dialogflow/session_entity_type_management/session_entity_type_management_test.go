// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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

func TestSessionEntityTypeManagement(t *testing.T) {
	t.Skip("flaky")
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID
	sessionID := fmt.Sprintf("test-session-%v", time.Now().Unix())

	// Create and defer closure of entity types client used in this test
	ctx := context.Background()
	entityTypesClient, err := dialogflow.NewEntityTypesClient(ctx)
	if err != nil {
		t.Error(err)
	}
	defer entityTypesClient.Close()

	// Create and defer deletion of entity type used in this test
	displayNames := [...]string{fmt.Sprintf("entityType-%v-1", time.Now().Unix()), fmt.Sprintf("entityType-%v-2", time.Now().Unix())}
	values := [][]string{{"value-1"}, {"value-2-1", "value-2-2"}}
	overrideModes := []string{"OVERRIDE", "SUPPLEMENT"}
	parent := fmt.Sprintf("projects/%s/agent", projectID)
	for _, displayName := range displayNames {
		entityType := dialogflowpb.EntityType{DisplayName: displayName, Kind: dialogflowpb.EntityType_KIND_MAP}
		creationRequest := dialogflowpb.CreateEntityTypeRequest{Parent: parent, EntityType: &entityType}
		response, err := entityTypesClient.CreateEntityType(ctx, &creationRequest)
		if err != nil {
			t.Error(err)
		}
		entityName := response.GetName()
		deletionRequest := dialogflowpb.DeleteEntityTypeRequest{Name: entityName}
		defer entityTypesClient.DeleteEntityType(ctx, &deletionRequest)
	}

	var sessionEntityTypeNames []string

	initialSessionEntityTypes, err := ListSessionEntityTypes(projectID, sessionID)

	if err != nil {
		t.Error("Unsuccessful initial ListSessionEntityTypes")
	}

	for i, displayName := range displayNames {
		name, err := CreateSessionEntityType(projectID, sessionID, displayName, overrideModes[i], values[i])
		if err != nil {
			t.Errorf("Unsuccessful entityType creation: %s", displayName)
		}
		sessionEntityTypeNames = append(sessionEntityTypeNames, name)
	}

	intermediateSessionEntityTypes, err := ListSessionEntityTypes(projectID, sessionID)

	if err != nil {
		t.Error("Unsuccessful intermediate ListSessionEntityTypes")
	}

	if len(intermediateSessionEntityTypes) != len(initialSessionEntityTypes)+len(displayNames) {
		t.Errorf("len(intermediateSessionEntityTypes) = %d; want %d", len(intermediateSessionEntityTypes), len(initialSessionEntityTypes)+len(displayNames))
	}

	for _, name := range sessionEntityTypeNames {
		route := strings.Split(name, "/")
		displayName := route[len(route)-1]
		err = DeleteSessionEntityType(projectID, sessionID, displayName)
		if err != nil {
			t.Errorf("Unsuccessful entityType deletion %s", displayName)
		}
	}

	finalSessionEntityTypes, err := ListSessionEntityTypes(projectID, sessionID)

	if err != nil {
		t.Error("Unsuccessful final ListSessionEntityTypes")
	}

	if len(finalSessionEntityTypes) != len(initialSessionEntityTypes) {
		t.Errorf("Actual len(finalSessionEntityTypes) = %d; want %d", len(finalSessionEntityTypes), len(initialSessionEntityTypes))
	}
}
