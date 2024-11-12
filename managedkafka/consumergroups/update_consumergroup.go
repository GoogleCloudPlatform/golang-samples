// Copyright 2024 Google LLC
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

package consumergroups

// [START managedkafka_update_consumergroup]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func updateConsumerGroup(w io.Writer, projectID, region, clusterID, consumerGroupID, topicPath string, partitionOffsets map[int32]int64, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// clusterID := "my-cluster"
	// consumerGroupID := "my-consumer-group"
	// topicPath := "my-topic-path"
	// partitionOffsets := map[int32]int64{1: 10, 2: 20, 3: 30}
	ctx := context.Background()
	client, err := managedkafka.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewClient got err: %w", err)
	}
	defer client.Close()

	clusterPath := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, region, clusterID)
	consumerGroupPath := fmt.Sprintf("%s/consumerGroups/%s", clusterPath, consumerGroupID)

	partitionMetadata := make(map[int32]*managedkafkapb.ConsumerPartitionMetadata)
	for partition, offset := range partitionOffsets {
		partitionMetadata[partition] = &managedkafkapb.ConsumerPartitionMetadata{
			Offset: offset,
		}
	}
	topicConfig := map[string]*managedkafkapb.ConsumerTopicMetadata{
		topicPath: {
			Partitions: partitionMetadata,
		},
	}
	consumerGroupConfig := managedkafkapb.ConsumerGroup{
		Name:   consumerGroupPath,
		Topics: topicConfig,
	}
	paths := []string{"topics"}
	updateMask := &fieldmaskpb.FieldMask{
		Paths: paths,
	}

	req := &managedkafkapb.UpdateConsumerGroupRequest{
		UpdateMask:    updateMask,
		ConsumerGroup: &consumerGroupConfig,
	}
	consumerGroup, err := client.UpdateConsumerGroup(ctx, req)
	if err != nil {
		return fmt.Errorf("client.UpdateConsumerGroup got err: %w", err)
	}
	fmt.Fprintf(w, "Updated consumer group: %#v\n", consumerGroup)
	return nil
}

// [END managedkafka_update_consumergroup]
