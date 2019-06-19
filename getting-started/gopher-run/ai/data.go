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

//Package ai moves the Firestore play data into a csv in a cloud storage bucket
package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	ml "google.golang.org/api/ml/v1"
)

func main() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()
	f, err := os.Create("data/playdata.csv")
	if err != nil {
		// log.Fatalf("os.Create: %v", err)
	}
	defer f.Close()

	mlService, err := ml.NewService(ctx)
	if err != nil {
		log.Printf("ml.NewService: %v", err)
	}
	jobsService := mlService.Projects.Jobs
	job := &ml.GoogleCloudMlV1__Job{
		JobId: "automatic",
		PredictionInput: &ml.GoogleCloudMlV1__PredictionInput{
			DataFormat: "CSV",
			InputPaths: []string{},
		},
		TrainingInput: &ml.GoogleCloudMlV1__TrainingInput{
			PackageUris:  []string{},
			PythonModule: "",
			Region:       "us-west1",
			ScaleTier:    "STANDARD_1",
		}}
	resp, err := jobsService.Create("projects/"+projectID, job).Do()
	if err != nil {
		log.Fatalf("Create: %v", err)
	}
	fmt.Println(resp)
}
