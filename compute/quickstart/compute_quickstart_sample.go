// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
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
	mongoport = "27017"
	port      = "80"
)

func main() {
	mongohost := os.Getenv("HOST")

	if mongohost == "" {
		log.Fatalf("HOST Environmental variable was not set. Should be the ip of a vm running mongo on port %s", mongoport)
	}

	ctx := context.Background()
	var err error
	tm, err := newTrainerManager(ctx, mongohost, mongoport)
	if err != nil {
		log.Fatalf("error connecting to mongo: %s", err)
	}

	trainers := []trainer{
		{Name: "Ash", Age: 10, City: "Pallet Town"},
		{Name: "Misty", Age: 10, City: "Cerulean City"},
		{Name: "Brock", Age: 15, City: "Pewter City"},
	}

	if err := tm.load(ctx, trainers); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", listHandler(tm))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func listHandler(tm *trainerManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trainers, err := tm.list(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		j, err := json.Marshal(trainers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)
		return
	}
}

type trainer struct {
	Name string
	Age  int
	City string
}

// newTrainerManager spins up a new TrainerManager for interacting with MongoDB.
func newTrainerManager(ctx context.Context, host, port string) (*trainerManager, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", host, port)
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo: %s", err)
	}

	collection := client.Database("test").Collection("trainers")

	return &trainerManager{
		client:     client,
		collection: collection,
	}, nil
}

type trainerManager struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// load pushes a collection of trainers into a mongoDB instance
func (tm *trainerManager) load(ctx context.Context, trainers []trainer) error {
	t := make([]interface{}, len(trainers))
	for i, tdata := range trainers {
		t[i] = tdata
	}

	if _, err := tm.collection.InsertMany(ctx, t); err != nil {
		return fmt.Errorf("error inserting records to mongo: %s", err)
	}

	return nil
}

// list retrieves the total collection of trainers from a mongoDB instance
func (tm *trainerManager) list(ctx context.Context) ([]*trainer, error) {
	var results []*trainer

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

	return results, cur.Err()
}
