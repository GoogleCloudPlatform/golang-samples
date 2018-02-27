// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> <OPERATION>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project id\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of list, create, delete\n")
	}

	var projectId, sessionId, contextId string
	flag.StringVar(&projectId, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&sessionId, "session-id", "", "Dialogflow session ID")
	flag.StringVar(&contextId, "context-id", "", "Dialogflow context ID")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)

	var err error

	switch operation {
	case "list":
		fmt.Printf("Contexts under projects/%s/agent/sessions/%s:\n", projectId, sessionId)
		var contexts []*dialogflowpb.Context
		contexts, err = listContexts(projectId, sessionId)
		if err != nil {
			log.Fatal(err)
		}
		for _, context := range contexts {
			fmt.Printf("Path: %s, Lifespan: %d\n", context.Name, context.LifespanCount)
		}
	case "create":
		fmt.Printf("Creating context projects/%s/agent/sessions/%s/contexts/%s...\n", projectId, sessionId, contextId)
		err = createContext(projectId, sessionId, contextId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	case "delete":
		fmt.Printf("Deleting context projects/%s/agent/sessions/%s/contexts/%s...\n", projectId, sessionId, contextId)
		err = deleteContext(projectId, sessionId, contextId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func listContexts(projectId string, sessionId string) ([]*dialogflowpb.Context, error) {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if (clientErr != nil) {
		return nil, clientErr
	}

	if (projectId == "" || sessionId == "") {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectId, sessionId))
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, sessionId)

	request := dialogflowpb.ListContextsRequest{Parent: parent}

	contextIterator := contextsClient.ListContexts(ctx, &request)
	var contexts []*dialogflowpb.Context


	for context, status := contextIterator.Next(); status != iterator.Done; {
		contexts = append(contexts, context)
		context, status = contextIterator.Next()
	}

	return contexts, nil
}

func createContext(projectId string, sessionId string, contextId string) error {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if (clientErr != nil) {
		return clientErr
	}

	if (projectId == "" || sessionId == "" || contextId == "") {
		return errors.New(fmt.Sprintf("Received empty project (%s) or session (%s) or context (%s)", projectId, sessionId, contextId))
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, sessionId)
	targetPath := fmt.Sprintf("%s/contexts/%s", parent, contextId)
	target := dialogflowpb.Context{Name: targetPath, LifespanCount: 10}

	request := dialogflowpb.CreateContextRequest{Parent: parent, Context: &target}

	_, requestErr := contextsClient.CreateContext(ctx, &request)
	if (requestErr != nil) {
		return requestErr
	}

	return nil
}

func deleteContext(projectId string, sessionId string, contextId string) error {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if (clientErr != nil) {
		return clientErr
	}

	if (projectId == "" || sessionId == "" || contextId == "") {
		return errors.New(fmt.Sprintf("Received empty project (%s) or session (%s) or context (%s)", projectId, sessionId, contextId))
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, sessionId)
	targetPath := fmt.Sprintf("%s/contexts/%s", parent, contextId)

	request := dialogflowpb.DeleteContextRequest{Name: targetPath}

	requestErr := contextsClient.DeleteContext(ctx, &request)
	if (requestErr != nil) {
		return requestErr
	}

	return nil
}
