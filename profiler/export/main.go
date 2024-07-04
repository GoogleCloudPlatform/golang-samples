// Copyright 2023 Google LLC
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

// [START profiler_export_profiles]

// Sample export shows how ListProfiles API can be used to download
// existing pprof profiles for a given project from GCP.
package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	cloudprofiler "cloud.google.com/go/cloudprofiler/apiv2"
	pb "cloud.google.com/go/cloudprofiler/apiv2/cloudprofilerpb"
	"google.golang.org/api/iterator"
)

var project = flag.String("project", "", "GCP project ID from which profiles should be fetched")
var pageSize = flag.Int("page_size", 100, "Number of profiles fetched per page. Maximum 1000.")
var pageToken = flag.String("page_token", "", "PageToken from a previous ListProfiles call. If empty, the listing will start from the begnning. Invalid page tokens result in error.")
var maxProfiles = flag.Int("max_profiles", 1000, "Maximum number of profiles to fetch across all pages. If this is <= 0, will fetch all available profiles")

const ProfilesDownloadedSuccessfully = "Read max allowed profiles"

// This function reads profiles for a given project and stores them into locally created files.
// The profile metadata gets stored into a 'metdata.csv' file, while the individual pprof files
// are created per profile.
func downloadProfiles(ctx context.Context, w io.Writer, project, pageToken string, pageSize, maxProfiles int) error {
	client, err := cloudprofiler.NewExportClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	log.Printf("Attempting to fetch %v profiles with a pageSize of %v for %v\n", maxProfiles, pageSize, project)

	// Initial request for the ListProfiles API
	request := &pb.ListProfilesRequest{
		Parent:    fmt.Sprintf("projects/%s", project),
		PageSize:  int32(pageSize),
		PageToken: pageToken,
	}

	// create a folder for storing profiles & metadata
	profilesDirName := fmt.Sprintf("profiles_%v", time.Now().Unix())
	if err := os.Mkdir(profilesDirName, 0750); err != nil {
		log.Fatal(err)
	}
	// create a file for storing profile metadata
	metadata, err := os.Create(fmt.Sprintf("%s/metadata.csv", profilesDirName))
	if err != nil {
		return err
	}
	defer metadata.Close()

	writer := csv.NewWriter(metadata)
	defer writer.Flush()

	writer.Write([]string{"File", "Name", "ProfileType", "Target", "Duration", "Labels"})

	profileCount := 0
	// Keep calling ListProfiles API till all profile pages are fetched or max pages reached
	profilesIterator := client.ListProfiles(ctx, request)
	for {
		// Read individual profile - the client will automatically make API calls to fetch next pages
		profile, err := profilesIterator.Next()

		if err == iterator.Done {
			log.Println("Read all available profiles")
			break
		}
		if err != nil {
			return fmt.Errorf("error reading profile from response: %w", err)
		}
		profileCount++

		filename := fmt.Sprintf("%s/profile%06d.pb.gz", profilesDirName, profileCount)
		err = os.WriteFile(filename, profile.ProfileBytes, 0640)

		if err != nil {
			return fmt.Errorf("unable to write file %s: %w", filename, err)
		}
		fmt.Fprintf(w, "deployment target: %v\n", profile.Deployment.Labels)

		labelBytes, err := json.Marshal(profile.Labels)
		if err != nil {
			return err
		}

		err = writer.Write([]string{filename, profile.Name, profile.Deployment.Target, profile.Duration.String(), string(labelBytes)})
		if err != nil {
			return err
		}

		if maxProfiles > 0 && profileCount >= maxProfiles {
			fmt.Fprintf(w, "result: %v", ProfilesDownloadedSuccessfully)
			break
		}

		if profilesIterator.PageInfo().Remaining() == 0 {
			// This signifies that the client will make a new API call internally
			log.Printf("next page token: %v\n", profilesIterator.PageInfo().Token)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	// validate project ID
	if *project == "" {
		log.Fatalf("No project ID provided, please provide the GCP project ID via '-project' flag")
	}
	var writer bytes.Buffer
	if err := downloadProfiles(context.Background(), &writer, *project, *pageToken, *pageSize, *maxProfiles); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished reading all profiles")
}

// [END profiler_export_profiles]
