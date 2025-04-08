// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package snippets

// [START memorystore_connect_to_database]
import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"time"

	memorystore "cloud.google.com/go/redis/apiv1"
	redispb "cloud.google.com/go/redis/apiv1/redispb"
	"github.com/go-redis/redis/v8"
)

// ConnectToDatabase demonstrates how to use go-redis library to connect to a
// Memorystore Redis instance.
func ConnectToDatabase(w io.Writer, projectID, location, instanceID string) error {

	// Instantiate a Redis administrative client
	ctx := context.Background()
	adminClient, err := memorystore.NewCloudRedisClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	req := &redispb.GetInstanceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID),
	}

	instance, err := adminClient.GetInstance(ctx, req)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, instance)

	// Load CA cert
	caCerts := instance.GetServerCaCerts()
	if len(caCerts) == 0 {
		return errors.New("memorystore: no server CA certs for instance")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(caCerts[0].Cert))

	// Setup Redis Connection pool
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", instance.Host, instance.Port),
		Password:     "PASSWORD",
		PoolSize:     1,
		MinIdleConns: 1,
		PoolTimeout:  0,
		IdleTimeout:  20 * time.Second,
		DialTimeout:  2 * time.Second,
		TLSConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	})

	p, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Response:\n%s", p)

	return nil
}

// [END memorystore_connect_to_database]
