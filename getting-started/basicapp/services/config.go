// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
