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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// [END import_libraries]

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> -session-id <SESSION ID> -language-code <LANGUAGE CODE> <OPERATION> <INPUTS>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project id\n")
		fmt.Fprintf(os.Stderr, "<SESSION ID> must be a Dialogflow session ID\n")
		fmt.Fprintf(os.Stderr, "<LANGUAGE CODE> must be a language code from https://dialogflow.com/docs/reference/language; defaults to en\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of text, audio, stream\n")
		fmt.Fprintf(os.Stderr, "<INPUTS> can be a series of text inputs if <OPERATION> is text, or a path to an audio file if <OPERATION> is audio or stream\n")
	}

	var projectID, sessionID, languageCode string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&sessionID, "session-id", "", "Dialogflow session ID")
	flag.StringVar(&languageCode, "language-code", "en", "Dialogflow language code from https://dialogflow.com/docs/reference/language; defaults to en")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	operation := args[0]
	inputs := args[1:]

	switch operation {
	case "text":
		fmt.Printf("Responses:\n")
		for _, query := range inputs {
			fmt.Printf("\nInput: %s\n", query)
			response, err := DetectIntentText(projectID, sessionID, query, languageCode)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Output: %s\n", response)
		}
	case "audio":
		if len(inputs) != 1 {
			log.Fatal("audio intent detection expects a single audio file as input")
		}
		audioFile := inputs[0]
		response, err := DetectIntentAudio(projectID, sessionID, audioFile, languageCode)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Response: %s\n", response)
	case "stream":
		if len(inputs) != 1 {
			log.Fatal("audio intent detection expects a single audio file as input")
		}
		audioFile := inputs[0]
		response, err := DetectIntentStream(projectID, sessionID, audioFile, languageCode)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Response: %s\n", response)

	default:
		flag.Usage()
		os.Exit(1)
	}
}

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

	audioBytes, err := ioutil.ReadFile(audioFile)
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
