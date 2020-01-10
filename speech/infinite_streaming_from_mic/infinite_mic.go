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

// To run this program, SoX could be used to capture mic input:
// sox -d -e signed -b 16 -c 1 -q mic.wav | go run infinite_mic.go

/*
The infinite_streaming_from_mic program reads audio data from microphone
and pipes it to Google Speech API to output the transcript infinitely.
*/
package main

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	color "github.com/fatih/color"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

const (
	sampleRate            = 44100                                      // Audio sample rate hertz.
	inputChannel          = 1                                          // Number of audio input channel.
	outputChannel         = 0                                          // Number of audio output channel.
	bytesPerSample        = 2                                          // Number of bytes each sample consists of.
	bytesPerSecond        = sampleRate * inputChannel * bytesPerSample // Number of bytes each second audio consists of.
	streamTimeLimit       = 290 * time.Second                          // Streaming API Limit( ≈ 5 mins).
	sampleDuration  int64 = 1000                                       // The duration of each sample is 1000ms.
)

func main() {
	// Define data storage buffer.
	framesPerBuffer := make([]byte, bytesPerSecond)
	processingBuffer := make([][]byte, 0)

	// Define time variable.
	streamDeadline := time.Now().Add(streamTimeLimit)
	var processingBufferStart int64 = 0
	restartCounter := 0

	// Define flag variable.
	streamIsNull := true
	exitFlag := false
	inputIsNull := false

	// Define speech stream variable.
	var ctx context.Context
	var client *speech.Client
	var stream speechpb.Speech_StreamingRecognizeClient
	var err error

	// Open established mic recording audio file and delay a while to avoid file non-existing error.
	time.Sleep(1 * time.Second)
	file, err := os.Open("mic.wav")
	if err != nil {
		log.Printf("Couldn't open mic recording file: %v", err)
		return
	}

	// There are three functions in the loop:  (Initialize Stream) ->   [ Read Audio | Receive Responses ].
	for !exitFlag {
		// Exit from loop if exitFlag condition is met.

		// Wait for Goroutines to finish before going to next loop.
		var wg sync.WaitGroup
		wg.Add(1)

		// Start stream.
		if streamIsNull {
			// Initialize speech stream.
			ctx = context.Background()

			client, err = speech.NewClient(ctx)
			if err != nil {
				log.Fatal(err)
			}

			stream, err = client.StreamingRecognize(ctx)
			if err != nil {
				log.Fatal(err)
			}

			// Start timing.
			streamIsNull = false
			startTime := time.Now()
			streamDeadline = startTime.Add(streamTimeLimit)

			restartCounter++
			color.Red("Start New Stream %v!\n", restartCounter)

			// Send the initial configuration message.
			if err := stream.Send(&speechpb.StreamingRecognizeRequest{
				StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
					StreamingConfig: &speechpb.StreamingRecognitionConfig{
						Config: &speechpb.RecognitionConfig{
							Encoding:        speechpb.RecognitionConfig_LINEAR16,
							SampleRateHertz: sampleRate,
							LanguageCode:    "en-US",
						},
					},
				},
			}); err != nil {
				log.Fatal(err)
			}
		}

		// Read contents from audio file and send requests.
		go func() {
			// Make sure that concurrent thread has finished before going to next loop.
			defer wg.Done()

			// Send leftover audio chunks which have not obtained final transcipts from last stream before sending new requests.
			color.Green("Writing %d chunks into the new stream.\n", len(processingBuffer))
			for _, chunk := range processingBuffer {
				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: chunk,
					},
				}); err != nil {
					log.Printf("Could not send audio: %v", err)
				}
			}

			// Read data from audio file and send new requests periodically(1000 ms).
			for {
				// If it has reached 5 mins streaming API time limit, then close the current stream.
				if time.Now().After(streamDeadline) {
					// Nothing else to pipe, close the stream.
					color.Green("Closing stream before it times out")
					if err := stream.CloseSend(); err != nil {
						log.Fatalf("Could not close stream: %v", err)
					}
					streamIsNull = true
					return
				}

				// Read audio data from file.
				if _, err := file.Read(framesPerBuffer); err != nil {
					log.Printf("Could not read any data: %v", err)
					inputIsNull = true
				}
				// Append the new audio data to processingBuffer for future comparison.
				processingBuffer = append(processingBuffer, [][]byte{framesPerBuffer}...)

				// Send new requests.
				if len(framesPerBuffer) > 0 {
					if err := stream.Send(&speechpb.StreamingRecognizeRequest{
						StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
							AudioContent: framesPerBuffer,
						},
					}); err != nil {
						log.Printf("Could not send audio: %v", err)
					}
				}

				// Control sending rate(send 1 s audio data per second) and make sure it is the same as the actual audio play speed.
				time.Sleep(1 * time.Second)
			}
		}()

		// Receive and process responses.
		exitFlag = func() bool {
			for {
				// Receive responses from server.
				resp, err := stream.Recv()

				// Exception Handling:
				// If server has sent all the responses , then it will return EOF(End of File) error to the client.
				if err == io.EOF {
					break
				}

				// If server couldn't send back responses, then it will return detailed error information.
				if err != nil {
					log.Fatalf("Cannot stream results: %v", err)
				}

				// If server couldn't recognize audio data, then it will return detailed error information.
				if err := resp.Error; err != nil {
					log.Fatalf("Could not recognize: %v", err)
				}

				// Process responses.
				for _, result := range resp.Results {
					// If server has sent back final transcript for specific audio chunks, then calculate the end time pointer
					// and remove audio chunks which have obtained final transcript already from processingBuffer.
					if result.GetIsFinal() {
						transcript := result.Alternatives[0].Transcript
						color.Blue("Transcript: %v\n", transcript)

						// If "exit exit exit" or "quit quit quit" command from user microphone input has been recognized, then the program will terminate.
						if strings.Contains(strings.ToLower(transcript), "exit exit exit") || strings.Contains(strings.ToLower(transcript), "quit quit quit") {
							log.Fatalln("Program has stopped running!")
							return true
						}

						// Calculate the correct end time of the specific response which contains final transcript based on restart times.
						resultEndTime := result.GetResultEndTime().Seconds*1000 + int64(math.Round(float64(result.GetResultEndTime().Nanos/1000000))) + int64(int(streamTimeLimit.Seconds())*1000*(restartCounter-1))

						// Remove audio chunks which have obtained final transcript already.
						for len(processingBuffer) > 0 {
							// Every element of processingBuffer is one sample which stores 1 second's audio data, so we need to calculate the first sample end time in processingBuffer.
							sampleEnd := processingBufferStart + sampleDuration

							// Compare first sample end time with the end time of the specific response to determine when to stop removing.
							if sampleEnd > resultEndTime {
								break
							}

							// Move the sample end time pointer.
							processingBufferStart = sampleEnd
							// Remove the first sample(audio chunk) from the head of processingBuffer like Queue(FIFO).
							processingBuffer = processingBuffer[1:]
						}

					}
				}
			}
			// If there is no more data to read from input, then the program will terminate.
			return inputIsNull
		}()
		wg.Wait()
	}
}

// [END speech_transcribe_streaming_mic]
