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
	"strings"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// [END import_libraries]

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> <OPERATION> <ADDITIONAL ARGUMENTS>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project ID\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of list, create, delete\n")
		fmt.Fprintf(os.Stderr, "<ADDITIONAL ARGUMENTS> For the create subcommand, you are expected to pass -synonyms, as well as an entity value. For the delete subcommand, you are expected to pass an entity value.\n")
	}

	var projectID, entityTypeID string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&entityTypeID, "entity-type-id", "", "Unique ID of entity type corresponding to the entity/entities you are working with")

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)

	var err error

	switch operation {
	case "list":
		fmt.Printf("Entities under projects/%s/agent:\n", projectID)
		var entities []*dialogflowpb.EntityType_Entity
		entities, err = ListEntities(projectID, entityTypeID)
		if err != nil {
			log.Fatal(err)
		}
		for _, entity := range entities {
			fmt.Printf("Value: %s\n", entity.GetValue())
			fmt.Println("Synonyms:")
			for _, synonym := range entity.GetSynonyms() {
				fmt.Printf("\t- %s\n", synonym)
			}
			fmt.Println("")
		}
	case "create":
		creationFlagSet := flag.NewFlagSet("create", flag.ExitOnError)
		var synonymsRaw string
		creationFlagSet.StringVar(&synonymsRaw, "synonyms", "", "Comma-separated list of synonyms for the given entity: <SYNONYM_1>,<SYNONYM_2>,...,<SYNONYM_N>")
		creationFlagSet.Parse(flag.Args()[1:])

		if len(creationFlagSet.Args()) != 1 {
			log.Fatal("No entity value passed to create")
		}
		entityValue := creationFlagSet.Arg(0)
		synonyms := strings.Split(synonymsRaw, ",")

		fmt.Printf("Creating entity %s...\n", entityValue)
		err := CreateEntity(projectID, entityTypeID, entityValue, synonyms)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Entity type %s created under type %s\n", entityValue, entityTypeID)
	case "delete":
		if len(flag.Args()) != 2 {
			log.Fatal("No entity value passed to delete")
		}
		entityValue := flag.Arg(1)

		fmt.Printf("Deleting values %s under projects/%s/agent/entityTypes/%s...\n", entityValue, projectID, entityTypeID)
		err = DeleteEntity(projectID, entityTypeID, entityValue)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func ListEntities(projectID, entityTypeID string) ([]*dialogflowpb.EntityType_Entity, error) {
	ctx := context.Background()

	entityTypesClient, clientErr := dialogflow.NewEntityTypesClient(ctx)
	if clientErr != nil {
		return nil, clientErr
	}
	defer entityTypesClient.Close()

	if projectID == "" || entityTypeID == "" {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s) or entity type (%s)", projectID, entityTypeID))
	}

	entityName := fmt.Sprintf("projects/%s/agent/entityTypes/%s", projectID, entityTypeID)

	request := dialogflowpb.GetEntityTypeRequest{Name: entityName}

	entityType, err := entityTypesClient.GetEntityType(ctx, &request)
	if err != nil {
		return []*dialogflowpb.EntityType_Entity{}, err
	}

	return entityType.GetEntities(), nil
}

// [START dialogflow_create_entity]
func CreateEntity(projectID, entityTypeID, entityValue string, synonyms []string) error {
	ctx := context.Background()

	entityTypesClient, clientErr := dialogflow.NewEntityTypesClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer entityTypesClient.Close()

	if projectID == "" || entityTypeID == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or entity type (%s)", projectID, entityTypeID))
	}

	parent := fmt.Sprintf("projects/%s/agent/entityTypes/%s", projectID, entityTypeID)
	entity := dialogflowpb.EntityType_Entity{Value: entityValue, Synonyms: synonyms}
	entities := []*dialogflowpb.EntityType_Entity{&entity}

	request := dialogflowpb.BatchCreateEntitiesRequest{Parent: parent, Entities: entities}

	creationOp, requestErr := entityTypesClient.BatchCreateEntities(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	err := creationOp.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

// [END dialogflow_create_entity]

// [START dialogflow_delete_entity]
func DeleteEntity(projectID, entityTypeID, entityValue string) error {
	ctx := context.Background()

	entityTypesClient, clientErr := dialogflow.NewEntityTypesClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer entityTypesClient.Close()

	if projectID == "" || entityTypeID == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or entity type (%s)", projectID, entityTypeID))
	}

	parent := fmt.Sprintf("projects/%s/agent/entityTypes/%s", projectID, entityTypeID)
	entityValues := []string{entityValue}
	request := dialogflowpb.BatchDeleteEntitiesRequest{Parent: parent, EntityValues: entityValues}

	deletionOp, requestErr := entityTypesClient.BatchDeleteEntities(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	err := deletionOp.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

// [END dialogflow_delete_entity]
