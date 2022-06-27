// Copyright 2022 Google LLC
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

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoport = "80"
	port      = "80"
)

var tm mongoManager

func main() {
	var err error
	mongohost := os.Getenv("HOST")

	tm, err = newTrainerManager(mongohost, mongoport)
	if err != nil {
		log.Fatal(fmt.Printf("error connecting to mongo: %s", err))
	}

	trainers := []trainer{
		{"Ash", 10, "Pallet Town"},
		{"Misty", 10, "Cerulean City"},
		{"Brock", 15, "Pewter City"},
	}

	if err := tm.load(trainers); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", listHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	trainers, err := tm.list()
	if err != nil {
		httpResponse(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}

	j, err := json.Marshal(trainers)
	if err != nil {
		httpResponse(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}
	httpResponse(w, http.StatusOK, j)
	return
}

func httpResponse(w http.ResponseWriter, status int, msg []byte) {
	w.WriteHeader(status)
	w.Write(msg)
	return
}

type trainer struct {
	Name string
	Age  int
	City string
}

// newTrainerManager spins up a new TrainerManager for interacting with MongoDB.
func newTrainerManager(host, port string) (*trainerManager, error) {
	var err error
	tm := trainerManager{}
	ctx := context.Background()

	uri := fmt.Sprintf("mongodb://%s:%s", host, port)
	clientOptions := options.Client().ApplyURI(uri)

	tm.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo: %s", err)
	}

	tm.collection = tm.client.Database("test").Collection("trainers")

	return &tm, nil
}

type mongoManager interface {
	load([]trainer) error
	list() ([]*trainer, error)
	setCollection(*mongo.Collection)
}

type trainerManager struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (tm *trainerManager) setCollection(c *mongo.Collection) {
	tm.collection = c
}

// load pushes a collection of trainers into a mongoDB instance
func (tm *trainerManager) load(trainers []trainer) error {
	t := make([]interface{}, len(trainers))
	for i, tdata := range trainers {
		t[i] = tdata
	}

	ctx := context.Background()
	if _, err := tm.collection.InsertMany(ctx, t); err != nil {
		return fmt.Errorf("error inserting records to mongo: %s", err)
	}

	return nil
}

// list retrieves the total collection of trainers from a mongoDB instance
func (tm *trainerManager) list() ([]*trainer, error) {
	var results []*trainer
	ctx := context.Background()

	cur, err := tm.collection.Find(ctx, bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem trainer
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
