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

// Configure database connection

package services

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var messageService MessageService

func getDBConnection() (db *sql.DB, err error) {
	conStr, ok := os.LookupEnv("MYSQL_CONNECTION")
	if ok {
		return sql.Open("mysql", conStr)
	} else {
		dbUser, ok := os.LookupEnv("DB_USER")
		if ok {
			dbPass, _ := os.LookupEnv("DB_PASSWORD")
			conStr = fmt.Sprintf("%s:%s@tcp(localhost:3306)/messagesdb", dbUser,
				dbPass)
			return sql.Open("mysql", conStr)
		} else {
			return db, errors.New("No database connection information provided")
		}
	}
}

// Gets an the MessageService if already instantiated, or creates a new one
func GetMessageService() MessageService {
	if messageService == nil {
		messageService = newMessageService()
	}
	return messageService
}

// Instantiates a MessageService for use in the app
func newMessageService() MessageService {
	log.Printf("newMessageService, enter\n")
	mService, ok := os.LookupEnv("MESSAGE_SERVICE")
	if ok && mService == "mock" {
		return MockMessageService{}
	}
	dbConn, err := getDBConnection()
	if err != nil {
		log.Fatal("service.NewMessageService: error, ", err)
	}
	return SQLMessagingService{dbConn}
}
