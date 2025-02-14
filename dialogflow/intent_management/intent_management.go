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

package intentmgmt

import (
	"context"
	"errors"
	"fmt"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"google.golang.org/api/iterator"
)

// [START dialogflow_list_intents]

func ListIntents(projectID string) ([]*dialogflowpb.Intent, error) {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx)
	if clientErr != nil {
		return nil, clientErr
	}
	defer intentsClient.Close()

	if projectID == "" {
		return nil, fmt.Errorf("Received empty project (%s)", projectID)
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
