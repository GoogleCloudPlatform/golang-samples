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

// [START import_libraries]
import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/iterator"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// [END import_libraries]

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> <OPERATION> <ADDITIONAL ARGUMENTS>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project ID\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of list, create, delete\n")
		fmt.Fprintf(os.Stderr, "<ADDITIONAL ARGUMENTS> must be a display name in the case of the create subcommand and an entity type ID in the case of the delete subcommand\n")
	}

	var projectID string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)

	var err error

	switch operation {
	case "list":
		fmt.Printf("EntityTypes under projects/%s/agent:\n", projectID)
		var entityTypes []*dialogflowpb.EntityType
		entityTypes, err = ListEntityTypes(projectID)
		if err != nil {
			log.Fatal(err)
		}
		for _, entityType := range entityTypes {
			fmt.Printf("Path: %s, Display name: %s, Kind: %s\n", entityType.GetName(), entityType.GetDisplayName(), entityType.GetKind())
		}
	case "create":
		creationFlagSet := flag.NewFlagSet("create", flag.ExitOnError)
		var kind string
		creationFlagSet.StringVar(&kind, "kind", "KIND_MAP", "Should be either KIND_MAP (default) or KIND_LIST")
		creationFlagSet.Parse(flag.Args()[1:])

		if len(creationFlagSet.Args()) != 1 {
			log.Fatal("The create subcommand should be called with a single display name")
		}
		displayName := creationFlagSet.Arg(0)

		fmt.Printf("Creating entityType %s...\n", displayName)
		entityTypeName, err := CreateEntityType(projectID, displayName, kind)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Entity type %s created as %s\n", displayName, entityTypeName)
	case "delete":
		if len(flag.Args()[1:]) != 1 {
			log.Fatal("The delete subcommand should be called with a single entity type ID")
		}
		entityTypeID := flag.Arg(1)

		fmt.Printf("Deleting entityType projects/%s/agent/entityTypes/%s...\n", projectID, entityTypeID)
		err = DeleteEntityType(projectID, entityTypeID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func ListEntityTypes(projectID string) ([]*dialogflowpb.EntityType, error) {
	ctx := context.Background()

	entityTypesClient, clientErr := dialogflow.NewEntityTypesClient(ctx)
	if clientErr != nil {
		return nil, clientErr
	}
	defer entityTypesClient.Close()

	if projectID == "" {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s)", projectID))
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	request := dialogflowpb.ListEntityTypesRequest{Parent: parent}

	entityTypeIterator := entityTypesClient.ListEntityTypes(ctx, &request)
	var entityTypes []*dialogflowpb.EntityType

	for entityType, status := entityTypeIterator.Next(); status != iterator.Done; {
		entityTypes = append(entityTypes, entityType)
		entityType, status = entityTypeIterator.Next()
	}

	return entityTypes, nil
}

// [START dialogflow_create_entity_type]
func CreateEntityType(projectID, displayName, kind string) (string, error) {
	ctx := context.Background()

	entityTypesClient, clientErr := dialogflow.NewEntityTypesClient(ctx)
	if clientErr != nil {
		return "", clientErr
	}
	defer entityTypesClient.Close()

	if projectID == "" || displayName == "" {
		return "", errors.New(fmt.Sprintf("Received empty project (%s) or displayName (%s)", projectID, displayName))
	}

	var kindValue dialogflowpb.EntityType_Kind
	switch kind {
	case "KIND_MAP":
		kindValue = dialogflowpb.EntityType_KIND_MAP
	case "KIND_LIST":
		kindValue = dialogflowpb.EntityType_KIND_LIST
	default:
		return "", errors.New(fmt.Sprintf("Received invalid kind argument: %s; acceptable values are \"KIND_MAP\" and \"KIND_LIST\"", kind))
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)
	target := dialogflowpb.EntityType{DisplayName: displayName, Kind: kindValue}

	request := dialogflowpb.CreateEntityTypeRequest{Parent: parent, EntityType: &target}

	response, requestErr := entityTypesClient.CreateEntityType(ctx, &request)
	if requestErr != nil {
		return "", requestErr
	}

	return response.GetName(), nil
}

// [END dialogflow_delete_entity_type]

// [START dialogflow_delete_entity_type]
func DeleteEntityType(projectID, entityTypeID string) error {
	ctx := context.Background()

	entityTypesClient, clientErr := dialogflow.NewEntityTypesClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer entityTypesClient.Close()

	if projectID == "" || entityTypeID == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or entityType (%s)", projectID, entityTypeID))
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)
	targetPath := fmt.Sprintf("%s/entityTypes/%s", parent, entityTypeID)

	request := dialogflowpb.DeleteEntityTypeRequest{Name: targetPath}

	requestErr := entityTypesClient.DeleteEntityType(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

// [END dialogflow_delete_entity_type]
