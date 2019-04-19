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
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> -session-id <SESSION ID> -context-id <CONTEXT ID> <OPERATION>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project ID\n")
		fmt.Fprintf(os.Stderr, "<SESSION ID> must be a Dialogflow session ID\n")
		fmt.Fprintf(os.Stderr, "<CONTEXT ID> must be a Dialogflow context ID - only required for create and delete calls\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of list, create, delete\n")
	}

	var projectID, sessionID, contextID string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&sessionID, "session-id", "", "Dialogflow session ID")
	flag.StringVar(&contextID, "context-id", "", "Dialogflow context ID")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)

	var err error

	switch operation {
	case "list":
		fmt.Printf("Contexts under projects/%s/agent/sessions/%s:\n", projectID, sessionID)
		var contexts []*dialogflowpb.Context
		contexts, err = ListContexts(projectID, sessionID)
		if err != nil {
			log.Fatal(err)
		}
		for _, context := range contexts {
			fmt.Printf("Path: %s, Lifespan: %d\n", context.Name, context.LifespanCount)
		}
	case "create":
		fmt.Printf("Creating context projects/%s/agent/sessions/%s/contexts/%s...\n", projectID, sessionID, contextID)
		err = CreateContext(projectID, sessionID, contextID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	case "delete":
		fmt.Printf("Deleting context projects/%s/agent/sessions/%s/contexts/%s...\n", projectID, sessionID, contextID)
		err = DeleteContext(projectID, sessionID, contextID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func ListContexts(projectID, sessionID string) ([]*dialogflowpb.Context, error) {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if clientErr != nil {
		return nil, clientErr
	}
	defer contextsClient.Close()

	if projectID == "" || sessionID == "" {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	request := dialogflowpb.ListContextsRequest{Parent: parent}

	contextIterator := contextsClient.ListContexts(ctx, &request)
	var contexts []*dialogflowpb.Context

	for context, status := contextIterator.Next(); status != iterator.Done; {
		contexts = append(contexts, context)
		context, status = contextIterator.Next()
	}

	return contexts, nil
}

// [START dialogflow_create_context]
func CreateContext(projectID, sessionID, contextID string) error {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer contextsClient.Close()

	if projectID == "" || sessionID == "" || contextID == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or session (%s) or context (%s)", projectID, sessionID, contextID))
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	targetPath := fmt.Sprintf("%s/contexts/%s", parent, contextID)
	target := dialogflowpb.Context{Name: targetPath, LifespanCount: 10}

	request := dialogflowpb.CreateContextRequest{Parent: parent, Context: &target}

	_, requestErr := contextsClient.CreateContext(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

// [END dialogflow_create_context]

// [START dialogflow_delete_context]
func DeleteContext(projectID, sessionID, contextID string) error {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer contextsClient.Close()

	if projectID == "" || sessionID == "" || contextID == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or session (%s) or context (%s)", projectID, sessionID, contextID))
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	targetPath := fmt.Sprintf("%s/contexts/%s", parent, contextID)

	request := dialogflowpb.DeleteContextRequest{Name: targetPath}

	requestErr := contextsClient.DeleteContext(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

// [END dialogflow_delete_context]
