// Copyright 2022 Google LLC
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
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

// Job describes a Cloud Run Job
// The typical usage flow of a Job is to call the following methods, which
// call the corresponding "gcloud run jobs" commands:
// Build(), Create(), Run().
// Note: The LogEntries() method cannot differentiate between executions at this
// time, so it is not recommended to call Run() multiple times on a single Job
// object.
type Job struct {
	// Name is an ID, used for logging and to generate a unique version to this run.
	Name string

	// The directory containing the job's code.
	Dir string

	// The container image name to deploy. If left blank the container will be built
	// and pushed to gcr.io/[ProjectID]/[Name]:[Revision]
	Image string

	// The project to deploy to.
	ProjectID string

	// The Region to deploy to.
	// We do not use Platform objects, since Jobs do not have a '--platform'
	// flag in gcloud.
	Region string

	// Additional runtime environment variable overrides for the app.
	Env EnvVars

	// Additional flags to be passed at Job creation.
	ExtraCreateFlags []string

	// Build this Image as a BuildPack, without using a Dockerfile
	AsBuildpack bool

	built   bool // True if container image has been built.
	created bool // True if job has been created.
	started bool // true if the Job has been started.
}

// NewJob creates a new Job to be run with Cloud Run Jobs.
// It will default to the ManagedPlatform in region us-central1,
// and build a container image as needed for deployment.
func NewJob(name, projectID string) *Job {
	return &Job{
		Name:      name,
		ProjectID: projectID,
		Region:    "us-central1",
	}
}

func (j *Job) CommonGCloudFlags() []string {
	return []string{
		"--region", j.Region,
		"--project", j.ProjectID,
	}

}

// validate confirms all required job properties are present.
func (j *Job) validate() error {
	if j.ProjectID == "" {
		return errors.New("Project ID missing")
	}
	if j.Region == "" {
		return errors.New("Region missing")
	}
	if err := j.Env.Validate(); err != nil {
		return err
	}

	return nil
}

// version returns the execution that the service will be deployed to.
// This identifier is also used to locate relevant log messages.
func (j *Job) version() string {
	return j.Name + "-" + runID
}

// Creates the Cloud Run job, but does not start it.
// If an image has not been specified or previously built, it will call Build.
func (j *Job) Create() error {
	// Don't deploy unless we're certain everything is ready for deployment
	// (i.e. admin client is authenticated and authorized)
	if err := j.validate(); err != nil {
		return err
	}

	if j.Image == "" && !j.built {
		if err := j.Build(); err != nil {
			return err
		}
	}

	if _, err := gcloud(fmt.Sprintf("%s: Creating Cloud Run Job", j.version()), j.createCmd()); err != nil {
		return fmt.Errorf("gcloud: %s: %q", j.version(), err)
	}

	j.created = true
	return nil
}

// Build builds a container image if one has not already been built.
// If service.Image is specified and this is directly called, it
// could overwrite an existing container image.
func (j *Job) Build() error {
	if err := j.validate(); err != nil {
		return err
	}
	if j.built {
		return fmt.Errorf("container image already built")
	}
	if j.Image == "" {
		ensureDefaultImageRepo(j.ProjectID, j.Region)
		j.Image = fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s:%s",
			j.Region, j.ProjectID, defaultRegistryName, j.Name, runID)
	}

	if _, err := gcloud(fmt.Sprintf("%s: Building image %s", j.version(), j.Image), j.buildCmd()); err != nil {
		return fmt.Errorf("gcloud: %s: %q", j.Image, err)
	}
	j.built = true

	return nil
}

// Run starts the Job in Cloud Run Jobs.
// This method will call Build and Create if necessary.
func (j *Job) Run() error {
	if err := j.validate(); err != nil {
		return err
	}
	// Create() checks that the image was built
	if !j.created {
		if err := j.Create(); err != nil {
			return err
		}
	}
	if _, err := gcloud(fmt.Sprintf("%s: Running cloud run job", j.version()), j.runCmd()); err != nil {
		return fmt.Errorf("gcloud: %v: %q", j.version(), err)
	}
	return nil
}

// Clean deletes the created Cloud Run service.
func (j *Job) Clean() error {
	// NOTE: don't check whether j.created is set.
	// We may want to attempt to clean up if creation failed (i.e. to clean up the image).

	if err := j.validate(); err != nil {
		return err
	}

	if _, err := gcloud(fmt.Sprintf("%s: Deleting cloud run job", j.version()), j.deleteJobCmd()); err != nil {
		return fmt.Errorf("gcloud: %v: %q", j.version(), err)
	}
	j.created = false

	// If built is false, no image was created or is not managed by cloudrun-ci.
	if j.built {
		_, err := gcloud(fmt.Sprintf("%s: Deleting Image %s", j.version(), j.Image), j.deleteImageCmd())
		if err != nil {
			return fmt.Errorf("gcloud: %v: %q", j.version(), err)
		}
		j.built = false
	}

	return nil
}

func (j *Job) createCmd() *exec.Cmd {
	args := append([]string{
		"--quiet",
		"alpha",
		"run",
		"jobs",
		"create",
		j.version(),
		"--image",
		j.Image,
	}, j.CommonGCloudFlags()...)

	if j.Env != nil {
		for k := range j.Env {
			args = append(args, "--set-env-vars", j.Env.Variable(k))
		}
	}

	args = append(args, j.ExtraCreateFlags...)

	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = j.Dir
	return cmd
}

func (j *Job) buildCmd() *exec.Cmd {
	args := []string{
		"--quiet",
		"builds",
		"submit",
		"--project",
		j.ProjectID,
	}
	if j.AsBuildpack {
		args = append(args, "--pack=image="+j.Image)
	} else {
		args = append(args, "--tag", j.Image)
	}

	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = j.Dir
	return cmd
}

// runCmd returns the gcloud command needed to start this RunJob
func (j *Job) runCmd() *exec.Cmd {
	args := append([]string{
		"--quiet",
		"alpha",
		"run",
		"jobs",
		"execute",
		j.version(),
		"--wait", // Waits for job to complete before returning.
	}, j.CommonGCloudFlags()...)

	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = j.Dir
	return cmd
}

func (j *Job) deleteImageCmd() *exec.Cmd {
	args := []string{
		"--quiet",
		"container",
		"images",
		"delete",
		j.Image,
		"--force-delete-tags",
	}

	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = j.Dir
	return cmd
}

func (j *Job) deleteJobCmd() *exec.Cmd {
	args := append([]string{
		"--quiet",
		"alpha",
		"run",
		"jobs",
		"delete",
		j.version(),
	}, j.CommonGCloudFlags()...)

	cmd := exec.Command(gcloudBin, args...)
	cmd.Dir = j.Dir
	return cmd
}

func (j *Job) LogEntries(filter string, find string, maxAttempts int) (bool, error) {
	ctx := context.Background()
	client, err := logadmin.NewClient(ctx, j.ProjectID)
	if err != nil {
		return false, fmt.Errorf("logadmin.NewClient: %w", err)
	}
	defer client.Close()

	preparedFilter := fmt.Sprintf(`resource.type="cloud_run_job" resource.labels.job_name="%s" %s`, j.version(), filter)
	fmt.Printf("Using log filter: %s\n", preparedFilter)

	for i := 1; i < maxAttempts; i++ {
		fmt.Printf("Attempt #%d\n", i)
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
				fmt.Printf("entry.Payload: %v\n", entry.Payload)
			}
			if strings.Contains(payload, find) {
				fmt.Printf("%q log entry found.\n", find)
				return true, nil
			}
		}
		time.Sleep(30 * time.Second)
	}
	return false, nil
}
