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

// Service interface definition and implementation

package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Interface for sending messages
type MessageService interface {

	// Gets the messages that have been sent to a user
	GetMessages(userTo string) ([]Message, error)

	// Send a message to a user
	SendMessage(userFrom, userTo, formattedMessage string) error
}

// Encapsulates a message from a User to her or his Friend with message Text
type Message struct {
	User, Friend, Text string
	Id                 int
}

// An implemementation of MessageService using a SQL database
type SQLMessagingService struct{ DBConn *sql.DB }

// Gets messages from a SQL database
func (service SQLMessagingService) GetMessages(userTo string) ([]Message,
	error) {
	log.Printf("SQLMessagingService.GetMessages, userTo: %s\n", userTo)
	messages := []Message{}
	rows, err := service.DBConn.Query(
		"SELECT user_from, text, id FROM messages WHERE user_to = ?",
		userTo)
	if err != nil {
		log.Printf("SQLMessagingService.GetMessages, Error in query: %v\n", err)
		return nil, errors.New("Due to an error, we could not get your messages.")
	}
	defer rows.Close()
	for rows.Next() {
		message := Message{}
		if err := rows.Scan(&message.User, &message.Text, &message.Id); err != nil {
			log.Printf("SQLMessagingService.GetMessages, Error in scan: %v\n", err)
			return nil, errors.New("Due to an error, we could not get all your " +
				"messages.")
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// Saves a message to the SQL database
func (service SQLMessagingService) SendMessage(userFrom, userTo,
	text string) error {
	log.Printf("SQLMessagingService.SendMessage, Message: %s\n", text)
	result, err := service.DBConn.Exec(
		"INSERT INTO messages (user_from, user_to, text) VALUES (?, ?, ?)",
		userFrom, userTo, text)
	if err != nil {
		log.Printf("SQLMessagingService.SendMessage, Error: %v\n", err)
		return errors.New("Due to an error, we could not send your message")
	} else {
		rows, _ := result.RowsAffected()
		id, _ := result.LastInsertId()
		log.Printf("SQLMessagingService.SendMessage, Rows affected: %d, id: %d\n",
			rows, id)
	}
	return nil
}

// Formats a user message
func FormatMessage(user, friend, message string) string {
	return fmt.Sprintf("Hi %s! %s! From %s!", friend, message, user)
}

// Checks user messages, with the given MessageService and user id
func CheckMessages(messageService MessageService,
	userTo string) ([]Message, error) {
	log.Printf("CheckMessages, Message: %s\n", userTo)
	return messageService.GetMessages(userTo)
}

// Formats and sends a user message
func SendUserMessage(messageService MessageService, message Message) error {
	text := FormatMessage(message.User, message.Friend, message.Text)
	error := messageService.SendMessage(message.Friend, message.Friend, text)
	return error
}
