// Copyright 2020 Google LLC
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

package functest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

// runID is an identifier that changes between runs.
var runID = time.Now().Format("20060102-150405")

// CloudFunction describes a Cloud Function.
type CloudFunction struct {
	// Name is an ID, with a numeric suffix provides the function name.
	Name string

	// Function name (if different from Name)
	Entrypoint string

	// The root directory containing the service's source code.
	Dir string

	// The function region.
	Region string

	// The functions runtime.
	Runtime string

	// The project to deploy to.
	ProjectID string

	// DeployCommand can override the default deploy command.
	DeployCommand string

	client     *http.Client // http client for requests.
	deployed   bool         // Whether the service has been deployed.
	deployName string       // Provides a uniqified name for deployment.
	url        *url.URL     // URL of the Cloud Function.

}

// NewCloudFunction creates a new CloudFunction.
// The entrypoint is derived from the name.
// The function name uses the provided name concatenated with a timestamp.
// Region defaults to us-central1, runtime defaults to go113.
func NewCloudFunction(name, projectID string) *CloudFunction {
	return &CloudFunction{
		deployName: name + runID,
		Entrypoint: name,
		Name:       name,
		ProjectID:  projectID,
		Region:     "us-central1",
		Runtime:    "go113",
	}
}

// Deploy the function to Cloud Functions.
func (f *CloudFunction) Deploy() error {
	f.log("deploying function...")

	if f.deployName == "" {
		f.deployName = f.Name + runID
	}

	if f.DeployCommand == "" {
		f.DeployCommand = fmt.Sprintf("gcloud functions deploy %s --trigger-http %s", f.DeployName(), f.deployFlagOverrides())
	} else {
		if strings.Contains(f.DeployCommand, "--runtime") || strings.Contains(f.DeployCommand, "--entrypoint") {
			f.log("--runtime and --entry-point flags are overridden from CloudFuntions.Runtime, CloudFunctions.Entrypoint")
		}

		// If --runtime or --entrypoint are set already they will be overridden.
		// These should be set by updating the CloudFunction struct property.
		// TODO: Parse the DeployCommand to configure Runtime and Entrypoint.
		f.DeployCommand = f.DeployCommand + " " + f.deployFlagOverrides()
	}

	cmd := exec.Command("/bin/sh", "-c", f.DeployCommand)
	cmd.Env = append(os.Environ(), f.gcloudEnvOverrides()...)
	if _, err := runCommand(cmd, f.Dir); err != nil {
		return fmt.Errorf("runCommand: %w", err)
	}
	f.deployed = true
	return nil
}

// Teardown deletes the deployed Cloud Function.
func (f *CloudFunction) Teardown() error {
	if !f.deployed {
		return errors.New("(no-op) called before deploy")
	}

	cmdStr := fmt.Sprintf("gcloud functions delete %s", f.deployName)
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Env = append(os.Environ(), f.gcloudEnvOverrides()...)

	f.log("Tearing down function...")
	if _, err := runCommand(cmd, f.Dir); err != nil {
		return fmt.Errorf("runCommand: %w", err)
	}
	f.deployed = false
	return nil
}

// deployFlagOverrides ensures the necessary deploy flags are present for deploy.
// Not available as gcloud properties so unavailable as environment variables.
func (f *CloudFunction) deployFlagOverrides() string {
	return fmt.Sprintf("--runtime %s --entry-point %s", f.Runtime, f.Entrypoint)
}

// gcloudEnvOverrides defines function-specific property overrides.
// Used to minimize complexity and modification of gcloud commands.
func (f *CloudFunction) gcloudEnvOverrides() []string {
	return []string{
		"CLOUDSDK_CORE_DISABLE_PROMPTS=TRUE",
		"CLOUDSDK_CORE_PROJECT=" + f.ProjectID,
		"CLOUDSDK_FUNCTIONS_REGION=" + f.Region,
	}
}

// Deployed reports whether the service has been deployed.
func (f *CloudFunction) Deployed() bool {
	return f.deployed
}

// DeployName is the uniq-ified name of the Cloud Function.
func (f *CloudFunction) DeployName() string {
	return f.deployName
}

func (f *CloudFunction) logf(format string, a ...interface{}) {
	f.log(fmt.Sprintf(format, a...))
}

func (f *CloudFunction) log(a ...interface{}) {
	log.Println(append([]interface{}{fmt.Sprintf("functest[%s]:", f.Name)}, a...)...)
}
