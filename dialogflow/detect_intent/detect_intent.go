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

package detect

// [START import_libraries]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

// [END import_libraries]

// [START dialogflow_detect_intent_text]
func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}

// [END dialogflow_detect_intent_text]

// [START dialogflow_detect_intent_audio]
func DetectIntentAudio(projectID, sessionID, audioFile, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	// In this example, we hard code the encoding and sample rate for simplicity.
	audioConfig := dialogflowpb.InputAudioConfig{AudioEncoding: dialogflowpb.AudioEncoding_AUDIO_ENCODING_LINEAR_16, SampleRateHertz: 16000, LanguageCode: languageCode}

	queryAudioInput := dialogflowpb.QueryInput_AudioConfig{AudioConfig: &audioConfig}

	audioBytes, err := os.ReadFile(audioFile)
	if err != nil {
		return "", err
	}

	queryInput := dialogflowpb.QueryInput{Input: &queryAudioInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput, InputAudio: audioBytes}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}

// [END dialogflow_detect_intent_audio]

// [START dialogflow_detect_intent_streaming]
func DetectIntentStream(projectID, sessionID, audioFile, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	// In this example, we hard code the encoding and sample rate for simplicity.
	audioConfig := dialogflowpb.InputAudioConfig{AudioEncoding: dialogflowpb.AudioEncoding_AUDIO_ENCODING_LINEAR_16, SampleRateHertz: 16000, LanguageCode: languageCode}

	queryAudioInput := dialogflowpb.QueryInput_AudioConfig{AudioConfig: &audioConfig}

	queryInput := dialogflowpb.QueryInput{Input: &queryAudioInput}

	streamer, err := sessionClient.StreamingDetectIntent(ctx)
	if err != nil {
		return "", err
	}

	f, err := os.Open(audioFile)
	if err != nil {
		return "", err
	}

	defer f.Close()

	go func() {
		audioBytes := make([]byte, 1024)

		request := dialogflowpb.StreamingDetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}
		err = streamer.Send(&request)
		if err != nil {
			log.Fatal(err)
		}

		for {
			_, err := f.Read(audioBytes)
			if err == io.EOF {
				streamer.CloseSend()
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			request = dialogflowpb.StreamingDetectIntentRequest{InputAudio: audioBytes}
			err = streamer.Send(&request)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	var queryResult *dialogflowpb.QueryResult

	for {
		response, err := streamer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		recognitionResult := response.GetRecognitionResult()
		transcript := recognitionResult.GetTranscript()
		log.Printf("Recognition transcript: %s\n", transcript)

		queryResult = response.GetQueryResult()
	}

	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}

// [END dialogflow_detect_intent_streaming]
