// Copyright 2025 Google LLC
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
	"io/ioutil"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

const (
	schemaPrefix     = "test-schema-"
	avroFilePath     = "../../schemas/resources/us-states.avsc"
	protoFilePath    = "../../schemas/resources/us-states.proto"
	avroRevFilePath  = "../../schemas/resources/us-states-plus.avsc"
	protoRevFilePath = "../../schemas/resources/us-states-plus.proto"

	topicPrefix = "test-topic-"
	subPrefix   = "test-sub-"
)

// once guards cleanup related operations in setup. No need to set up and tear
// down every time, so this speeds things up.
var once sync.Once

func setup(t *testing.T) (*pubsub.Client, *pubsub.SchemaClient) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	schemaClient, err := pubsub.NewSchemaClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create schema client: %v", err)
	}

	// Cleanup resources from the previous tests.
	// This includes schemas, topics, and subscriptions.
	once.Do(func() {
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			scs, err := listSchemas(ioutil.Discard, tc.ProjectID)
			if err != nil {
				fmt.Printf("failed to list schemas: %v", err)
			}
			for _, sc := range scs {
				schemaName := strings.Split(sc.Name, "/")
				schemaID := schemaName[len(schemaName)-1]
				if strings.HasPrefix(schemaID, schemaPrefix) {
					deleteSchema(ioutil.Discard, tc.ProjectID, schemaID)
				}
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			topicIter := client.Topics(ctx)
			for {
				topic, err := topicIter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Printf("topicIter.Next got err: %v", err)
				}
				if strings.HasPrefix(topic.ID(), topicPrefix) {
					if err := topic.Delete(ctx); err != nil {
						fmt.Printf("topic.Delete got err: %v", err)
					}
				}
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			subIter := client.Subscriptions(ctx)
			for {
				sub, err := subIter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Printf("subIter.Next got err: %v", err)
				}
				if strings.HasPrefix(sub.ID(), subPrefix) {
					if err := sub.Delete(ctx); err != nil {
						fmt.Printf("sub.Delete got err: %v", err)
					}
				}
			}
			wg.Done()
		}()
		wg.Wait()
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
	var protoSchema *pubsub.SchemaConfig
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
			protoSchema, err = sc.Schema(ctx, protoSchemaID, pubsub.SchemaViewFull)
			if err != nil {
				r.Errorf("failed to get schema: %v\n", err)
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
			if err := rollbackSchema(buf, tc.ProjectID, protoSchemaID, protoSchema.RevisionID); err != nil {
				r.Errorf("rollbackSchema err: %v\n", err)
			}
			got := buf.String()
			want := "Rolled back a schema"
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
			schemaRev := fmt.Sprintf("%s@%s", protoSchemaID, protoSchema.RevisionID)
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
			// Account for more schemas being created because of retries.
			if len(schemas) < 2 {
				r.Errorf("expected at least 2 schemas, got %d", len(schemas))
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
			err := createTopicWithSchemaRevisions(buf, tc.ProjectID, topicID, protoSchemaID, protoSchema.RevisionID, protoSchema.RevisionID, pubsub.EncodingBinary)
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
			if err := deleteSchemaRevision(buf, tc.ProjectID, protoSchemaID, protoSchema.RevisionID); err != nil {
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
	_, err := defaultSchemaConfig(tc.ProjectID, avroSchemaID, avroFilePath, pubsub.SchemaAvro)
	if err != nil {
		t.Fatalf("defaultSchemaConfig err: %v", err)
	}
	subID := subPrefix + uuid.NewString()

	t.Run("createTopicWithSchema", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := createAvroSchema(ioutil.Discard, tc.ProjectID, avroSchemaID, avroFilePath); err != nil {
				r.Errorf("createAvroSchema err: %v", err)
			}

			buf := new(bytes.Buffer)
			err := createTopicWithSchema(buf, tc.ProjectID, topicID, avroSchemaID, pubsub.EncodingJSON)
			if err != nil {
				r.Errorf("createTopicWithSchema: %v", err)
			}
			got := buf.String()
			want := "Topic with schema created"
			if !strings.Contains(got, want) {
				r.Errorf("createTopicWithSchema mismatch\ngot: %v\nwant: %v\n", got, want)
			}

			subCfg := pubsub.SubscriptionConfig{
				Topic: client.Topic(topicID),
			}
			if _, err = client.CreateSubscription(ctx, subID, subCfg); err != nil {
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
			err = subscribeWithAvroSchemaRevisions(buf, tc.ProjectID, subID, avroFilePath)
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

	deleteSchema(ioutil.Discard, tc.ProjectID, avroSchemaID)
	client.Subscription(subID).Delete(ctx)
	client.Topic(topicID).Delete(ctx)
}

func TestSchemas_ProtoSchemaAll(t *testing.T) {
	client, _ := setup(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	topicID := topicPrefix + uuid.NewString()
	protoSchemaID := schemaPrefix + "proto-" + uuid.NewString()
	_, err := defaultSchemaConfig(tc.ProjectID, protoSchemaID, avroFilePath, pubsub.SchemaAvro)
	if err != nil {
		t.Fatalf("defaultSchemaConfig err: %v", err)
	}
	subID := subPrefix + uuid.NewString()

	t.Run("createResources", func(t *testing.T) {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := createProtoSchema(ioutil.Discard, tc.ProjectID, protoSchemaID, protoFilePath); err != nil {
				r.Errorf("createProtoSchema err: %v", err)
			}

			buf := new(bytes.Buffer)
			err := createTopicWithSchema(buf, tc.ProjectID, topicID, protoSchemaID, pubsub.EncodingJSON)
			if err != nil {
				r.Errorf("createTopicWithSchema: %v", err)
			}
			got := buf.String()
			want := "Topic with schema created"
			if !strings.Contains(got, want) {
				r.Errorf("createTopicWithSchema mismatch\ngot: %v\nwant: %v\n", got, want)
			}

			subCfg := pubsub.SubscriptionConfig{
				Topic: client.Topic(topicID),
			}
			if _, err = client.CreateSubscription(ctx, subID, subCfg); err != nil {
				r.Errorf("client.CreateSubscription err: %v", err)
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
			err := subscribeWithProtoSchema(buf, tc.ProjectID, subID, protoFilePath)
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

	deleteSchema(ioutil.Discard, tc.ProjectID, protoSchemaID)
	client.Subscription(subID).Delete(ctx)
	client.Topic(topicID).Delete(ctx)
}

func TestSchemas_UpdateTopicSchema(t *testing.T) {
	_, schemaClient := setup(t)
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	topicID := topicPrefix + uuid.NewString()
	protoSchemaID := schemaPrefix + "proto-" + uuid.NewString()
	protoSchemaID2 := schemaPrefix + "proto-" + uuid.NewString()

	protoSource, err := ioutil.ReadFile(protoFilePath)
	if err != nil {
		t.Fatalf("error reading from file: %s", protoFilePath)
	}
	schema, err := schemaClient.CreateSchema(ctx, protoSchemaID, pubsub.SchemaConfig{
		Type:       pubsub.SchemaProtocolBuffer,
		Definition: string(protoSource),
	})
	if err != nil {
		t.Fatalf("createProtoSchema err: %v", err)
	}

	_, err = schemaClient.CreateSchema(ctx, protoSchemaID2, pubsub.SchemaConfig{
		Type:       pubsub.SchemaProtocolBuffer,
		Definition: string(protoSource),
	})
	if err != nil {
		t.Fatalf("createProtoSchema err: %v", err)
	}

	if err := createTopicWithSchema(ioutil.Discard, tc.ProjectID, topicID, protoSchemaID, pubsub.EncodingJSON); err != nil {
		t.Fatalf("createTopicWithSchema: %v", err)
	}

	buf := new(bytes.Buffer)
	if err := updateTopicSchema(buf, tc.ProjectID, topicID, schema.RevisionID, schema.RevisionID); err != nil {
		t.Fatalf("updateTopicSchema err : %v", err)
	}
}

func defaultSchemaConfig(projectID, schemaID, schemaFile string, schemaType pubsub.SchemaType) (*pubsub.SchemaConfig, error) {
	schemaSource, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}
	cfg := &pubsub.SchemaConfig{
		Name:       fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
		Type:       schemaType,
		Definition: string(schemaSource),
	}
	return cfg, nil
}
