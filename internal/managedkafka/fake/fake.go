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

package fake

import (
	"context"
	"fmt"
	"net"
	"testing"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	longrunningpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
)

const (
	clusterID         = "fake-cluster"
	topicID           = "fake-topic"
	consumerGroupID   = "fake-consumergroup"
	connectClusterID  = "fake-connect-cluster"
	connectorID       = "fake-connector"
)

// The reason why we have a fake server is because testing end-to-end will exceed the deadline of 10 minutes.
// There is currently no strong support available for maintaining persistent resources either.
type fakeManagedKafkaServer struct {
	managedkafkapb.UnimplementedManagedKafkaServer
}

type fakeManagedKafkaConnectServer struct {
	managedkafkapb.UnimplementedManagedKafkaConnectServer
}

func Options(t *testing.T) []option.ClientOption {
	server := &fakeManagedKafkaServer{}
	connectServer := &fakeManagedKafkaConnectServer{}
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	gsrv := grpc.NewServer()
	managedkafkapb.RegisterManagedKafkaServer(gsrv, server)
	managedkafkapb.RegisterManagedKafkaConnectServer(gsrv, connectServer)
	fakeServerAddr := listener.Addr().String()
	go func() {
		if err := gsrv.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return []option.ClientOption{
		option.WithEndpoint(fakeServerAddr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
	}
}

func (f *fakeManagedKafkaServer) CreateCluster(ctx context.Context, req *managedkafkapb.CreateClusterRequest) (*longrunningpb.Operation, error) {
	anypb := &anypb.Any{}
	err := anypb.MarshalFrom(req.Cluster)
	if err != nil {
		return nil, fmt.Errorf("anypb.MarshalFrom got err: %w", err)
	}
	return &longrunningpb.Operation{
		Done: true,
		Result: &longrunningpb.Operation_Response{
			Response: anypb,
		},
	}, nil
}

func (f *fakeManagedKafkaServer) DeleteCluster(ctx context.Context, req *managedkafkapb.DeleteClusterRequest) (*longrunningpb.Operation, error) {
	return &longrunningpb.Operation{
		Done: true,
		Result: &longrunningpb.Operation_Response{
			Response: &anypb.Any{},
		},
	}, nil
}

func (f *fakeManagedKafkaServer) GetCluster(ctx context.Context, req *managedkafkapb.GetClusterRequest) (*managedkafkapb.Cluster, error) {
	return &managedkafkapb.Cluster{
		Name: clusterID,
	}, nil
}

func (f *fakeManagedKafkaServer) ListClusters(ctx context.Context, req *managedkafkapb.ListClustersRequest) (*managedkafkapb.ListClustersResponse, error) {
	return &managedkafkapb.ListClustersResponse{
		Clusters: []*managedkafkapb.Cluster{{
			Name: clusterID,
		}},
	}, nil
}

func (f *fakeManagedKafkaServer) UpdateCluster(ctx context.Context, req *managedkafkapb.UpdateClusterRequest) (*longrunningpb.Operation, error) {
	anypb := &anypb.Any{}
	err := anypb.MarshalFrom(req.Cluster)
	if err != nil {
		return nil, fmt.Errorf("anypb.MarshalFrom got err: %w", err)
	}
	return &longrunningpb.Operation{
		Done: true,
		Result: &longrunningpb.Operation_Response{
			Response: anypb,
		},
	}, nil
}

func (f *fakeManagedKafkaServer) CreateTopic(ctx context.Context, req *managedkafkapb.CreateTopicRequest) (*managedkafkapb.Topic, error) {
	return req.Topic, nil
}

func (f *fakeManagedKafkaServer) DeleteTopic(ctx context.Context, req *managedkafkapb.DeleteTopicRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeManagedKafkaServer) GetTopic(ctx context.Context, req *managedkafkapb.GetTopicRequest) (*managedkafkapb.Topic, error) {
	return &managedkafkapb.Topic{
		Name: topicID,
	}, nil
}

func (f *fakeManagedKafkaServer) ListTopics(ctx context.Context, req *managedkafkapb.ListTopicsRequest) (*managedkafkapb.ListTopicsResponse, error) {
	return &managedkafkapb.ListTopicsResponse{
		Topics: []*managedkafkapb.Topic{{
			Name: topicID,
		}},
	}, nil
}

func (f *fakeManagedKafkaServer) UpdateTopic(ctx context.Context, req *managedkafkapb.UpdateTopicRequest) (*managedkafkapb.Topic, error) {
	return &managedkafkapb.Topic{
		Name: topicID,
	}, nil
}

func (f *fakeManagedKafkaServer) DeleteConsumerGroup(ctx context.Context, req *managedkafkapb.DeleteConsumerGroupRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (f *fakeManagedKafkaServer) GetConsumerGroup(ctx context.Context, req *managedkafkapb.GetConsumerGroupRequest) (*managedkafkapb.ConsumerGroup, error) {
	return &managedkafkapb.ConsumerGroup{
		Name: consumerGroupID,
	}, nil
}

func (f *fakeManagedKafkaServer) ListConsumerGroups(ctx context.Context, req *managedkafkapb.ListConsumerGroupsRequest) (*managedkafkapb.ListConsumerGroupsResponse, error) {
	return &managedkafkapb.ListConsumerGroupsResponse{
		ConsumerGroups: []*managedkafkapb.ConsumerGroup{{
			Name: consumerGroupID,
		}},
	}, nil
}

func (f *fakeManagedKafkaServer) UpdateConsumerGroup(ctx context.Context, req *managedkafkapb.UpdateConsumerGroupRequest) (*managedkafkapb.ConsumerGroup, error) {
	return &managedkafkapb.ConsumerGroup{
		Name: consumerGroupID,
	}, nil
}

// Connect server methods
func (f *fakeManagedKafkaConnectServer) CreateConnectCluster(ctx context.Context, req *managedkafkapb.CreateConnectClusterRequest) (*longrunningpb.Operation, error) {
	anypb := &anypb.Any{}
	err := anypb.MarshalFrom(req.ConnectCluster)
	if err != nil {
		return nil, fmt.Errorf("anypb.MarshalFrom got err: %w", err)
	}
	return &longrunningpb.Operation{
		Done: true,
		Result: &longrunningpb.Operation_Response{
			Response: anypb,
		},
	}, nil
}

func (f *fakeManagedKafkaConnectServer) DeleteConnectCluster(ctx context.Context, req *managedkafkapb.DeleteConnectClusterRequest) (*longrunningpb.Operation, error) {
	return &longrunningpb.Operation{
		Done: true,
		Result: &longrunningpb.Operation_Response{
			Response: &anypb.Any{},
		},
	}, nil
}

func (f *fakeManagedKafkaConnectServer) GetConnectCluster(ctx context.Context, req *managedkafkapb.GetConnectClusterRequest) (*managedkafkapb.ConnectCluster, error) {
	return &managedkafkapb.ConnectCluster{
		Name: connectClusterID,
	}, nil
}

func (f *fakeManagedKafkaConnectServer) ListConnectClusters(ctx context.Context, req *managedkafkapb.ListConnectClustersRequest) (*managedkafkapb.ListConnectClustersResponse, error) {
	return &managedkafkapb.ListConnectClustersResponse{
		ConnectClusters: []*managedkafkapb.ConnectCluster{{
			Name: connectClusterID,
		}},
	}, nil
}

func (f *fakeManagedKafkaConnectServer) UpdateConnectCluster(ctx context.Context, req *managedkafkapb.UpdateConnectClusterRequest) (*longrunningpb.Operation, error) {
	anypb := &anypb.Any{}
	err := anypb.MarshalFrom(req.ConnectCluster)
	if err != nil {
		return nil, fmt.Errorf("anypb.MarshalFrom got err: %w", err)
	}
	return &longrunningpb.Operation{
		Done: true,
		Result: &longrunningpb.Operation_Response{
			Response: anypb,
		},
	}, nil
}

// Connector methods
func (f *fakeManagedKafkaConnectServer) CreateConnector(ctx context.Context, req *managedkafkapb.CreateConnectorRequest) (*managedkafkapb.Connector, error) {
	return req.Connector, nil
}

func (f *fakeManagedKafkaConnectServer) GetConnector(ctx context.Context, req *managedkafkapb.GetConnectorRequest) (*managedkafkapb.Connector, error) {
	return &managedkafkapb.Connector{
		Name: connectorID,
		Configs: map[string]string{
			"connector.class": "test.connector",
		},
	}, nil
}

func (f *fakeManagedKafkaConnectServer) ListConnectors(ctx context.Context, req *managedkafkapb.ListConnectorsRequest) (*managedkafkapb.ListConnectorsResponse, error) {
	return &managedkafkapb.ListConnectorsResponse{
		Connectors: []*managedkafkapb.Connector{{
			Name: connectorID,
			Configs: map[string]string{
				"connector.class": "test.connector",
			},
		}},
	}, nil
}

func (f *fakeManagedKafkaConnectServer) UpdateConnector(ctx context.Context, req *managedkafkapb.UpdateConnectorRequest) (*managedkafkapb.Connector, error) {
	return req.Connector, nil
}

func (f *fakeManagedKafkaConnectServer) DeleteConnector(ctx context.Context, req *managedkafkapb.DeleteConnectorRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (f *fakeManagedKafkaConnectServer) PauseConnector(ctx context.Context, req *managedkafkapb.PauseConnectorRequest) (*managedkafkapb.PauseConnectorResponse, error) {
	return &managedkafkapb.PauseConnectorResponse{}, nil
}

func (f *fakeManagedKafkaConnectServer) ResumeConnector(ctx context.Context, req *managedkafkapb.ResumeConnectorRequest) (*managedkafkapb.ResumeConnectorResponse, error) {
	return &managedkafkapb.ResumeConnectorResponse{}, nil
}

func (f *fakeManagedKafkaConnectServer) StopConnector(ctx context.Context, req *managedkafkapb.StopConnectorRequest) (*managedkafkapb.StopConnectorResponse, error) {
	return &managedkafkapb.StopConnectorResponse{}, nil
}

func (f *fakeManagedKafkaConnectServer) RestartConnector(ctx context.Context, req *managedkafkapb.RestartConnectorRequest) (*managedkafkapb.RestartConnectorResponse, error) {
	return &managedkafkapb.RestartConnectorResponse{}, nil
}
