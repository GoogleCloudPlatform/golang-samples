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
	"strconv"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	location                 = "us-central1"
	templateID               = "my-go-test-template"
	deleteTemplateReponse    = "Deleted job template"
	deleteJobReponse         = "Deleted job"
	jobRunningState          = "RUNNING"
	jobSucceededState        = "SUCCEEDED"
	testBucketName           = "cloud-samples-data"
	testBucketDirName        = "media/"
	testVideoFileName        = "ChromeCast.mp4"
	testConcatFileName       = "ForBiggerEscapes.mp4"
	testOverlayImageFileName = "overlay.jpg"
	testCaptionsFileName     = "captions.srt"
	testSubtitlesFileName1   = "subtitles-en.srt"
	testSubtitlesFileName2   = "subtitles-es.srt"
	smallSpriteSheetFileName = "small-sprite-sheet0000000000.jpeg"
	largeSpriteSheetFileName = "large-sprite-sheet0000000000.jpeg"
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

	bucketName := testutil.TestBucket(ctx, t, tc.ProjectID, "golang-samples-transcoder")
	inputURI := "gs://" + bucketName + "/" + testBucketDirName + testVideoFileName
	inputConcatURI := "gs://" + bucketName + "/" + testBucketDirName + testConcatFileName
	inputOverlayImageURI := "gs://" + bucketName + "/" + testBucketDirName + testOverlayImageFileName
	inputCaptionsURI := "gs://" + bucketName + "/" + testBucketDirName + testCaptionsFileName
	inputSubtitles1URI := "gs://" + bucketName + "/" + testBucketDirName + testSubtitlesFileName1
	inputSubtitles2URI := "gs://" + bucketName + "/" + testBucketDirName + testSubtitlesFileName2
	outputURIForPreset := "gs://" + bucketName + "/test-output-preset/"
	outputURIForPresetBatchMode := "gs://" + bucketName + "/test-output-preset-batch-mode/"
	outputURIForTemplate := "gs://" + bucketName + "/test-output-template/"
	outputURIForAdHoc := "gs://" + bucketName + "/test-output-adhoc/"
	outputURIForStaticOverlay := "gs://" + bucketName + "/test-output-static-overlay/"
	outputURIForAnimatedOverlay := "gs://" + bucketName + "/test-output-animated-overlay/"
	outputDirForSetNumberSpritesheet := "test-output-set-number-spritesheet/"
	outputURIForSetNumberSpritesheet := "gs://" + bucketName + "/" + outputDirForSetNumberSpritesheet
	outputDirForPeriodicSpritesheet := "test-output-periodic-spritesheet/"
	outputURIForPeriodicSpritesheet := "gs://" + bucketName + "/" + outputDirForPeriodicSpritesheet
	outputURIForConcat := "gs://" + bucketName + "/test-output-concat/"
	outputURIForEmbeddedCaptions := "gs://" + bucketName + "/test-output-embedded-captions/"
	outputURIForStandaloneCaptions := "gs://" + bucketName + "/test-output-standalone-captions/"

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
	writeTestGCSFiles(t, bucketName)
	t.Logf("\nwriteTestGCSFiles() completed\n")
	testJobFromPreset(t, projectNumber, inputURI, outputURIForPreset)
	t.Logf("\ntestJobFromPreset() completed\n")
	testJobFromPresetBatchMode(t, projectNumber, inputURI, outputURIForPresetBatchMode)
	t.Logf("\ntestJobFromPresetBatchMode() completed\n")
	testJobFromTemplate(t, projectNumber, inputURI, outputURIForTemplate)
	t.Logf("\ntestJobFromTemplate() completed\n")
	testJobFromAdHoc(t, projectNumber, inputURI, outputURIForAdHoc)
	t.Logf("\ntestJobFromAdHoc() completed\n")
	testJobWithStaticOverlay(t, projectNumber, inputURI, inputOverlayImageURI, outputURIForStaticOverlay)
	t.Logf("\ntestJobWithStaticOverlay() completed\n")
	testJobWithAnimatedOverlay(t, projectNumber, inputURI, inputOverlayImageURI, outputURIForAnimatedOverlay)
	t.Logf("\ntestJobWithAnimatedOverlay() completed\n")

	testJobWithSetNumberImagesSpritesheet(t, projectNumber, inputURI, outputURIForSetNumberSpritesheet)
	t.Logf("\ntestJobWithSetNumberImagesSpritesheet() completed\n")
	// Check if the spritesheets exist.
	checkGCSFileExists(t, bucketName, outputDirForSetNumberSpritesheet+smallSpriteSheetFileName)
	checkGCSFileExists(t, bucketName, outputDirForSetNumberSpritesheet+largeSpriteSheetFileName)

	testJobWithPeriodicImagesSpritesheet(t, projectNumber, inputURI, outputURIForPeriodicSpritesheet)
	t.Logf("\ntestJobWithPeriodicImagesSpritesheet() completed\n")
	// Check if the spritesheets exist.
	checkGCSFileExists(t, bucketName, outputDirForPeriodicSpritesheet+smallSpriteSheetFileName)
	checkGCSFileExists(t, bucketName, outputDirForPeriodicSpritesheet+largeSpriteSheetFileName)

	testJobWithConcatenatedInputs(t, projectNumber, inputURI, 0*time.Second, 8*time.Second+100*time.Millisecond, inputConcatURI, 3*time.Second+500*time.Millisecond, 15*time.Second, outputURIForConcat)
	t.Logf("\ntestJobWithConcatenatedInputs() completed\n")

	testJobWithEmbeddedCaptions(t, projectNumber, inputURI, inputCaptionsURI, outputURIForEmbeddedCaptions)
	t.Logf("\ntestJobWithEmbeddedCaptions() completed\n")
	testJobWithStandaloneCaptions(t, projectNumber, inputURI, inputSubtitles1URI, inputSubtitles2URI, outputURIForStandaloneCaptions)
	t.Logf("\ntestJobWithStandaloneCaptions() completed\n")
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

func writeTestGCSFiles(t *testing.T, bucketName string) {
	t.Helper()
	writeTestGCSFile(t, bucketName, testBucketName, testBucketDirName+testVideoFileName)
	writeTestGCSFile(t, bucketName, testBucketName, testBucketDirName+testOverlayImageFileName)
	writeTestGCSFile(t, bucketName, testBucketName, testBucketDirName+testConcatFileName)
	writeTestGCSFile(t, bucketName, testBucketName, testBucketDirName+testCaptionsFileName)
	writeTestGCSFile(t, bucketName, testBucketName, testBucketDirName+testSubtitlesFileName1)
	writeTestGCSFile(t, bucketName, testBucketName, testBucketDirName+testSubtitlesFileName2)
}

// writeTestGCSFile deletes the GCS test bucket and uploads a test video file to it.
func writeTestGCSFile(t *testing.T, dstBucket string, srcBucket string, srcObject string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dstObject := srcObject
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		t.Fatalf("Object(%q).CopierFrom(%q).Run: %v", dstObject, srcObject, err)
	}
}

func checkGCSFileExists(t *testing.T, bucketName string, fileName string) {
	var err error
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		var objAttrs *storage.ObjectAttrs
		objAttrs, err = client.Bucket(bucketName).Object(fileName).Attrs(ctx)
		if err == nil && objAttrs != nil {
			return
		}
		if err == storage.ErrObjectNotExist {
			return
		}
		if err != nil {
			r.Errorf("Error getting bucket attrs: %v", err)
		}
	})
	if err != nil {
		t.Fatal(err)
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
	if err := createJobFromPreset(buf, tc.ProjectID, location, inputURI, outputURIForPreset); err != nil {
		t.Errorf("createJobFromPreset got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobFromPreset got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobFromPresetBatchMode tests major operations on a job created from a
// preset in batch mode. It will wait until the job successfully completes as
// part of the test.
func testJobFromPresetBatchMode(t *testing.T, projectNumber string, inputURI string, outputURIForPresetBatchMode string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobFromPresetBatchMode(buf, tc.ProjectID, location, inputURI, outputURIForPresetBatchMode); err != nil {
		t.Errorf("createJobFromPresetBatchMode got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobFromPresetBatchMode got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobWithStaticOverlay tests major operations on a job created from an ad-hoc configuration that
// includes a static overlay. It will wait until the job successfully completes as part of the test.
func testJobWithStaticOverlay(t *testing.T, projectNumber string, inputURI string, inputOverlayImageURI string, outputURIForStaticOverlay string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobWithStaticOverlay(buf, tc.ProjectID, location, inputURI, inputOverlayImageURI, outputURIForStaticOverlay); err != nil {
		t.Errorf("createJobWithStaticOverlay got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobWithStaticOverlay got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobWithAnimatedOverlay tests major operations on a job created from an ad-hoc configuration that
// includes an animated overlay. It will wait until the job successfully completes as part of the test.
func testJobWithAnimatedOverlay(t *testing.T, projectNumber string, inputURI string, inputOverlayImageURI string, outputURIForAnimatedOverlay string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobWithAnimatedOverlay(buf, tc.ProjectID, location, inputURI, inputOverlayImageURI, outputURIForAnimatedOverlay); err != nil {
		t.Errorf("testJobWithAnimatedOverlay got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("testJobWithAnimatedOverlay got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobWithSetNumberImagesSpritesheet tests major operations on a job created from an ad-hoc configuration that
// generates two spritesheets. It will wait until the job successfully completes as part of the test.
func testJobWithSetNumberImagesSpritesheet(t *testing.T, projectNumber string, inputURI string, outputURIForSetNumberSpritesheet string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobWithSetNumberImagesSpritesheet(buf, tc.ProjectID, location, inputURI, outputURIForSetNumberSpritesheet); err != nil {
		t.Errorf("createJobWithSetNumberImagesSpritesheet got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobWithSetNumberImagesSpritesheet got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobWithPeriodicImagesSpritesheet tests major operations on a job created from an ad-hoc configuration that
// generates two spritesheets. It will wait until the job successfully completes as part of the test.
func testJobWithPeriodicImagesSpritesheet(t *testing.T, projectNumber string, inputURI string, outputURIForPeriodicSpritesheet string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	testutil.Retry(t, 3, 2*time.Second, func(r *testutil.R) {
		buf.Reset()
		if err := createJobWithPeriodicImagesSpritesheet(buf, tc.ProjectID, location, inputURI, outputURIForPeriodicSpritesheet); err != nil {
			r.Errorf("createJobWithPeriodicImagesSpritesheet got err: %v", err)
		}
		got := buf.String()

		if !strings.Contains(got, jobName) {
			r.Errorf("createJobWithPeriodicImagesSpritesheet got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
		}
		strSlice := strings.Split(got, "/")
		jobID = strSlice[len(strSlice)-1]
	})
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobWithConcatenatedInputs tests major operations on a job created from an ad-hoc configuration that
// concatenates two inputs videos. It will wait until the job successfully completes as part of the test.
func testJobWithConcatenatedInputs(t *testing.T, projectNumber string, input1URI string, startTimeInput1 time.Duration, endTimeInput1 time.Duration, input2URI string, startTimeInput2 time.Duration, endTimeInput2 time.Duration, outputURIForConcat string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobWithConcatenatedInputs(buf, tc.ProjectID, location, input1URI, startTimeInput1, endTimeInput1, input2URI, startTimeInput2, endTimeInput2, outputURIForConcat); err != nil {
		t.Errorf("testJobWithConcatenatedInputs got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("testJobWithConcatenatedInputs got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobWithEmbeddedCaptions tests major operations on a job created from an ad-hoc configuration that
// embeds captions in the output video. It will wait until the job successfully completes as part of the test.
func testJobWithEmbeddedCaptions(t *testing.T, projectNumber string, inputVideoURI string, inputCaptionsURI string, outputURIForEmbeddedCaptions string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobWithEmbeddedCaptions(buf, tc.ProjectID, location, inputVideoURI, inputCaptionsURI, outputURIForEmbeddedCaptions); err != nil {
		t.Errorf("createJobWithEmbeddedCaptions got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobWithEmbeddedCaptions got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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

// testJobWithStandaloneCaptions tests major operations on a job created from an ad-hoc configuration that
// can use captions from a standalone file. It will wait until the job successfully completes as part of the test.
func testJobWithStandaloneCaptions(t *testing.T, projectNumber string, inputVideoURI string, inputSubtitles1URI string, inputSubtitles2URI string, outputURIForStandaloneCaptions string) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	jobID := ""

	// Create the job.
	jobName := fmt.Sprintf("projects/%s/locations/%s/jobs/", projectNumber, location)
	if err := createJobWithStandaloneCaptions(buf, tc.ProjectID, location, inputVideoURI, inputSubtitles1URI, inputSubtitles2URI, outputURIForStandaloneCaptions); err != nil {
		t.Errorf("createJobWithStandaloneCaptions got err: %v", err)
	}
	got := buf.String()

	if !strings.Contains(got, jobName) {
		t.Errorf("createJobWithStandaloneCaptions got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobName)
	}
	strSlice := strings.Split(got, "/")
	jobID = strSlice[len(strSlice)-1]
	buf.Reset()

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
	buf.Reset()

	// Get the job state, which should be succeeded. If the job is still running on the last attempt, pass the test.
	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := getJobState(buf, tc.ProjectID, location, jobID); err != nil {
			r.Errorf("getJobState got err: %v", err)
		}
		got := buf.String()

		if r.Attempt == 3 {
			if !strings.Contains(got, jobSucceededState) && !strings.Contains(got, jobRunningState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v or %v\n----\n", got, jobSucceededState, jobRunningState)
			}
		} else {
			if !strings.Contains(got, jobSucceededState) {
				r.Errorf("getJobState got\n----\n%v\n----\nWant to contain:\n----\n%v\n----\n", got, jobSucceededState)
			}
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
