// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
