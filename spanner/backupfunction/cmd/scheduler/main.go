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
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
)

const defaultLocation = "us-central1"
const pubsubTopic = "cloud-spanner-scheduled-backups"
const jobPrefix = "spanner-backup"

// Project contains the information of a GCP project.
type Project struct {
	Name      string     `yaml:"name"`
	Instances []Instance `yaml:"instances"`
}

// Instance contains the information of an instance.
type Instance struct {
	Name      string     `yaml:"name"`
	Databases []Database `yaml:"databases"`
}

// Database contains the backup schedule configuration for a database.
type Database struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Expire   string `yaml:"expire"`
	Location string `yaml:"location"`
	TimeZone string `yaml:"time_zone"`
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
		log.Fatalf("Failed to parse the config file: %v", err)
	}

	ctx := context.Background()
	client, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create a scheduler client: %v", err)
	}
	defer client.Close()

	topicPath := fmt.Sprintf("projects/%s/topics/%s", project.Name, pubsubTopic)

	for _, instance := range project.Instances {
		for _, db := range instance.Databases {
			dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", project.Name, instance.Name, db.Name)
			// Get the specified location. If not given, use the default one.
			loc := db.Location
			if loc == "" {
				loc = defaultLocation
			}
			locPath := fmt.Sprintf("projects/%s/locations/%s", project.Name, loc)
			jobID := fmt.Sprintf("%s-%s", jobPrefix, db.Name)
			jobName := fmt.Sprintf("%s/jobs/%s", locPath, jobID)

			err = updateJob(ctx, client, jobName, locPath, dbPath, topicPath, db)
			if err != nil {
				if errCode(err) == codes.NotFound {
					// Create a new job if the job does not exist.
					createJob(ctx, client, jobName, locPath, dbPath, topicPath, db)
				} else {
					log.Printf("Failed to update a job: %v\n", err)
				}
			}
		}
	}
}

// errCode extracts the canonical error code from a Go error.
func errCode(err error) codes.Code {
	s, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}
	return s.Code()
}

func updateJob(ctx context.Context, client *scheduler.CloudSchedulerClient, jobName, locPath, dbPath, topicPath string, db Database) error {
	meta := backupfunction.Meta{Database: dbPath, Expire: db.Expire}
	data, err := json.Marshal(meta)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	// Update a job.
	req := &schedulerpb.UpdateJobRequest{
		Job: &schedulerpb.Job{
			Name: jobName,
			Target: &schedulerpb.Job_PubsubTarget{
				PubsubTarget: &schedulerpb.PubsubTarget{
					TopicName: topicPath,
					Data:      data,
				},
			},
			Schedule: db.Schedule,
			TimeZone: db.TimeZone,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"schedule", "pubsub_target.data", "time_zone"},
		},
	}
	_, err = client.UpdateJob(ctx, req)
	if err == nil {
		log.Printf("Update the job %v.", jobName)
	}
	return err
}

func createJob(ctx context.Context, client *scheduler.CloudSchedulerClient, jobName, locPath, dbPath, topicPath string, db Database) {
	meta := backupfunction.Meta{Database: dbPath, Expire: db.Expire}
	data, err := json.Marshal(meta)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	// Create a new job.
	req := &schedulerpb.CreateJobRequest{
		Parent: locPath,
		Job: &schedulerpb.Job{
			Name:        jobName,
			Description: fmt.Sprintf("A scheduling job for Cloud Spanner database %s", dbPath),
			Target: &schedulerpb.Job_PubsubTarget{
				PubsubTarget: &schedulerpb.PubsubTarget{
					TopicName: topicPath,
					Data:      data,
				},
			},
			Schedule: db.Schedule,
			TimeZone: db.TimeZone,
		},
	}
	resp, err := client.CreateJob(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create a cloud scheduler job: %v", err)
	}
	log.Printf("Create a scheduled backup job: %v\n", resp)
}
