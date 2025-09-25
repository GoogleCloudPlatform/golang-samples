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

package connectors

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/managedkafka/fake"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	connectorPrefix  = "connector"
	connectClusterID = "test-connect-cluster"
	region           = "us-central1"
)

func TestConnectors(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	connectorID := fmt.Sprintf("%s-%d", connectorPrefix, time.Now().UnixNano())
	options := fake.Options(t)

	t.Run("CreateMirrorMaker2SourceConnector", func(t *testing.T) {
		buf.Reset()
		sourceBootstrapServers := "source-dns:9092"
		targetBootstrapServers := "target-dns:9092"
		tasksMax := "3"
		sourceClusterAlias := "source"
		targetClusterAlias := "target"
		topics := ".*"
		topicsExclude := "mm2.*.internal,.*.replica,__.*"
		if err := createMirrorMaker2SourceConnector(buf, tc.ProjectID, region, connectClusterID, connectorID+"-mm2", sourceBootstrapServers, targetBootstrapServers, tasksMax, sourceClusterAlias, targetClusterAlias, topics, topicsExclude, options...); err != nil {
			t.Fatalf("failed to create MirrorMaker 2.0 source connector: %v", err)
		}
		got := buf.String()
		want := "Created MirrorMaker 2.0 Source connector"
		if !strings.Contains(got, want) {
			t.Fatalf("createMirrorMaker2SourceConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("CreatePubSubSourceConnector", func(t *testing.T) {
		buf.Reset()
		kafkaTopic := "test-topic"
		cpsSubscription := "test-subscription"
		tasksMax := "3"
		valueConverter := "org.apache.kafka.connect.converters.ByteArrayConverter"
		keyConverter := "org.apache.kafka.connect.storage.StringConverter"
		if err := createPubSubSourceConnector(buf, tc.ProjectID, region, connectClusterID, connectorID+"-pubsub-source", kafkaTopic, cpsSubscription, tc.ProjectID, tasksMax, valueConverter, keyConverter, options...); err != nil {
			t.Fatalf("failed to create Pub/Sub source connector: %v", err)
		}
		got := buf.String()
		want := "Created Pub/Sub source connector"
		if !strings.Contains(got, want) {
			t.Fatalf("createPubSubSourceConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("CreatePubSubSinkConnector", func(t *testing.T) {
		buf.Reset()
		topics := "test-topic"
		valueConverter := "org.apache.kafka.connect.storage.StringConverter"
		keyConverter := "org.apache.kafka.connect.storage.StringConverter"
		cpsTopic := "test-pubsub-topic"
		tasksMax := "3"
		if err := createPubSubSinkConnector(buf, tc.ProjectID, region, connectClusterID, connectorID+"-pubsub-sink", topics, valueConverter, keyConverter, cpsTopic, tc.ProjectID, tasksMax, options...); err != nil {
			t.Fatalf("failed to create Pub/Sub sink connector: %v", err)
		}
		got := buf.String()
		want := "Created Pub/Sub sink connector"
		if !strings.Contains(got, want) {
			t.Fatalf("createPubSubSinkConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("CreateCloudStorageSinkConnector", func(t *testing.T) {
		buf.Reset()
		topics := "test-topic"
		gcsBucketName := "test-bucket"
		tasksMax := "3"
		formatOutputType := "json"
		valueConverter := "org.apache.kafka.connect.json.JsonConverter"
		valueConverterSchemasEnable := "false"
		keyConverter := "org.apache.kafka.connect.storage.StringConverter"
		gcsCredentialsDefault := "true"
		if err := createCloudStorageSinkConnector(buf, tc.ProjectID, region, connectClusterID, connectorID+"-gcs-sink", topics, gcsBucketName, tasksMax, formatOutputType, valueConverter, valueConverterSchemasEnable, keyConverter, gcsCredentialsDefault, options...); err != nil {
			t.Fatalf("failed to create Cloud Storage sink connector: %v", err)
		}
		got := buf.String()
		want := "Created Cloud Storage sink connector"
		if !strings.Contains(got, want) {
			t.Fatalf("createCloudStorageSinkConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("CreateBigQuerySinkConnector", func(t *testing.T) {
		buf.Reset()
		topics := "test-topic"
		tasksMax := "3"
		keyConverter := "org.apache.kafka.connect.storage.StringConverter"
		valueConverter := "org.apache.kafka.connect.json.JsonConverter"
		valueConverterSchemasEnable := "false"
		defaultDataset := "test-dataset"
		if err := createBigQuerySinkConnector(buf, tc.ProjectID, region, connectClusterID, connectorID+"-bq-sink", topics, tasksMax, keyConverter, valueConverter, valueConverterSchemasEnable, defaultDataset, options...); err != nil {
			t.Fatalf("failed to create BigQuery sink connector: %v", err)
		}
		got := buf.String()
		want := "Created BigQuery sink connector"
		if !strings.Contains(got, want) {
			t.Fatalf("createBigQuerySinkConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("GetConnector", func(t *testing.T) {
		buf.Reset()
		if err := getConnector(buf, tc.ProjectID, region, connectClusterID, connectorID, options...); err != nil {
			t.Fatalf("failed to get connector: %v", err)
		}
		got := buf.String()
		want := "Got connector"
		if !strings.Contains(got, want) {
			t.Fatalf("getConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("UpdateConnector", func(t *testing.T) {
		buf.Reset()
		config := map[string]string{"tasks.max": "2"}
		if err := updateConnector(buf, tc.ProjectID, region, connectClusterID, connectorID, config, options...); err != nil {
			t.Fatalf("failed to update connector: %v", err)
		}
		got := buf.String()
		want := "Updated connector"
		if !strings.Contains(got, want) {
			t.Fatalf("updateConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("ListConnectors", func(t *testing.T) {
		buf.Reset()
		if err := listConnectors(buf, tc.ProjectID, region, connectClusterID, options...); err != nil {
			t.Fatalf("failed to list connectors: %v", err)
		}
		got := buf.String()
		want := "Got connector"
		if !strings.Contains(got, want) {
			t.Fatalf("listConnectors() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("PauseConnector", func(t *testing.T) {
		buf.Reset()
		if err := pauseConnector(buf, tc.ProjectID, region, connectClusterID, connectorID, options...); err != nil {
			t.Fatalf("failed to pause connector: %v", err)
		}
		got := buf.String()
		want := "Paused connector"
		if !strings.Contains(got, want) {
			t.Fatalf("pauseConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("ResumeConnector", func(t *testing.T) {
		buf.Reset()
		if err := resumeConnector(buf, tc.ProjectID, region, connectClusterID, connectorID, options...); err != nil {
			t.Fatalf("failed to resume connector: %v", err)
		}
		got := buf.String()
		want := "Resumed connector"
		if !strings.Contains(got, want) {
			t.Fatalf("resumeConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("StopConnector", func(t *testing.T) {
		buf.Reset()
		if err := stopConnector(buf, tc.ProjectID, region, connectClusterID, connectorID, options...); err != nil {
			t.Fatalf("failed to stop connector: %v", err)
		}
		got := buf.String()
		want := "Stopped connector"
		if !strings.Contains(got, want) {
			t.Fatalf("stopConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("RestartConnector", func(t *testing.T) {
		buf.Reset()
		if err := restartConnector(buf, tc.ProjectID, region, connectClusterID, connectorID, options...); err != nil {
			t.Fatalf("failed to restart connector: %v", err)
		}
		got := buf.String()
		want := "Restarted connector"
		if !strings.Contains(got, want) {
			t.Fatalf("restartConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})

	t.Run("DeleteConnector", func(t *testing.T) {
		buf.Reset()
		if err := deleteConnector(buf, tc.ProjectID, region, connectClusterID, connectorID, options...); err != nil {
			t.Fatalf("failed to delete connector: %v", err)
		}
		got := buf.String()
		want := "Deleted connector"
		if !strings.Contains(got, want) {
			t.Fatalf("deleteConnector() mismatch got: %s\nwant: %s", got, want)
		}
	})
}
