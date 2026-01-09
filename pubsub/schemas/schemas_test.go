// Copyright 2021 Google LLC
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

// Package schema is a tool to manage Google Cloud Pub/Sub schemas by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview
package schema

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"
	schema "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

const (
	schemaPrefix     = "test-schema-"
	avroFilePath     = "./resources/us-states.avsc"
	protoFilePath    = "./resources/us-states.proto"
	avroRevFilePath  = "./resources/us-states-plus.avsc"
	protoRevFilePath = "./resources/us-states-plus.proto"

	topicPrefix = "test-topic-"
	subPrefix   = "test-sub-"
)

// once guards cleanup related operations in setup. No need to set up and tear
// down every time, so this speeds things up.
var once sync.Once

func setup(t *testing.T) (*pubsub.Client, *schema.SchemaClient) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	schemaClient, err := schema.NewSchemaClient(ctx)
	if err != nil {
		t.Fatalf("failed to create schema client: %v", err)
	}

	// Cleanup schema resources from the previous tests.
	once.Do(func() {
		scs, err := listSchemas(io.Discard, tc.ProjectID)
		if err != nil {
			fmt.Printf("failed to list schemas: %v", err)
		}
		for _, sc := range scs {
			schemaName := strings.Split(sc.Name, "/")
			schemaID := schemaName[len(schemaName)-1]
			if strings.HasPrefix(schemaID, schemaPrefix) {
				deleteSchema(io.Discard, tc.ProjectID, schemaID)
			}
		}
	})

	return client, schemaClient
}

func TestSchemas_Admin(t *testing.T) {
	_, sc := setup(t)
	tc := testutil.SystemTest(t)

	avroSchemaID := schemaPrefix + "avro-" + uuid.NewString()
	t.Run("createAvroSchema", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			if err := createAvroSchema(buf, tc.ProjectID, avroSchemaID, avroFilePath); err != nil {
				r.Errorf("createAvroSchema err: %v", err)
			}
			got := buf.String()
			want := "Schema created"
			if !strings.Contains(got, want) {
				r.Errorf("createAvroSchema() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	t.Run("commitAvroSchema", func(t *testing.T) {
		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			if err := commitAvroSchema(buf, tc.ProjectID, avroSchemaID, avroRevFilePath); err != nil {
				r.Errorf("commitAvroSchema err: %v\n", err)
			}
			got := buf.String()
			want := "Committed a schema using an Avro schema"
			if !strings.Contains(got, want) {
				r.Errorf("commitAvroSchema() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	protoSchemaID := schemaPrefix + "proto-" + uuid.NewString()
	protoSchemaName := fmt.Sprintf("projects/%s/schemas/%s", tc.ProjectID, protoSchemaID)
	var protoSchema *pubsubpb.Schema
	t.Run("createProtoSchema", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			if err := createProtoSchema(buf, tc.ProjectID, protoSchemaID, protoFilePath); err != nil {
				r.Errorf("create err: %v", err)
			}
			got := buf.String()
			want := "Schema created"
			if !strings.Contains(got, want) {
				r.Errorf("createProtoSchema() got: %q\nwant: %q\n", got, want)
			}

			ctx := context.Background()
			var err error
			req := &pubsubpb.GetSchemaRequest{
				Name: protoSchemaName,
				View: pubsubpb.SchemaView_FULL,
			}
			protoSchema, err = sc.GetSchema(ctx, req)
			if err != nil {
				r.Errorf("failed to get proto schema: %v\n", err)
			}
		})
	})

	t.Run("commitProtoSchema", func(t *testing.T) {
		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			if err := commitProtoSchema(buf, tc.ProjectID, protoSchemaID, protoRevFilePath); err != nil {
				r.Errorf("commitProtoSchema err: %v\n", err)
			}
			got := buf.String()
			want := "Committed a schema using a protobuf schema"
			if !strings.Contains(got, want) {
				r.Errorf("commitAvroSchema() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	t.Run("rollbackSchema", func(t *testing.T) {
		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			if err := rollbackSchema(buf, tc.ProjectID, protoSchemaID, protoSchema.RevisionId); err != nil {
				r.Errorf("rollbackSchema err: %v\n", err)
			}
			got := buf.String()
			want := "Rolled back schema"
			if !strings.Contains(got, want) {
				r.Errorf("rollbackSchema() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	t.Run("getSchema", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := getSchema(buf, tc.ProjectID, avroSchemaID)
			if err != nil {
				r.Errorf("getSchema err: %v", err)
			}
			got := buf.String()
			want := "Got schema"
			if !strings.Contains(got, want) {
				r.Errorf("getSchema() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	t.Run("getSchemaRevision", func(t *testing.T) {
		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			schemaRev := fmt.Sprintf("%s@%s", protoSchemaID, protoSchema.RevisionId)
			err := getSchemaRevision(buf, tc.ProjectID, schemaRev)
			if err != nil {
				r.Errorf("getSchemaRevision err: %v", err)
			}
			got := buf.String()
			want := "Got schema revision"
			if !strings.Contains(got, want) {
				r.Errorf("getSchemaRevision() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	t.Run("listSchemas", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			schemas, err := listSchemas(buf, tc.ProjectID)
			if err != nil {
				r.Errorf("failed to list schemas: %v", err)
			}
			if len(schemas) != 2 {
				r.Errorf("expected 2 schemas, got %d", len(schemas))
			}
		})
	})

	t.Run("listSchemaRevisions", func(t *testing.T) {
		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			_, err := listSchemaRevisions(buf, tc.ProjectID, protoSchemaID)
			if err != nil {
				r.Errorf("failed to list schemas: %v", err)
			}
			got := buf.String()
			want := "Got schema revision"
			if !strings.Contains(got, want) {
				r.Errorf("listSchemaRevisions() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	topicID := topicPrefix + uuid.NewString()
	t.Run("createTopicWithSchemaRevisions", func(t *testing.T) {
		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := createTopicWithSchemaRevisions(buf, tc.ProjectID, topicID, protoSchemaID, protoSchema.RevisionId, protoSchema.RevisionId)
			if err != nil {
				r.Errorf("createTopicWithSchemaRevisions err: %v", err)
			}
			got := buf.String()
			want := "Created topic with schema revision"
			if !strings.Contains(got, want) {
				r.Errorf("createTopicWithSchemaRevisions() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	t.Run("deleteSchemaRevision", func(t *testing.T) {
		testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			if err := deleteSchemaRevision(buf, tc.ProjectID, protoSchemaID, protoSchema.RevisionId); err != nil {
				r.Errorf("deleteSchemaRevision err: %v", err)
			}
			got := buf.String()
			want := "Deleted a schema revision"
			if !strings.Contains(got, want) {
				r.Errorf("deleteSchemaRevision() got: %q\nwant: %q\n", got, want)
			}
		})
	})

	t.Run("deleteSchema", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			if err := deleteSchema(buf, tc.ProjectID, avroSchemaID); err != nil {
				r.Errorf("deleteSchema err: %v", err)
			}
			if err := deleteSchema(buf, tc.ProjectID, protoSchemaID); err != nil {
				r.Errorf("deleteSchema err: %v", err)
			}
		})
	})
}

func TestSchemas_AvroSchemaAll(t *testing.T) {
	client, _ := setup(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	topicID := topicPrefix + uuid.NewString()
	avroSchemaID := schemaPrefix + "avro-" + uuid.NewString()
	_, err := defaultSchemaConfig(tc.ProjectID, avroSchemaID, avroFilePath, pubsubpb.Schema_AVRO)
	if err != nil {
		t.Fatalf("defaultSchemaConfig err: %v", err)
	}
	subID := subPrefix + uuid.NewString()

	t.Run("createTopicWithSchema", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := createAvroSchema(io.Discard, tc.ProjectID, avroSchemaID, avroFilePath); err != nil {
				r.Errorf("createAvroSchema err: %v", err)
			}

			buf := new(bytes.Buffer)
			err := createTopicWithSchema(buf, tc.ProjectID, topicID, avroSchemaID, pubsubpb.Encoding_JSON)
			if err != nil {
				r.Errorf("createTopicWithSchema: %v", err)
			}
			got := buf.String()
			want := "Topic with schema created"
			if !strings.Contains(got, want) {
				r.Errorf("createTopicWithSchema mismatch\ngot: %v\nwant: %v\n", got, want)
			}

			sub := &pubsubpb.Subscription{
				Name:  fmt.Sprintf("projects/%s/subscriptions/%s", tc.ProjectID, subID),
				Topic: fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, topicID),
			}
			if _, err = client.SubscriptionAdminClient.CreateSubscription(ctx, sub); err != nil {
				r.Errorf("client.CreateSubscription err: %v", err)
			}
		})
	})

	t.Run("publishAvroRecords", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := publishAvroRecords(buf, tc.ProjectID, topicID, avroFilePath)
			if err != nil {
				r.Errorf("publishAvroRecords: %v", err)
			}
			got := buf.String()
			want := "Published avro record: {\"name\":\"Alaska\",\"post_abbr\":\"AK\"}\n"
			if diff := cmp.Diff(want, got); diff != "" {
				r.Errorf("publishAvroRecords() mismatch: -want, +got:\n%s", diff)
			}
		})
	})

	t.Run("subscribeWithAvroRecords", func(t *testing.T) {
		testutil.Retry(t, 3, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := subscribeWithAvroSchema(buf, tc.ProjectID, subID, avroFilePath)
			if err != nil {
				r.Errorf("subscribeWithAvroSchema: %v", err)
			}
			got := buf.String()
			want := " is abbreviated as "
			if !strings.Contains(got, want) {
				r.Errorf("subscribeWithAvroSchema mismatch\ngot: %v\nwant: %v\n", got, want)
			}
		})
	})

	t.Run("subscribeWithAvroSchemaRevisions", func(t *testing.T) {
		testutil.Retry(t, 3, time.Second, func(r *testutil.R) {
			err := publishAvroRecords(io.Discard, tc.ProjectID, topicID, avroFilePath)
			if err != nil {
				r.Errorf("publishAvroRecords: %v", err)
			}
			buf := new(bytes.Buffer)
			err = subscribeWithAvroSchemaRevisions(buf, tc.ProjectID, subID)
			if err != nil {
				r.Errorf("subscribeWithAvroSchemaRevisions: %v", err)
			}
			got := buf.String()
			want := " is abbreviated as "
			if !strings.Contains(got, want) {
				r.Errorf("subscribeWithAvroSchemaRevisions mismatch\ngot: %v\nwant: %v\n", got, want)
			}
		})
	})

	deleteSchema(io.Discard, tc.ProjectID, avroSchemaID)
	dsr := &pubsubpb.DeleteSubscriptionRequest{
		Subscription: fmt.Sprintf("projects/%s/subscriptions/%s", tc.ProjectID, subID),
	}
	client.SubscriptionAdminClient.DeleteSubscription(ctx, dsr)
	dtr := &pubsubpb.DeleteTopicRequest{
		Topic: fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, topicID),
	}
	client.TopicAdminClient.DeleteTopic(ctx, dtr)
}

func TestSchemas_ProtoSchemaAll(t *testing.T) {
	client, _ := setup(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	topicID := topicPrefix + uuid.NewString()
	protoSchemaID := schemaPrefix + "proto-" + uuid.NewString()
	_, err := defaultSchemaConfig(tc.ProjectID, protoSchemaID, avroFilePath, pubsubpb.Schema_AVRO)
	if err != nil {
		t.Fatalf("defaultSchemaConfig err: %v", err)
	}
	subID := subPrefix + uuid.NewString()

	t.Run("createResources", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := createProtoSchema(io.Discard, tc.ProjectID, protoSchemaID, protoFilePath); err != nil {
				r.Errorf("createProtoSchema err: %v", err)
			}

			buf := new(bytes.Buffer)
			err := createTopicWithSchema(buf, tc.ProjectID, topicID, protoSchemaID, pubsubpb.Encoding_JSON)
			if err != nil {
				r.Errorf("createTopicWithSchema: %v", err)
			}
			got := buf.String()
			want := "Topic with schema created"
			if !strings.Contains(got, want) {
				r.Errorf("createTopicWithSchema mismatch\ngot: %v\nwant: %v\n", got, want)
			}

			sub := &pubsubpb.Subscription{
				Name:  fmt.Sprintf("projects/%s/subscriptions/%s", tc.ProjectID, subID),
				Topic: fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, topicID),
			}
			if _, err = client.SubscriptionAdminClient.CreateSubscription(ctx, sub); err != nil {
				r.Errorf("failed to create subscription: %v", err)
			}
		})
	})

	t.Run("publishProtoMessages", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := publishProtoMessages(buf, tc.ProjectID, topicID)
			if err != nil {
				r.Errorf("publishProtoMessages: %v", err)
			}
			got := buf.String()
			want := "Published proto message"
			if !strings.Contains(got, want) {
				r.Errorf("createTopicWithSchema mismatch\ngot: %v\nwant: %v\n", got, want)
			}
		})
	})

	t.Run("subscribeProtoMessages", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			buf := new(bytes.Buffer)
			err := subscribeWithProtoSchema(buf, tc.ProjectID, subID)
			if err != nil {
				r.Errorf("subscribeWithProtoSchema: %v", err)
			}
			got := buf.String()
			want := " is abbreviated as "
			if !strings.Contains(got, want) {
				r.Errorf("subscribeWithProtoSchema mismatch\ngot: %v\nwant: %v\n", got, want)
			}
		})
	})

	deleteSchema(io.Discard, tc.ProjectID, protoSchemaID)
	dsr := &pubsubpb.DeleteSubscriptionRequest{
		Subscription: fmt.Sprintf("projects/%s/subscriptions/%s", tc.ProjectID, subID),
	}
	client.SubscriptionAdminClient.DeleteSubscription(ctx, dsr)
	dtr := &pubsubpb.DeleteTopicRequest{
		Topic: fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, topicID),
	}
	client.TopicAdminClient.DeleteTopic(ctx, dtr)
}

func TestSchemas_UpdateTopicSchema(t *testing.T) {
	pubsubClient, schemaClient := setup(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	topicID := topicPrefix + uuid.NewString()

	protoSchemaID := schemaPrefix + "proto-" + uuid.NewString()

	protoSource, err := os.ReadFile(protoFilePath)
	if err != nil {
		t.Fatalf("error reading from file: %s", protoFilePath)
	}
	csr := &pubsubpb.CreateSchemaRequest{
		Parent:   fmt.Sprintf("projects/%s", tc.ProjectID),
		SchemaId: protoSchemaID,
		Schema: &pubsubpb.Schema{
			Type:       pubsubpb.Schema_PROTOCOL_BUFFER,
			Definition: string(protoSource),
		},
	}
	schema, err := schemaClient.CreateSchema(ctx, csr)
	if err != nil {
		t.Fatalf("createProtoSchema err: %v", err)
	}

	pubsubClient.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, topicID),
		SchemaSettings: &pubsubpb.SchemaSettings{
			Schema:   schema.GetName(),
			Encoding: pubsubpb.Encoding_BINARY,
		},
	})

	buf := new(bytes.Buffer)
	if err := updateTopicSchema(buf, tc.ProjectID, topicID, schema.RevisionId, schema.RevisionId); err != nil {
		t.Fatalf("updateTopicSchema err : %v", err)
	}
}

func defaultSchemaConfig(projectID, schemaID, schemaFile string, schemaType pubsubpb.Schema_Type) (*pubsubpb.Schema, error) {
	schemaSource, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}
	s := &pubsubpb.Schema{
		Name:       fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
		Type:       schemaType,
		Definition: string(schemaSource),
	}
	return s, nil
}
