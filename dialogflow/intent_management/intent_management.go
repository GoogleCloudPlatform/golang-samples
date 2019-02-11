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
	"google.golang.org/api/iterator"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// [END import_libraries]

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> <OPERATION> <SUBCOMMAND ARGUMENTS>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project ID\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of list, create, delete\n")
		fmt.Fprintf(os.Stderr, "<SUBCOMMAND ARGUMENTS> can be passed if <OPERATION> is create; pass with flags -training-phrases-parts <PART_1>,<PART_2>,...,<PART_M> -message-texts=<TEXT_1>,<TEXT_2>,...,<TEXT_N>, where <PARTS_i> and <TEXT_j> are strings\n")
	}

	var projectID string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	operation := flag.Arg(0)

	var err error

	switch operation {
	case "list":
		fmt.Printf("Intents under projects/%s/agent:\n", projectID)
		var intents []*dialogflowpb.Intent
		intents, err = ListIntents(projectID)
		if err != nil {
			log.Fatal(err)
		}

		// Ugly code for beautiful output
		for _, intent := range intents {
			fmt.Printf("Intent name: %s\nDisplay name: %s\n", intent.GetName(), intent.GetDisplayName())
			fmt.Printf("Action: %s\n", intent.GetAction())
			fmt.Printf("Root followup intent: %s\nParent followup intent: %s\n", intent.GetRootFollowupIntentName(), intent.GetParentFollowupIntentName())
			fmt.Printf("Input contexts: %s\n", strings.Join(intent.GetInputContextNames(), ", "))
			fmt.Println("Output contexts:")
			for _, outputContext := range intent.GetOutputContexts() {
				fmt.Printf("\tName: %s\n", outputContext.GetName())
			}
			fmt.Println("---")
		}
	case "create":
		creationFlagSet := flag.NewFlagSet("create", flag.ExitOnError)
		var trainingPhrasesPartsRaw, messageTextsRaw string
		creationFlagSet.StringVar(&trainingPhrasesPartsRaw, "training-phrases-parts", "", "Parts of phrases associated with the intent you are creating")
		creationFlagSet.StringVar(&messageTextsRaw, "message-texts", "", "Messages that the Dialogflow agent should respond to the intent with")

		creationFlagSet.Parse(flag.Args()[1:])
		creationArgs := creationFlagSet.Args()
		if len(creationArgs) != 1 {
			log.Fatalf("Please pass a display name for the intent you wish to create")
		}

		displayName := creationArgs[0]
		trainingPhrasesParts := strings.Split(trainingPhrasesPartsRaw, ",")
		messageTexts := strings.Split(messageTextsRaw, ",")

		fmt.Printf("Creating intent %s under projects/%s/agent...\n", displayName, projectID)
		err = CreateIntent(projectID, displayName, trainingPhrasesParts, messageTexts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	case "delete":
		deletionFlagSet := flag.NewFlagSet("delete", flag.ExitOnError)
		var intentID string
		deletionFlagSet.StringVar(&intentID, "intent-id", "", "Path to intent you would like to delete")

		deletionFlagSet.Parse(flag.Args()[1:])

		if intentID == "" {
			log.Fatal("Expected non-empty -intention-id argument")
		}

		fmt.Printf("Deleting intent projects/%s/agent/intents/%s...\n", projectID, intentID)
		err = DeleteIntent(projectID, intentID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Done!\n")
	default:
		flag.Usage()
		os.Exit(1)
	}
}

// [START dialogflow_list_intents]

func ListIntents(projectID string) ([]*dialogflowpb.Intent, error) {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx)
	if clientErr != nil {
		return nil, clientErr
	}
	defer intentsClient.Close()

	if projectID == "" {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s)", projectID))
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	request := dialogflowpb.ListIntentsRequest{Parent: parent}

	intentIterator := intentsClient.ListIntents(ctx, &request)
	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

// [END dialogflow_list_intents]

// [START dialogflow_create_intent]
func CreateIntent(projectID, displayName string, trainingPhraseParts, messageTexts []string) error {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer intentsClient.Close()

	if projectID == "" || displayName == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or intent (%s)", projectID, displayName))
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	var targetTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase
	var targetTrainingPhraseParts []*dialogflowpb.Intent_TrainingPhrase_Part
	for _, partString := range trainingPhraseParts {
		part := dialogflowpb.Intent_TrainingPhrase_Part{Text: partString}
		targetTrainingPhraseParts = []*dialogflowpb.Intent_TrainingPhrase_Part{&part}
		targetTrainingPhrase := dialogflowpb.Intent_TrainingPhrase{Type: dialogflowpb.Intent_TrainingPhrase_EXAMPLE, Parts: targetTrainingPhraseParts}
		targetTrainingPhrases = append(targetTrainingPhrases, &targetTrainingPhrase)
	}

	intentMessageTexts := dialogflowpb.Intent_Message_Text{Text: messageTexts}
	wrappedIntentMessageTexts := dialogflowpb.Intent_Message_Text_{Text: &intentMessageTexts}
	intentMessage := dialogflowpb.Intent_Message{Message: &wrappedIntentMessageTexts}

	target := dialogflowpb.Intent{DisplayName: displayName, WebhookState: dialogflowpb.Intent_WEBHOOK_STATE_UNSPECIFIED, TrainingPhrases: targetTrainingPhrases, Messages: []*dialogflowpb.Intent_Message{&intentMessage}}

	request := dialogflowpb.CreateIntentRequest{Parent: parent, Intent: &target}

	_, requestErr := intentsClient.CreateIntent(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

// [END dialogflow_create_intent]

// [START dialogflow_delete_intent]
func DeleteIntent(projectID, intentID string) error {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx)
	if clientErr != nil {
		return clientErr
	}
	defer intentsClient.Close()

	if projectID == "" || intentID == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or intent (%s)", projectID, intentID))
	}

	targetPath := fmt.Sprintf("projects/%s/agent/intents/%s", projectID, intentID)

	request := dialogflowpb.DeleteIntentRequest{Name: targetPath}

	requestErr := intentsClient.DeleteIntent(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

// [END dialogflow_delete_intent]
