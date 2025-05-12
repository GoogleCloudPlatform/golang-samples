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

package connectgateway

// [START connectgateway_get_namespace]

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	gateway "cloud.google.com/go/gkeconnect/gateway/apiv1"
	gatewaypb "cloud.google.com/go/gkeconnect/gateway/apiv1/gatewaypb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// getNamespace retrieves the Connect Gateway URL associated with the input
// membership. It then creates a kubernetes client using the retrieved Gateway
// URL to make requests to the underlying cluster.
func getNamespace(w io.Writer, membershipName, region string) error {
	// Use Gateway Control to retrieve the Connect Gateway URL to be used as the
	// host of the kubernetes client.
	ctx := context.Background()
	// If the membership location is regional, then the regional endpoint needs to be set for the client.
	// Global memberships do not require this override as the default endpoint is global.
	opts := option.WithEndpoint(fmt.Sprintf("%v-connectgateway.googleapis.com", region))
	gatewayClient, err := gateway.NewGatewayControlRESTClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create Connect Gateway client: %v", err)
	}
	defer gatewayClient.Close()

	req := &gatewaypb.GenerateCredentialsRequest{
		Name: membershipName,
	}
	resp, err := gatewayClient.GenerateCredentials(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to fetch Connect Gateway URL for membership %s: %w", membershipName, err)
	}
	gatewayURL := resp.Endpoint
	fmt.Printf("Connect Gateway Endpoint: %s\n", gatewayURL)

	// Configure the kubernetes client library using the Connect Gateway URL and
	// application default credentials.
	scopes := "https://www.googleapis.com/auth/cloud-platform"
	tokenSource, err := google.DefaultTokenSource(ctx, scopes)
	if err != nil {
		return fmt.Errorf("failed to get default credentials: %w", err)
	}
	wrapTransport := func(rt http.RoundTripper) http.RoundTripper {
		return &oauth2.Transport{
			Source: tokenSource,
			Base:   rt,
		}
	}
	config := &rest.Config{
		Host:          gatewayURL,
		WrapTransport: wrapTransport,
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Call GetNamespace using the kubernetes client.
	namespace, err := kubeClient.CoreV1().Namespaces().Get(context.Background(), "default", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}
	fmt.Fprintf(w, "\nDefault Namespace:\n%#v", namespace)
	return nil
}

// [END connectgateway_get_namespace]
