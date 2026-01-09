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
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"strings"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

// labels are used in operation-related logs.
const (
	labelOperationDeploy        = "deploy service"
	labelOperationBuild         = "build container image"
	labelOperationDeleteService = "delete service"
	labelOperationDeleteImage   = "delete container image"
	labelOperationGetURL        = "get url"
	defaultRegistryName         = "cloudrunci"
)

// HTTPGetProbe describes a probe definition using HTTP Get.
type HTTPGetProbe struct {
	Path string
	Port int
}

// GRPCProbe describes a probe definition using gRPC.
type GRPCProbe struct {
	Port    int
	Service string
}

// ReadinessProbe describes the readiness probe for a Cloud Run service.
type ReadinessProbe struct {
	TimeoutSeconds   int
	PeriodSeconds    int
	SuccessThreshold int
	FailureThreshold int
	HttpGet          *HTTPGetProbe
	GRPC             *GRPCProbe
}

// Service describes a Cloud Run service
type Service struct {
	// Name is an ID, used for logging and to generate a unique version to this run.
	Name string

	// The root directory containing the service's source code.
	Dir string

	// The container image name to deploy. If left blank the container will be built
	// and pushed to gcr.io/[ProjectID]/cloudrunci/[Name]:[Revision]
	Image string

	// The project to deploy to.
	ProjectID string

	// Allow unauthenticated request.
	AllowUnauthenticated bool

	// The platform to deploy to.
	Platform Platform

	// Additional runtime environment variable overrides for the app.
	Env EnvVars

	// Build the image without Dockerfile, using Google Cloud buildpacks.
	AsBuildpack bool

	// Strictly HTTP/2 serving
	HTTP2 bool

	deployed bool     // Whether the service has been deployed.
	built    bool     // Whether the container image has been built.
	url      *url.URL // The url of the deployed service.

	// Location to deploy the Service, and related artifacts
	Location string

	// Readiness probe definition for the containers in this service.
	Readiness *ReadinessProbe
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
		Location:  "us-central1",
	}
}

// Deployed reports whether the service has been deployed.
func (s *Service) Deployed() bool {
	return s.deployed
}

// RetryOptions holds options for Service.Request's retry behavior
type RetryOptions struct {
	MaxAttempts  int
	Delay        time.Duration
	ShouldAccept func(*http.Response) bool
}

func getDefaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxAttempts:  5,
		Delay:        20 * time.Second,
		ShouldAccept: Accept2xx,
	}
}

// Accept2xx returns true for responses in the 200 class of http response codes
func Accept2xx(r *http.Response) bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// AcceptNonServerError returns true for any non-500 http response
func AcceptNonServerError(r *http.Response) bool {
	return r.StatusCode < 500
}

func WithAttempts(n int) func(*RetryOptions) {
	return func(r *RetryOptions) {
		r.MaxAttempts = n
	}
}
func WithDelay(d time.Duration) func(*RetryOptions) {
	return func(r *RetryOptions) {
		r.Delay = d
	}
}
func WithAcceptFunc(f func(*http.Response) bool) func(*RetryOptions) {
	return func(r *RetryOptions) {
		r.ShouldAccept = f
	}
}

// Do executes the provided http.Request using the default http client
func (s *Service) Do(req *http.Request, opts ...func(*RetryOptions)) (*http.Response, error) {
	if !s.deployed {
		return nil, errors.New("Request called before Deploy")
	}
	options := getDefaultRetryOptions()
	for _, fn := range opts {
		fn(&options)
	}
	var lastSeen error
	resp := &http.Response{}
	for i := 0; i < options.MaxAttempts; i++ {
		defaultClient := &http.Client{}

		resp, lastSeen = defaultClient.Do(req)
		if lastSeen != nil {
			continue
		}
		if options.ShouldAccept(resp) {
			return resp, nil
		}
		time.Sleep(options.Delay)
	}
	// Too many attempts, return the last result.
	return resp, fmt.Errorf("no acceptable response after %d retries: %w", options.MaxAttempts, lastSeen)
}

// ImageRepoURL returns the base URL for building docker images
func (s *Service) ImageRepoURL() string {
	return fmt.Sprintf("%s-docker.pkg.dev/%s/%s", s.Location, s.ProjectID, defaultRegistryName)
}

// ensureDefaultImageRepo uses gcloud to create a default Image registry.
func (s *Service) ensureDefaultImageRepo() error {
	return ensureDefaultImageRepo(s.ProjectID, s.Location)
}

// Request issues an HTTP request to the deployed service.
func (s *Service) Request(method string, path string, opts ...func(*RetryOptions)) (*http.Response, error) {
	req, err := s.NewRequest(method, path)
	if err != nil {
		return &http.Response{}, err
	}
	return s.Do(req, opts...)
}

// NewRequest creates a new http.Request for the deployed service.
func (s *Service) NewRequest(method, path string) (*http.Request, error) {
	if !s.deployed {
		return nil, errors.New("NewRequest called before Deploy")
	}
	url, err := s.URL(path)
	if err != nil {
		return nil, fmt.Errorf("service.URL: %w", err)
	}
	return s.Platform.NewRequest(method, url)
}

// URL prepends the deployed service's base URL to the given path.
// Returns an error if the application has not been deployed.
func (s *Service) URL(p string) (string, error) {
	u, err := s.ParsedURL()
	if err != nil {
		return "", fmt.Errorf("service.ParsedURL: %w", err)
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
		return "", fmt.Errorf("service.ParsedURL: %w", err)
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
		out, err := gcloud(s.operationLabel(labelOperationGetURL), s.urlCmd())
		if err != nil {
			return nil, fmt.Errorf("gcloud: %s: %q", s.Name, err)
		}

		sURL := string(out)
		u, err := url.Parse(sURL)
		if err != nil {
			return nil, fmt.Errorf("url.Parse: %w", err)
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
func (s *Service) Version() string {
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

	if _, err := gcloud(s.operationLabel(labelOperationDeploy), s.deployCmd()); err != nil {
		return fmt.Errorf("gcloud: %s: %q", s.Version(), err)
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
		err := s.ensureDefaultImageRepo()
		if err != nil {
			return fmt.Errorf("failed to create image repository: %w", err)
		}
		s.Image = fmt.Sprintf("%s/%s:%s", s.ImageRepoURL(), s.Name, runID)
	}

	if out, err := gcloud(s.operationLabel(labelOperationBuild), s.buildCmd()); err != nil {
		log.Print(string(out))
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

	if _, err := gcloud(s.operationLabel(labelOperationDeleteService), s.deleteServiceCmd()); err != nil {
		return fmt.Errorf("gcloud: %v: %q", s.Version(), err)
	}
	s.deployed = false

	// If s.built is false no image was created or is not managed by cloudrun-ci.
	if s.built {
		_, err := gcloud(s.operationLabel("delete container image"), s.deleteImageCmd())
		if err != nil {
			return fmt.Errorf("gcloud: %v: %q", s.Version(), err)
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
		"alpha", // TODO until --use-http2 goes GA
		"run",
		"deploy",
		s.Version(),
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
	if s.HTTP2 {
		args = append(args, "--use-http2")
	}

	if s.Readiness != nil {
		var readinessProbeParts []string
		if s.Readiness.TimeoutSeconds > 0 {
			readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("timeoutSeconds=%d", s.Readiness.TimeoutSeconds))
		}
		if s.Readiness.PeriodSeconds > 0 {
			readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("periodSeconds=%d", s.Readiness.PeriodSeconds))
		}
		if s.Readiness.SuccessThreshold > 0 {
			readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("successThreshold=%d", s.Readiness.SuccessThreshold))
		}
		if s.Readiness.FailureThreshold > 0 {
			readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("failureThreshold=%d", s.Readiness.FailureThreshold))
		}
		if s.Readiness.HttpGet != nil {
			if s.Readiness.HttpGet.Path != "" {
				readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("httpGet.path=%s", s.Readiness.HttpGet.Path))
			}
			if s.Readiness.HttpGet.Port > 0 {
				readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("httpGet.port=%d", s.Readiness.HttpGet.Port))
			}
		} else if s.Readiness.GRPC != nil {
			if s.Readiness.GRPC.Service != "" {
				readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("grpc.service=%s", s.Readiness.GRPC.Service))
			}
			if s.Readiness.GRPC.Port > 0 {
				readinessProbeParts = append(readinessProbeParts, fmt.Sprintf("grpc.port=%d", s.Readiness.GRPC.Port))
			}
		}
		if len(readinessProbeParts) > 0 {
			args = append(args, "--readiness-probe="+strings.Join(readinessProbeParts, ","))
		}
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
		"beta", // TODO until --pack goes to GA
		"builds",
		"submit",
		"--project",
		s.ProjectID,
	}

	if !s.AsBuildpack {
		args = append(args, "--tag", s.Image)
	} else {
		args = append(args, "--pack=image="+s.Image)
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
		s.Version(),
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
		s.Version(),
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

func (s *Service) LogEntries(filter string, find string, maxAttempts int) (bool, error) {
	ctx := context.Background()
	client, err := logadmin.NewClient(ctx, s.ProjectID)
	if err != nil {
		return false, fmt.Errorf("logadmin.NewClient: %w", err)
	}
	defer client.Close()

	preparedFilter := fmt.Sprintf(`resource.type="cloud_run_revision" resource.labels.service_name="%s" %s`, s.Version(), filter)
	log.Printf("Using log filter: %s\n", preparedFilter)

	log.Println("Waiting for logs...")
	time.Sleep(3 * time.Minute)

	for i := 1; i < maxAttempts; i++ {
		log.Printf("Attempt #%d\n", i)
		it := client.Entries(ctx, logadmin.Filter(preparedFilter))
		for {
			entry, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return false, fmt.Errorf("it.Next: %w", err)
			}
			payload := fmt.Sprintf("%v", entry.Payload)
			if len(payload) > 0 {
				log.Printf("entry.Payload: %v\n", entry.Payload)
			}
			if strings.Contains(payload, find) {
				log.Printf("%q log entry found.\n", find)
				return true, nil
			}
		}
		time.Sleep(15 * time.Second)
	}
	return false, nil
}

// ensureDefaultImageRepo creates a default docker repo in the given project and location
// if it does not already exist.
func ensureDefaultImageRepo(project string, location string) error {
	cmd := exec.Command(gcloudBin,
		"artifacts", "repositories", "create", defaultRegistryName,
		"--project",
		project,
		"--repository-format=docker",
		"--location", location)
	o, err := gcloudWithoutRetry("ensure image repo", cmd)
	if err == nil {
		return nil
	}
	if strings.Contains(string(o), "ALREADY_EXISTS") {
		return nil
	}
	return err
}
