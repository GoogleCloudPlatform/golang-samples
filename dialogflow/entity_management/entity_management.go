// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

// [START import_libraries]
import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
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
