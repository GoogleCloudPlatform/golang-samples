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

// Sample spanner_arrays is a demonstration program which queries Google's Cloud Spanner
// and returns results containing arrays.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"google.golang.org/api/iterator"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

// Country describes a country and the cities inside it.
type Country struct {
	Name    string
	Colours []spanner.NullString
	Cities  []spanner.NullString
}

func main() {
	ctx := context.Background()

	dsn := flag.String("database", "projects/your-project-id/instances/your-instance-id/databases/your-database-id", "Cloud Spanner database name")
	flag.Parse()

	// Connect to the Spanner Admin API.
	admin, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		log.Fatalf("failed to create database admin client: %v", err)
	}
	defer admin.Close()

	err = createDatabase(ctx, admin, *dsn)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	defer removeDatabase(ctx, admin, *dsn)

	// Connect to database.
	client, err := spanner.NewClient(ctx, *dsn)
	if err != nil {
		log.Fatalf("Failed to create client %v", err)
	}
	defer client.Close()

	err = loadPresets(ctx, client)
	if err != nil {
		log.Fatalf("failed to load preset data: %v", err)
	}

	it := client.Single().Query(ctx, spanner.NewStatement(`
		SELECT a.Name AS Name, ARRAY(
			SELECT b.Name FROM Cities b WHERE a.CountryId = b.CountryId
		) AS Cities, Colours FROM Countries a
	`))
	defer it.Stop()

	for {
		row, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed to read results: %v", err)
		}

		var country Country
		if err = row.ToStruct(&country); err != nil {
			log.Fatalf("failed to read row into Country struct: %v", err)
		}

		var cities []string
		for _, c := range country.Cities {
			cities = append(cities, c.String())
		}

		var colours []string
		for _, c := range country.Colours {
			colours = append(colours, c.String())
		}
		log.Printf("%s (%s): %s", country.Name, strings.Join(colours, ", "), strings.Join(cities, ", "))
	}
}

// loadPresets inserts some demonstration data into the tables.
func loadPresets(ctx context.Context, db *spanner.Client) error {
	mx := []*spanner.Mutation{
		spanner.InsertMap("Countries", map[string]interface{}{
			"CountryId": 49,
			"Name":      "Germany",
			"Colours":   []string{"black", "red", "gold"},
		}),
		spanner.InsertMap("Cities", map[string]interface{}{
			"CountryId":  49,
			"CityId":     100,
			"Name":       "Berlin",
			"Population": 3605000,
		}),
		spanner.InsertMap("Cities", map[string]interface{}{
			"CountryId":  49,
			"CityId":     101,
			"Name":       "Hamburg",
			"Population": 1739117,
		}),
		spanner.InsertMap("Cities", map[string]interface{}{
			"CountryId":  49,
			"CityId":     102,
			"Name":       "Dresden",
			"Population": 486854,
		}),
		spanner.InsertMap("Countries", map[string]interface{}{
			"CountryId": 44,
			"Name":      "United Kingdom",
			"Colours":   []string{"white", "red", "blue"},
		}),
		spanner.InsertMap("Cities", map[string]interface{}{
			"CountryId":  44,
			"CityId":     200,
			"Name":       "London",
			"Population": 8788000,
		}),
		spanner.InsertMap("Cities", map[string]interface{}{
			"CountryId":  44,
			"CityId":     201,
			"Name":       "Liverpool",
			"Population": 465700,
		}),
		spanner.InsertMap("Cities", map[string]interface{}{
			"CountryId":  44,
			"CityId":     202,
			"Name":       "Bristol",
			"Population": 428100,
		}),
		spanner.InsertMap("Cities", map[string]interface{}{
			"CountryId":  44,
			"CityId":     203,
			"Name":       "Newcastle",
			"Population": 304636,
		}),
	}

	_, err := db.Apply(ctx, mx)
	return err
}

// createDatabase uses the Spanner database administration client to create the tables used in this demonstration.
func createDatabase(ctx context.Context, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		log.Fatalf("Invalid database id %s", db)
	}

	var (
		projectID    = matches[1]
		databaseName = matches[2]
	)

	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          projectID,
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseName),
		ExtraStatements: []string{
			`CREATE TABLE Countries (
				CountryId 	INT64 NOT NULL,
				Name   		STRING(1024) NOT NULL,
				Colours     ARRAY<STRING(1024)> NOT NULL
			) PRIMARY KEY (CountryId)`,
			`CREATE TABLE Cities (
				CountryId	INT64 NOT NULL,
				CityId		INT64 NOT NULL,
				Name			STRING(MAX) NOT NULL,
				Population  INT64 NOT NULL
			) PRIMARY KEY (CountryId, CityId),
			INTERLEAVE IN PARENT Countries ON DELETE CASCADE`,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err == nil {
		log.Printf("Created database [%s]", db)
	}
	return err
}

// removeDatabase deletes the database which this demonstration program created.
func removeDatabase(ctx context.Context, adminClient *database.DatabaseAdminClient, db string) {
	if err := adminClient.DropDatabase(ctx, &adminpb.DropDatabaseRequest{Database: db}); err != nil {
		log.Fatalf("Failed to remove database [%s]: %v", db, err)
	} else {
		log.Printf("Removed database [%s]", db)
	}
}
