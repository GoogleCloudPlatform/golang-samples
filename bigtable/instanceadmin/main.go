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

// [START bigtable_instanceadmin]

// Instance admin is a sample program demonstrating use of the Cloud Bigtable client
// library to perform basic instance admin operations.
package main

import (
	"context"
	"flag"
	"log"

	"cloud.google.com/go/bigtable"
)

func main() {
	projectID := "my-project-id" // The Google Cloud Platform project ID
	instanceID := "my-bigtable-instance"
	clusterID := "my-cluster"
	zone := "us-central1-b"

	// Allow overriding with flags
	flag.StringVar(&projectID, "project", projectID, "The Google Cloud Platform project ID.")
	flag.StringVar(&instanceID, "instance", instanceID, "The Google Cloud Bigtable instance ID to create.")
	flag.StringVar(&zone, "zone", zone, "The zone for the initial cluster.")
	flag.Parse()

	ctx := context.Background()

	instanceAdminClient, err := bigtable.NewInstanceAdminClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create instance admin client: %v", err)
	}
	defer instanceAdminClient.Close()

	// These are placeholders. You must create these in your GCP Organization/Project first.
	// See [Creating and managing tags](https://docs.cloud.google.com/bigtable/docs/tags) for
	// more information on creating a tag.
	// tagKey := "tagKeys/12345"
	// tagValue := "tagValues/56789"

	instanceConf := &bigtable.InstanceConf{
		InstanceId:   instanceID,
		DisplayName:  instanceID,
		ClusterId:    clusterID,
		NumNodes:     1,
		InstanceType: bigtable.PRODUCTION,
		StorageType:  bigtable.SSD,
		Zone:         zone,
		// The 'Tags' field is optional. Uncomment the lines below (and the tagKey, tagValue lines above)
		// only if you want to create an instance with tags.
		// Ensure the tagKey and tagValue exist in your GCP project.
		// Tags: map[string]string{tagKey: tagValue},
	}
	log.Printf("Creating instance %s with cluster %s in %s...", instanceID, clusterID, zone)
	err = instanceAdminClient.CreateInstance(ctx, instanceConf)

	if err != nil {
		log.Fatalf("Could not start create instance operation: %v", err)
	}

	log.Printf("Instance %s created successfully.\n", instanceID)
}
