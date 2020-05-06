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

// Package cloudrunci facilitates end-to-end testing against the production Cloud Run.
//
// This is a specialized tool that could be used in addition to unit tests. It
// calls the `gcloud beta run` command directly.
//
// gcloud (https://cloud.google.com/sdk) must be installed. You must be authorized via
// the gcloud command-line tool (`gcloud auth login`).
//
// You may specify the location of gcloud via the GCLOUD_BIN environment variable.
package cloudrunci

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"time"
)

// Service describes a Cloud Run service
type Service struct {
	// Name is an ID, used for logging and to generate a unique version to this run.
	Name string

	// The root directory containing the service's source code.
	Dir string

	// The container image name to deploy. If left blank the container will be built
	// and pushed to gcr.io/[ProjectID]/[Name]:[Revision]
	Image string

	// The project to deploy to.
	ProjectID string

	// Allow unauthenticated request.
	AllowUnauthenticated bool

	// The platform to deploy to.
	Platform Platform

	// Additional runtime environment variable overrides for the app.
	Env EnvVars

	deployed bool     // Whether the service has been deployed.
	built    bool     // Whether the container image has been built.
	url      *url.URL // The url of the deployed service.
}

// runID is an identifier that changes between runs.
var runID = time.Now().Format("20060102-150405")

// NewService creates a new Service based on the name and projectID provided.
// It will default to the ManagedPlatform in region us-central1,
// and build a container image as needed  for deployment.
func NewService(name, projectID string) *Service {
	return &Service{
		Name:      name,
		ProjectID: projectID,
		Platform:  ManagedPlatform{Region: "us-central1"},
	}
}

// Deployed reports whether the service has been deployed.
func (s *Service) Deployed() bool {
	return s.deployed
}

// Request issues an HTTP request to the deployed service.
func (s *Service) Request(method, path string) (*http.Response, error) {
	if !s.deployed {
		return nil, errors.New("Request called before Deploy")
	}
	req, err := s.NewRequest(method, path)
	if err != nil {
		return nil, err
	}
	defaultClient := &http.Client{}

	return defaultClient.Do(req)
}

// NewRequest creates a new http.Request for the deployed service.
func (s *Service) NewRequest(method, path string) (*http.Request, error) {
	if !s.deployed {
		return nil, errors.New("NewRequest called before Deploy")
	}
	url, err := s.URL(path)
	if err != nil {
		return nil, fmt.Errorf("service.URL: %v", err)
	}
	return s.Platform.NewRequest(method, url)
}

// URL prepends the deployed service's base URL to the given path.
// Returns an error if the application has not been deployed.
func (s *Service) URL(p string) (string, error) {
	u, err := s.ParsedURL()
	if err != nil {
		return "", fmt.Errorf("service.ParsedURL: %v", err)
	}
	modified := &url.URL{}
	*modified = *u
	modified.Path = path.Join(modified.Path, p)

	return modified.String(), nil
}

// Host returns the host:port of the service to facilitate new gRPC connections.
func (s *Service) Host() (string, error) {
	u, err := s.ParsedURL()
	if err != nil {
		return "", fmt.Errorf("service.ParsedURL: %v", err)
	}
	return u.Host + ":443", nil
}

// ParsedURL retrieves the parsed URL of the service.
// This URL is stored on the service struct for repeated retrieval.
func (s *Service) ParsedURL() (*url.URL, error) {
	if !s.deployed {
		return nil, errors.New("URL called before Deploy")
	}
	if s.url == nil {
		out, err := gcloud(s.operationLabel("get url"), s.urlCmd())
		if err != nil {
			return nil, fmt.Errorf("gcloud: %s: %q", s.Name, err)
		}

		sURL := string(out)
		u, err := url.Parse(sURL)
		if err != nil {
			return nil, fmt.Errorf("url.Parse: %v", err)
		}

		s.url = u
	}
	return s.url, nil
}

// validate confirms all required service properties are present.
func (s *Service) validate() error {
	if s.ProjectID == "" {
		return errors.New("Project ID missing")
	}
	if s.Platform == nil {
		return errors.New("Platform configuration missing")
	}
	if err := s.Platform.Validate(); err != nil {
		return err
	}
	if err := s.Env.Validate(); err != nil {
		return err
	}

	return nil
}

// revision returns the revision that the service will be deployed to.
// NOTE: Until traffic splitting is available, this will be used as the service name.
func (s *Service) version() string {
	return s.Name + "-" + runID
}

// Deploy deploys the service to Cloud Run.
// If an image has not been specified or previously built, it will call Build.
// If the deployment fails, it tries to clean up the failed deployment.
func (s *Service) Deploy() error {
	// Don't deploy unless we're certain everything is ready for deployment
	// (i.e. admin client is authenticated and authorized)
	if err := s.validate(); err != nil {
		return err
	}

	if s.Image == "" && !s.built {
		if err := s.Build(); err != nil {
			return err
		}
	}

	if _, err := gcloud(s.operationLabel("deploy service"), s.deployCmd()); err != nil {
		return fmt.Errorf("gcloud: %s: %q", s.version(), err)
	}

	s.deployed = true
	return nil
}

// Build builds a container image if one has not already been built.
// If service.Image is specified and this is directly called, it
// could overwrite an existing container image.
// If service.Image is not specified, the service.Deploy() function will
// call Build().
func (s *Service) Build() error {
	if err := s.validate(); err != nil {
		return err
	}
	if s.built {
		return fmt.Errorf("container image already built")
	}
	if s.Image == "" {
		s.Image = fmt.Sprintf("gcr.io/%s/%s:%s", s.ProjectID, s.Name, runID)
	}

	if _, err := gcloud(s.operationLabel("build container image"), s.buildCmd()); err != nil {
		return fmt.Errorf("gcloud: %s: %q", s.Image, err)
	}
	s.built = true

	return nil
}

// Clean deletes the created Cloud Run service.
func (s *Service) Clean() error {
	// NOTE: don't check whether p.deployed is set.
	// We may want to attempt to clean up if deployment failed.

	if err := s.validate(); err != nil {
		return err
	}

	if _, err := gcloud(s.operationLabel("delete service"), s.deleteServiceCmd()); err != nil {
		return fmt.Errorf("gcloud: %v: %q", s.version(), err)
	}
	s.deployed = false

	// If s.built is false no image was created or is not managed by cloudrun-ci.
	if s.built {
		_, err := gcloud(s.operationLabel("delete container image"), s.deleteImageCmd())
		if err != nil {
			return fmt.Errorf("gcloud: %v: %q", s.version(), err)
		}
		s.built = false
	}

	return nil
}

func (s *Service) operationLabel(op string) string {
	return fmt.Sprintf("operation [%s] for service [%s]", op, s.Name)
}

func (s *Service) deployCmd() *exec.Cmd {
	args := append([]string{
		"--quiet",
		"run",
		"deploy",
		s.version(),
		"--project",
		s.ProjectID,
		"--image",
		s.Image,
	}, s.Platform.CommandFlags()...)

	if s.Env != nil {
		for k := range s.Env {
			args = append(args, "--set-env-vars", s.Env.Variable(k))
		}
	}
	if s.AllowUnauthenticated {
		args = append(args, "--allow-unauthenticated")
	}

	// NOTE: if the "beta" component is not available, and this is run in parallel,
	// gcloud will attempt to install those components multiple
	// times and will eventually fail on IO.
	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = s.Dir
	return cmd
}

func (s *Service) buildCmd() *exec.Cmd {
	args := []string{
		"--quiet",
		"builds",
		"submit",
		"--project",
		s.ProjectID,
		"--tag",
		s.Image,
	}

	// NOTE: if the "beta" component is not available, and this is run in parallel,
	// gcloud will attempt to install those components multiple
	// times and will eventually fail on IO.
	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = s.Dir
	return cmd
}

func (s *Service) deleteImageCmd() *exec.Cmd {
	args := []string{
		"--quiet",
		"container",
		"images",
		"delete",
		s.Image,
	}

	// NOTE: if the "beta" component is not available, and this is run in parallel,
	// gcloud will attempt to install those components multiple
	// times and will eventually fail on IO.
	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = s.Dir
	return cmd
}

func (s *Service) deleteServiceCmd() *exec.Cmd {
	args := append([]string{
		"--quiet",
		"run",
		"services",
		"delete",
		s.version(),
		"--project",
		s.ProjectID,
	}, s.Platform.CommandFlags()...)

	// NOTE: if the "beta" component is not available, and this is run in parallel,
	// gcloud will attempt to install those components multiple
	// times and will eventually fail on IO.
	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = s.Dir
	return cmd
}

func (s *Service) urlCmd() *exec.Cmd {
	args := append([]string{
		"--quiet",
		"run",
		"services",
		"describe",
		s.version(),
		"--project",
		s.ProjectID,
		"--format",
		"value(status.url)",
	}, s.Platform.CommandFlags()...)

	// NOTE: if the "beta" component is not available, and this is run in parallel,
	// gcloud will attempt to install those components multiple
	// times and will eventually fail on IO.
	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = s.Dir
	return cmd
}
