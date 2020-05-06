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

package cloudrunci

import (
	"errors"
	"fmt"
	"net/http"
)

// Platform describes how platforms are defined.
type Platform interface {
	Name() string
	CommandFlags() []string
	Validate() error
	NewRequest(method, url string) (*http.Request, error)
}

type platformBase struct{}

func (p platformBase) NewRequest(method, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}

// ManagedPlatform defines the Cloud Run (fully managed) hosting platform.
type ManagedPlatform struct {
	platformBase
	Region string
}

// Name retrieves the ID for the full managed platform.
func (p ManagedPlatform) Name() string {
	return "managed"
}

// Validate confirms required properties are set.
func (p ManagedPlatform) Validate() error {
	if p.Region == "" {
		return errors.New("Region missing")
	}
	return nil
}

// CommandFlags retrieves the common gcloud flags for targeting the full managed platform in a given region.
func (p ManagedPlatform) CommandFlags() []string {
	return []string{"--platform", "managed", "--region", p.Region}
}

// NewRequest creates an HTTP request for a URL on the platform.
func (p ManagedPlatform) NewRequest(method, url string) (*http.Request, error) {
	req, err := p.platformBase.NewRequest(method, url)
	if err == nil {
		token, err := CreateIDToken(url)
		if err != nil {
			return req, err
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	return req, err
}

// GKEPlatform defines the GKE platform for Cloud Run services.
type GKEPlatform struct {
	platformBase
	Cluster         string
	ClusterLocation string
}

// Name retrieves the ID for the GKE platform.
func (p GKEPlatform) Name() string {
	return "gke"
}

// Validate confirms required properties are set.
func (p GKEPlatform) Validate() error {
	if p.Cluster == "" {
		return errors.New("Cluster missing")
	}
	if p.ClusterLocation == "" {
		return errors.New("Cluster Location missing")
	}
	return nil
}

// CommandFlags retrieves the common gcloud flags for targeting a GKE cluster.
func (p GKEPlatform) CommandFlags() []string {
	return []string{"--platform", "gke", "--cluster", p.Cluster, "--cluster-location", p.ClusterLocation}
}

// KubernetesPlatform defines use of a knative-compatible kubernetes cluster beyond GKE.
type KubernetesPlatform struct {
	platformBase
	Kubeconfig string
	Context    string
}

// Name retrieves the ID for the Kubernetes platform.
func (p KubernetesPlatform) Name() string {
	return "kubernetes"
}

// Validate confirms required properties are set.
func (p KubernetesPlatform) Validate() error {
	if p.Kubeconfig == "" {
		return errors.New("kubeconfig missing")
	}
	if p.Context == "" {
		return errors.New("context missing")
	}
	return nil
}

// CommandFlags retrieves the common gcloud flags for targeting a Kubernetes cluster.
func (p KubernetesPlatform) CommandFlags() []string {
	return []string{"--platform", "gke", "--kubeconfig", p.Kubeconfig, "--context", p.Context}
}
