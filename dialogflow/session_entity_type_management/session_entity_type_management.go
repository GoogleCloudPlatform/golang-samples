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

	var projectID, sessionID string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&sessionID, "session-id", "", "Dialogflow session ID")

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)

	var err error

	switch operation {
	case "list":
		fmt.Printf("SessionEntityTypes under projects/%s/agent/sessions/%s:\n", projectID, sessionID)
		var sessionEntityTypes []*dialogflowpb.SessionEntityType
		sessionEntityTypes, err = ListSessionEntityTypes(projectID, sessionID)
		if err != nil {
			log.Fatal(err)
		}
		for _, sessionEntityType := range sessionEntityTypes {
			overrideMode := int32(sessionEntityType.GetEntityOverrideMode())
			overrideModeString := dialogflowpb.SessionEntityType_EntityOverrideMode_name[overrideMode]
			fmt.Printf("Path: %s, Entity override mode: %s\n", sessionEntityType.GetName(), overrideModeString)
			fmt.Printf("Entities:\n")
			for _, entity := range sessionEntityType.GetEntities() {
				fmt.Printf("\t%s\n", entity.GetValue())
			}
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func ListSessionEntityTypes(projectID, sessionID string) ([]*dialogflowpb.SessionEntityType, error) {
	ctx := context.Background()

	sessionEntityTypesClient, clientErr := dialogflow.NewSessionEntityTypesClient(ctx)
	if clientErr != nil {
		return nil, clientErr
	}
	defer sessionEntityTypesClient.Close()

	if projectID == "" || sessionID == "" {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	request := dialogflowpb.ListSessionEntityTypesRequest{Parent: parent}

	sessionEntityTypeIterator := sessionEntityTypesClient.ListSessionEntityTypes(ctx, &request)
	var sessionEntityTypes []*dialogflowpb.SessionEntityType

	for sessionEntityType, status := sessionEntityTypeIterator.Next(); status != iterator.Done; {
		sessionEntityTypes = append(sessionEntityTypes, sessionEntityType)
		sessionEntityType, status = sessionEntityTypeIterator.Next()
	}

	return sessionEntityTypes, nil
}
