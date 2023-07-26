// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	memorystore "cloud.google.com/go/redis/apiv1"
	redispb "cloud.google.com/go/redis/apiv1/redispb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

var instanceID string

func TestMain(m *testing.M) {

	// Set up
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatal("test project not set up properly")
		return
	}
	parent := fmt.Sprintf("projects/%s/locations/us-central1/", tc.ProjectID)

	ctx := context.Background()
	id := uuid.New()
	instanceID = fmt.Sprintf("test-instance-%s", id.String())

	adminClient, err := memorystore.NewCloudRedisClient(ctx)
	if err != nil {
		log.Fatal("can't instantiate MemoryStore Redis admin client")
		return
	}

	req := &redispb.CreateInstanceRequest{
		Parent:     parent,
		InstanceId: instanceID,
		Instance: &redispb.Instance{
			Name: fmt.Sprintf("%s/%s", parent, instanceID),
		},
	}

	op, err := adminClient.CreateInstance(ctx, req)
	if err != nil {
		return
	}
	op.Wait(ctx)

	m.Run()

	// Teardown
	_, err = adminClient.DeleteInstance(ctx, &redispb.DeleteInstanceRequest{
		Name: fmt.Sprintf("%s/%s", parent, instanceID),
	})

	if err != nil {
		log.Fatalf("couldn't delete Redis instance: %s", err)
	}
}

func TestConnectToDatabase(t *testing.T) {

	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	err := ConnectToDatabase(&buf, tc.ProjectID, "us-central1", instanceID)
	if err != nil {
		t.Fatal(err)
	}

	want := "Response"
	got := buf.String()

	if !strings.Contains(want, got) {
		t.Errorf("wanted: %s; got: %s", want, got)
	}
}
