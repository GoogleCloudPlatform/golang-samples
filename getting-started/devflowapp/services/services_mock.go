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

// Mock objects for messaging service

package services

import (
	"log"
)

// Mock object that saves the messages in app memory
type MockMessageService map[string][]Message

// Gets messages from app memory
func (service MockMessageService) GetMessages(userTo string) ([]Message, error) {
	log.Printf("MockMicroservice.GetMessages, len: %d\n", len(service))
	messages, ok := service[userTo]
	if ok {
		return messages, nil
	}
	return []Message{}, nil
}

// Saves messages to app memory
func (service MockMessageService) SendMessage(userFrom, userTo,
	text string) error {
	log.Printf("MockMicroservice.SendMessage, Message: %s\n", text)
	message := Message{
		User:   userFrom,
		Friend: userTo,
		Text:   text,
		Id:     len(service),
	}
	messages, ok := service[userTo]
	if ok {
		messages = append(messages, message)
	} else {
		messages = []Message{message}
	}
	service[userTo] = messages
	return nil
}
