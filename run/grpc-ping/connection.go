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

package main

// [START cloudrun_grpc_conn]
// [START run_grpc_conn]

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewConn creates a new gRPC connection.
// host should be of the form domain:port, e.g., example.com:443
func NewConn(host string, insecure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}

	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	// Create a TokenSource
	// Using a TokenSource instead of a Token will ensure that the underlying Token
	// used will be refreshed when needed
	// A TokenSource automatically inject tokens with each underlying client request
	audience := "https://" + strings.Split(host, ":")[0]
	tokenSource, err := idtoken.NewTokenSource(context.Background(), audience,
		option.WithAudiences(audience))
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"NewTokenSource: %s", err,
		)
	}
	type grpcTokenSource struct {
		oauth.TokenSource
	}
	opts = append(opts, grpc.WithPerRPCCredentials(grpcTokenSource{
		TokenSource: oauth.TokenSource{
			TokenSource: tokenSource,
		},
	}))

	return grpc.Dial(host, opts...)
}

// [END run_grpc_conn]
// [END cloudrun_grpc_conn]
