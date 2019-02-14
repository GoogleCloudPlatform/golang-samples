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

// Unit tests for services package

package services

import (
	"fmt"
	"testing"
)

func TestCheckMessages(t *testing.T) {
	fmt.Print("Starting unit tests\n")
	messageService := MockMessageService{}
	_, err := CheckMessages(messageService, "user")
	if err != nil {
		t.Errorf("TestCheckMessages: Got an error: %v\n", err)
	}
}

func TestSendUserMessage(t *testing.T) {
	messageService := MockMessageService{}
	message := Message{
		User:   "Unit",
		Friend: "Test",
		Text:   "We mock you!",
		Id:     1}
	err := SendUserMessage(messageService, message)
	if err != nil {
		t.Errorf("TestSendUserMessage: Got an error: %v\n", err)
	}
	messages, err := CheckMessages(messageService, "Test")
	if err != nil {
		t.Errorf("TestSendUserMessage: Got an error on chec: %v\n", err)
	}
	expected := 1
	result := len(messages)
	if result != expected {
		t.Errorf("TestSendUserMessage: Expected: %d, got %d\n", expected, result)
	}
}
