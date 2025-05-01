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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	gateway "cloud.google.com/go/gkeconnect/gateway/apiv1"
	gatewaypb "cloud.google.com/go/gkeconnect/gateway/apiv1/gatewaypb"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	scopes = "https://www.googleapis.com/auth/cloud-platform"
)

func getNamespace(membershipName, membershipLocation, serviceAccountKeyPath string) (*v1.Namespace, error) {
	ctx := context.Background()

	gatewayURL, err := getGatewayURL(ctx, membershipName, membershipLocation)
	if err != nil {
		return nil, fmt.Errorf("error fetching Connect Gateway URL: %v", err)
	}

	kubeClient, err := configureKubernetesClient(ctx, gatewayURL, serviceAccountKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error configuring Kubernetes client: %v", err)
	}

	return callGetNamespace(kubeClient)
}

func getGatewayURL(ctx context.Context, membershipName, membershipLocation string) (string, error) {
	var opts option.ClientOption
	if membershipLocation != "global" {
		opts = option.WithEndpoint(fmt.Sprintf("%v-connectgateway.googleapis.com", membershipLocation))
	}
	gatewayClient, err := gateway.NewGatewayControlRESTClient(ctx, opts)
	if err != nil {
		return "", fmt.Errorf("failed to create Connect Gateway client: %w", err)
	}
	defer gatewayClient.Close()

	req := &gatewaypb.GenerateCredentialsRequest{
		Name: membershipName,
	}

	resp, err := gatewayClient.GenerateCredentials(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Connect Gateway URL for membership %s: %w", membershipName, err)
	}

	fmt.Printf("Connect Gateway Endpoint: %s\n", resp.Endpoint)
	return resp.Endpoint, nil
}

func configureKubernetesClient(ctx context.Context, gatewayURL string, serviceAccountKeyPath string) (*kubernetes.Clientset, error) {
	// Read the service account key file.
	keyBytes, err := ioutil.ReadFile(serviceAccountKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading service account key file: %v", err)
	}

	// Create Google credentials from the service account key.
	creds, err := google.CredentialsFromJSON(context.Background(), keyBytes, scopes)
	if err != nil {
		return nil, fmt.Errorf("error creating credentials: %v", err)
	}
	config := &rest.Config{
		Host: gatewayURL,
		WrapTransport: func(rt http.RoundTripper) http.RoundTripper {
			return &oauth2.Transport{
				Source: creds.TokenSource,
				Base:   rt,
			}
		},
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}
	return clientset, nil
}

func callGetNamespace(clientset *kubernetes.Clientset) (*v1.Namespace, error) {
	namespace, err := clientset.CoreV1().Namespaces().Get(context.Background(), "default", metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}
	return namespace, nil
}

func main() {
	membershipName := os.Getenv("MEMBERSHIP_NAME")
	membershipLocation := os.Getenv("MEMBERSHIP_LOCATION")
	serviceAccountKeyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if _, err := os.Stat(serviceAccountKeyPath); os.IsNotExist(err) {
		fmt.Printf("service account key file not found at %s\n", serviceAccountKeyPath)
		os.Exit(1)
	}

	namespace, err := getNamespace(membershipName, membershipLocation, serviceAccountKeyPath)
	if err != nil {
		fmt.Printf("failed to get namespace: %v", err)
		os.Exit(1)
	}
	fmt.Printf("\nDefault Namespace:\n%v", namespace)
}
