// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
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

	var project, session, context string
	flag.StringVar(&project, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&session, "session-id", "", "Dialogflow session ID")
	flag.StringVar(&context, "context-id", "", "Dialogflow context ID")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)

	switch operation {
	case "list":
		listContexts(project, session)
	case "create":
		createContext(project, session, context)
	case "delete":
		deleteContext(project, session, context)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func listContexts(projectId string, sessionId string) {
	ctx := context.Background()

	contextsClient, err := dialogflow.NewContextsClient(ctx)
	if (err != nil) {
		log.Fatal(err)
	}

	if (projectId == "" || sessionId == "") {
		log.Fatalf("Received empty project (%s) or session (%s)", projectId, sessionId)
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, sessionId)

	request := dialogflowpb.ListContextsRequest{Parent: parent}

	contextIterator := contextsClient.ListContexts(ctx, &request)

	fmt.Printf("Contexts under %s:\n", parent)

	for context, status := contextIterator.Next(); status != iterator.Done; {
		fmt.Printf("%v\n", context)
		context, status = contextIterator.Next()
	}
}

func createContext(projectId string, sessionId string, contextId string) {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if (clientErr != nil) {
		log.Fatal(clientErr)
	}

	if (projectId == "" || sessionId == "" || contextId == "") {
		log.Fatalf("Received empty project (%s) or session (%s) or context (%s)", projectId, sessionId, contextId)
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, sessionId)
	targetPath := fmt.Sprintf("%s/contexts/%s", parent, contextId)
	target := dialogflowpb.Context{Name: targetPath, LifespanCount: 10}

	request := dialogflowpb.CreateContextRequest{Parent: parent, Context: &target}

	fmt.Printf("Creating context %s...\n", targetPath)
	response, requestErr := contextsClient.CreateContext(ctx, &request)
	if (requestErr != nil) {
		log.Fatal(requestErr)
	}
	fmt.Printf("Context created: %v\n", response)
}

func deleteContext(projectId string, sessionId string, contextId string) {
	ctx := context.Background()

	contextsClient, clientErr := dialogflow.NewContextsClient(ctx)
	if (clientErr != nil) {
		log.Fatal(clientErr)
	}

	if (projectId == "" || sessionId == "" || contextId == "") {
		log.Fatalf("Received empty project (%s) or session (%s) or context (%s)", projectId, sessionId, contextId)
	}

	parent := fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, sessionId)
	targetPath := fmt.Sprintf("%s/contexts/%s", parent, contextId)

	request := dialogflowpb.DeleteContextRequest{Name: targetPath}

	fmt.Printf("Deleting context %s...\n", targetPath)

	requestErr := contextsClient.DeleteContext(ctx, &request)
	if (requestErr != nil) {
		log.Fatal(requestErr)
	}
	fmt.Printf("Context deleted: %s\n", targetPath)
}
