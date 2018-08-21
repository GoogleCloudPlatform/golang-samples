// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
