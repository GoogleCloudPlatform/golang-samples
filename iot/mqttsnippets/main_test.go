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

package mqttsnippets

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

var (
	projectID  string
	topicID    string
	topicName  string // topicName is the full path to the topic (e.g. project/{project}/topics/{topic}).
	registryID string
)

func TestMain(m *testing.M) {
	setup(m)
	s := m.Run()
	shutdown()
	os.Exit(s)
}

func setup(m *testing.M) {
	ctx := context.Background()
	tc, ok := testutil.ContextMain(m)

	// Retrieve project ID.
	if !ok {
		fmt.Fprintln(os.Stderr, "Project is not set up properly for system tests. Make sure GOLANG_SAMPLES_PROJECT_ID is set")
		return
	}
	projectID = tc.ProjectID

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Could not create pubsub Client:\n%v\n", err)
		return
	}

	topicID = "golang-iot-test-topic"
	topic := client.Topic(topicID)

	if ok, err = topic.Exists(ctx); err != nil {
		fmt.Printf("Error checking if topic exists:\n%v\n", err)
		return
	}
	if ok {
		if err := topic.Delete(ctx); err != nil {
			fmt.Printf("Could not cleanup existing topic:\n%v\n", err)
		}
	}

	newTopic, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		fmt.Printf("Could not create topic: %v", err)
		return
	}

	topicName = newTopic.String()
	fmt.Printf("Topic created: %v\n", topicName)

	registryID = "golang-iot-test-registry"

	// Delete the current registry if it exists.
	if _, err := deleteRegistry(os.Stdout, projectID, region, registryID); err != nil {
		if !strings.Contains(err.Error(), "Error 404") {
			fmt.Printf("Could not delete registry: %v\n", err)
			return
		}
	}

	if _, err = createRegistry(os.Stdout, projectID, region, registryID, topicName); err != nil {
		fmt.Printf("Could not create registry: %v\n", err)
		return
	}
}

func shutdown() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Could not create pubsub Client:\n%v\n", err)
		return
	}

	topic := client.Topic(topicID)
	if err := topic.Delete(ctx); err != nil {
		fmt.Printf("Could not delete topic: %v\n", err)
		return
	}
	fmt.Printf("Deleted topic: %v\n", topic)

	if _, err := deleteRegistry(os.Stdout, projectID, region, registryID); err != nil {
		fmt.Printf("Could not delete registry: %v\n", err)
	}
}

// deleteRegistry deletes a device registry if it is empty.
func deleteRegistry(w io.Writer, projectID string, region string, registryID string) (*cloudiot.Empty, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	client, err := cloudiot.New(httpClient)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Delete(name).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Deleted registry: %s\n", registryID)

	return response, nil
}
