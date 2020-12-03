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

package transcoder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	location              = "us-central1"
	templateID            = "my-go-test-template"
	deleteTemplateReponse = "Deleted job template"
	deleteJobReponse      = "Deleted job"
	jobSucceededState     = "SUCCEEDED"
	testVideoFileName     = "ChromeCast.mp4"
	testVideoFileLocation = "../testdata/"
	preset                = "preset/web-hd"
)

// To run the tests, do the following:
// Export the following env vars:
// *   GOOGLE_APPLICATION_CREDENTIALS
// *   GOLANG_SAMPLES_PROJECT_ID
// Enable the following APIs on the test project:
// *   Transcoder API
// *   Cloud Resource Manager API (needed for project number translation)

// TestJobTemplatesAndJobs tests major operations on job templates
// and jobs.
func TestJobTemplatesAndJobs(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := tc.ProjectID + "-golang-samples-transcoder-test"
	inputURI := "gs://" + bucketName + "/" + testVideoFileName
	outputURIForPreset := "gs://" + bucketName + "/test-output-preset/"
	outputURIForTemplate := "gs://" + bucketName + "/test-output-template/"
	outputURIForAdHoc := "gs://" + bucketName + "/test-output-adhoc/"

	// Get the project number
	cloudresourcemanagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		t.Fatalf("cloudresourcemanager.NewService: %v", err)
	}
	project, err := cloudresourcemanagerClient.Projects.Get(tc.ProjectID).Do()
	if err != nil {
		t.Fatalf("cloudresourcemanagerClient.Projects.Get.Do: %v", err)
	}
	projectNumber := strconv.FormatInt(project.ProjectNumber, 10)

	testJobTemplates(t, projectNumber)
	t.Logf("\ntestJobTemplates() completed\n")
	writeTestGCSFile(t, tc.ProjectID, bucketName)
	t.Logf("\nwriteTestGCSFile() completed\n")
	testJobFromPreset(t, projectNumber, inputURI, outputURIForPreset)
	t.Logf("\ntestJobFromPreset() completed\n")
	testJobFromTemplate(t, projectNumber, inputURI, outputURIForTemplate)
	t.Logf("\ntestJobFromTemplate() completed\n")
	testJobFromAdHoc(t, projectNumber, inputURI, outputURIForAdHoc)
	t.Logf("\ntestJobFromAdHoc() completed\n")
}

// testJobTemplates tests major operations on job templates. Create, get,
// and list operations check if the template resource name is returned. The
// delete operation checks for a hard-coded string response.
func testJobTemplates(t *testing.T, projectNumber string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}

	// Remove the default template if it exists
	if err := getJobTemplate(buf, tc.ProjectID, location, templateID); err == nil {
		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			if err := deleteJobTemplate(buf, tc.ProjectID, location, templateID); err != nil {
				r.Errorf("deleteJobTemplate got err: %v", err)
			}
		})
	}

	// Create a new job template.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		templateName := fmt.Sprintf("projects/%s/locations/%s/jobTemplates/%s", projectNumber, location, templateID)
		if err := createJobTemplate(buf, tc.ProjectID, location, templateID); err != nil {
			r.Errorf("createJobTemplate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, templateName) {
			r.Errorf("createJobTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, templateName)
		}
	})

	// Get the new job template.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		templateName := fmt.Sprintf("projects/%s/locations/%s/jobTemplates/%s", projectNumber, location, templateID)
		if err := getJobTemplate(buf, tc.ProjectID, location, templateID); err != nil {
			r.Errorf("getJobTemplate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, templateName) {
			r.Errorf("getJobTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, templateName)
		}
	})

	// List the job templates for a given location.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		templateName := fmt.Sprintf("projects/%s/locations/%s/jobTemplates/%s", projectNumber, location, templateID)
		if err := listJobTemplates(buf, tc.ProjectID, location); err != nil {
			r.Errorf("listJobTemplates got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, templateName) {
			r.Errorf("listJobTemplates got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, templateName)
		}
	})

	// Delete the job template.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteJobTemplate(buf, tc.ProjectID, location, templateID); err != nil {
			r.Errorf("deleteJobTemplate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteTemplateReponse) {
			r.Errorf("deleteJobTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteTemplateReponse)
		}
	})
}

// writeTestGCSFile deletes the GCS test bucket and uploads a test video file to it.
func writeTestGCSFile(t *testing.T, projectID string, bucketName string) {
	t.Helper()
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	testutil.CleanBucket(ctx, t, projectID, bucketName)

	// Open local test file.
	f, err := os.Open(testVideoFileLocation + testVideoFileName)
	if err != nil {
		t.Fatalf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucketName).Object(testVideoFileName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("Writer.Close: %v", err)
	}
}

// testJobFromPreset tests major operations on a job created from a preset. It
// will wait until the job successfully completes as part of the test.
func testJobFromPreset(t *testing.T, projectNumber string, inputURI string, outputURIForPreset string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobFromPreset(buf, tc.ProjectID, location, inputURI, outputURIForPreset, preset); err != nil {
		t.Errorf("createJobFromPreset got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobFromPreset got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]

	// Get the job by job ID.
	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectNumber, location, jobID)
		if err := getJob(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJob got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, jobName) {
			r.Errorf("getJob got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
		}
	})

	// Get the job state (should be succeeded).
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, jobSucceededState) {
			r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
		}
	})

	// List the jobs for a given location.
	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectNumber, location, jobID)
		if err := listJobs(buf, tc.ProjectID, location); err != nil {
			r.Errorf("listJobs got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, jobName) {
			r.Errorf("listJobs got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
		}
	})

	// Delete the job.
	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		if err := deleteJob(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("deleteJob got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteJobReponse) {
			r.Errorf("deleteJob got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteJobReponse)
		}
	})
}

// testJobFromTemplate tests major operations on a job created from a template. It
// will wait until the job successfully completes as part of the test.
func testJobFromTemplate(t *testing.T, projectNumber string, inputURI string, outputURIForTemplate string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create a job template.
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		templateName := fmt.Sprintf("projects/%s/locations/%s/jobTemplates/%s", projectNumber, location, templateID)
		if err := createJobTemplate(buf, tc.ProjectID, location, templateID); err != nil {
			r.Errorf("createJobTemplate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, templateName) {
			r.Errorf("createJobTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, templateName)
		}
	})

	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobFromTemplate(buf, tc.ProjectID, location, inputURI, outputURIForTemplate, templateID); err != nil {
		t.Errorf("createJobFromTemplate got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobFromTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]

	// Get the job state (should be succeeded).
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, jobSucceededState) {
			r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
		}
	})

	// Delete the job.
	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		if err := deleteJob(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("deleteJob got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteJobReponse) {
			r.Errorf("deleteJob got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteJobReponse)
		}
	})

	// Delete the job template
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		if err := deleteJobTemplate(buf, tc.ProjectID, location, templateID); err != nil {
			r.Errorf("deleteJobTemplate got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteTemplateReponse) {
			r.Errorf("deleteJobTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteTemplateReponse)
		}
	})
}

// testJobFromAdHoc tests major operations on a job created from an ad-hoc configuration. It
// will wait until the job successfully completes as part of the test.
func testJobFromAdHoc(t *testing.T, projectNumber string, inputURI string, outputURIForAdHoc string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobFromAdHoc(buf, tc.ProjectID, location, inputURI, outputURIForAdHoc); err != nil {
		t.Errorf("createJobFromAdHoc got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobFromAdHoc got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]

	// Get the job.
	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectNumber, location, jobID)
		if err := getJob(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJob got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, jobName) {
			r.Errorf("getJob got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
		}
	})

	// Get the job state (should be succeeded).
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, jobSucceededState) {
			r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
		}
	})

	// Delete the job.
	testutil.Retry(t, 3, 5*time.Second, func(r *testutil.R) {
		if err := deleteJob(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("deleteJob got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, deleteJobReponse) {
			r.Errorf("deleteJob got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, deleteJobReponse)
		}
	})
}
