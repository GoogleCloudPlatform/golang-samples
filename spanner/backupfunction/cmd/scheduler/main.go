// Copyright 2020 Google LLC
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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	scheduler "cloud.google.com/go/scheduler/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/spanner/backupfunction"
	schedulerpb "google.golang.org/genproto/googleapis/cloud/scheduler/v1"
	"gopkg.in/yaml.v2"
)

const defaultLocation = "us-central1"
const pubsubTopic = "cloud-spanner-scheduled-backups"
const jobPrefix = "spanner-backup"

type Project struct {
	Name      string `yaml:"name"`
	Instances []struct {
		Name      string `yaml:"name"`
		Databases []struct {
			Name     string `yaml:"name"`
			Schedule string `yaml:"schedule"`
			Expire   string `yaml:"expire"`
			Location string `yaml:"location"`
		}
	} `yaml:"instances"`
}

func main() {
	var filename string

	flag.StringVar(&filename, "config", "", "The file path of the config file in yaml format.")
	flag.Parse()

	if filename == "" {
		flag.Usage()
		os.Exit(2)
	}
	content, err := ioutil.ReadFile(filename)

	var project Project

	err = yaml.Unmarshal(content, &project)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	client, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create a scheduler client: %v", err)
	}
	defer client.Close()

	for _, instance := range project.Instances {
		for _, db := range instance.Databases {
			dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", project.Name, instance.Name, db.Name)
			// Get the specified location. If not given, use the default one.
			loc := db.Location
			if loc == "" {
				loc = defaultLocation
			}

			jobName := fmt.Sprintf("%s-%s", jobPrefix, db.Name)
			meta := backupfunction.Meta{Database: dbPath, Expire: db.Expire}
			data, err := json.Marshal(meta)
			if err != nil {
				log.Fatalf("Failed to marshal data: %v", err)
			}

			// Create a new job.
			req := &schedulerpb.CreateJobRequest{
				Parent: fmt.Sprintf("projects/%s/locations/%s", project.Name, loc),
				Job: &schedulerpb.Job{
					Name:        fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project.Name, loc, jobName),
					Description: fmt.Sprintf("A scheduling job for Cloud Spanner database %s", dbPath),
					Target: &schedulerpb.Job_PubsubTarget{
						PubsubTarget: &schedulerpb.PubsubTarget{
							TopicName: fmt.Sprintf("projects/%s/topics/%s", project.Name, pubsubTopic),
							Data:      data,
						},
					},
					Schedule: db.Schedule,
				},
			}
			resp, err := client.CreateJob(ctx, req)
			if err != nil {
				log.Fatalf("Failed to create a cloud scheduler job: %v", err)
			}
			log.Printf("Create a scheduled backup job: %v\n", resp)
		}
	}
}
