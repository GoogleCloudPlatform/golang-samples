// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestEntityTypeManagement(t *testing.T) {
	t.Skip("flaky")
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID

	parent := fmt.Sprintf("projects/%s/agents", projectID)

	displayNames := [...]string{fmt.Sprintf("entityType-%v-1", time.Now().Unix()), fmt.Sprintf("entityType-%v-2", time.Now().Unix())}
	var entityTypeNames []string

	initialEntityTypes, err := ListEntityTypes(projectID)

	if err != nil {
		t.Error("Unsuccessful initial ListEntityTypes")
	}

	for _, displayName := range displayNames {
		name, err := CreateEntityType(projectID, displayName, "KIND_MAP")
		if err != nil {
			t.Errorf("Unsuccessful entityType creation: %s", displayName)
		}
		entityTypeNames = append(entityTypeNames, name)
	}

	intermediateEntityTypes, err := ListEntityTypes(projectID)

	if err != nil {
		t.Error("Unsuccessful intermediate ListEntityTypes")
	}

	if len(intermediateEntityTypes) != len(initialEntityTypes)+len(displayNames) {
		t.Errorf("len(intermediateEntityTypes) = %d; want %d", len(intermediateEntityTypes), len(initialEntityTypes)+len(displayNames))
	}

	for _, entityTypeName := range entityTypeNames {
		route := strings.Split(entityTypeName, "/")
		entityTypeID := route[len(route)-1]
		err = DeleteEntityType(projectID, entityTypeID)
		if err != nil {
			t.Errorf("Unsuccessful entityType deletion %s/entityType/%s", parent, entityTypeID)
		}
	}

	finalEntityTypes, err := ListEntityTypes(projectID)

	if err != nil {
		t.Error("Unsuccessful final ListEntityTypes")
	}

	if len(finalEntityTypes) != len(initialEntityTypes) {
		t.Errorf("Actual len(finalEntityTypes) = %d; want %d", len(finalEntityTypes), len(initialEntityTypes))
	}
}
