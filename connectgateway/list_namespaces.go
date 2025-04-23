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
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	gateway "cloud.google.com/go/gkeconnect/gateway/apiv1"
	gatewaypb "cloud.google.com/go/gkeconnect/gateway/apiv1/gatewaypb"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// --- Configuration ---
var (
	scopes = "https://www.googleapis.com/auth/cloud-platform"
)

func listNamespaces(projectID, membershipID, membershipLocation, serviceAccountKeyPath string) (*v1.NamespaceList, error) {
	ctx := context.Background()

	gatewayURL, err := getGatewayURL(ctx, projectID, membershipID, membershipLocation)
	if err != nil {
		return nil, fmt.Errorf("error fetching GKE Connect Gateway URL: %v", err)
	}

	kubeClient, err := configureKubernetesClient(ctx, gatewayURL, serviceAccountKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error configuring Kubernetes client: %v", err)
	}

	return callListNamespaces(kubeClient)
}

func getGatewayURL(ctx context.Context, projectID, membershipID, membershipLocation string) (string, error) {
	gatewayClient, err := gateway.NewGatewayControlRESTClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create GKE Connect Gateway client: %w", err)
	}
	defer gatewayClient.Close()

	req := &gatewaypb.GenerateCredentialsRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/memberships/%s", projectID, membershipLocation, membershipID),
	}

	resp, err := gatewayClient.GenerateCredentials(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", fmt.Errorf("membership not found: %w", err)
		}
		return "", fmt.Errorf("failed to fetch GKE Connect Gateway URL: %w", err)
	}

	fmt.Printf("GKE Connect Gateway Endpoint: %s\n", resp.Endpoint)
	if resp.Endpoint == "" {
		return "", fmt.Errorf("error: GKE Connect Gateway Endpoint is empty")
	}
	return resp.Endpoint, nil
}

func configureKubernetesClient(ctx context.Context, gatewayURL string, serviceAccountKeyPath string) (*kubernetes.Clientset, error) {
	// Read the service account key file
	keyBytes, err := ioutil.ReadFile(serviceAccountKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading service account key file: %v", err)
	}

	// Create Google credentials from the service account key
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

func callListNamespaces(clientset *kubernetes.Clientset) (*v1.NamespaceList, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	return namespaces, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: Must include 4 arguments for projectID, membershipID, membershipLocation, and serviceAccountKey")
		os.Exit(1)
	}
	projectNumber := os.Args[1]
	membershipID := os.Args[2]
	membershipLocation := os.Args[3]
	serviceAccountKeyPath := os.Args[4]

	namespaces, err := listNamespaces(projectNumber, membershipID, membershipLocation, serviceAccountKeyPath)
	if err != nil {
		log.Fatalf("listNamespaces: %v", err)
	}
	if len(namespaces.Items) > 0 {
		fmt.Println("\n--- List of Namespaces ---")
		for _, namespace := range namespaces.Items {
			fmt.Printf("Name: %s\n", namespace.ObjectMeta.Name)
		}
	} else {
		fmt.Println("No namespaces found in the cluster.")
	}
}
