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

package inspect

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

const (
	topicName        = "dlp-inspect-test-topic-"
	subscriptionName = "dlp-inspect-test-sub-"

	ssnFileName = "fake_ssn.txt"
	bucketName  = "golang-samples-dlp-test2"

	jobTriggerIdPrefix                      = "dlp-job-trigger-unit-test-case-12345678"
	dataSetIDForHybridJob                   = "dlp_test_dataset"
	tableIDForHybridJob                     = "dlp_inspect_test_table_table_id"
	inspectsGCSTestFileName                 = "test.txt"
	filePathToUpload                        = "./testdata/test.txt"
	dirPathForInspectGCSSendToScc           = "dlp-go-lang-test-for-inspect-gcs-send-to-scc/"
	bucketnameForInspectGCSFileWithSampling = "dlp-job-go-lang-test-inspect-gcs-file-with-sampling"
)

func TestInspectDatastore(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	writeTestDatastoreFiles(t, tc.ProjectID)
	tests := []struct {
		kind string
		want string
	}{
		{
			kind: "SSNTask",
			want: "Created job",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.kind, func(t *testing.T) {
			t.Parallel()
			testutil.Retry(t, 5, 15*time.Second, func(r *testutil.R) {
				u := uuid.New().String()[:8]
				buf := new(bytes.Buffer)
				if err := inspectDatastore(buf, tc.ProjectID, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, topicName+u, subscriptionName+u, tc.ProjectID, "", test.kind); err != nil {
					r.Errorf("inspectDatastore(%s) got err: %v", test.kind, err)
					return
				}
				if got := buf.String(); !strings.Contains(got, test.want) {
					r.Errorf("inspectDatastore(%s) = %q, want %q substring", test.kind, got, test.want)
				}
			})
		})
	}
}

type SSNTask struct {
	Description string
}

func writeTestDatastoreFiles(t *testing.T, projectID string) {
	t.Helper()
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("datastore.NewClient: %v", err)
	}
	kind := "SSNTask"
	name := "ssntask1"
	ssnKey := datastore.NameKey(kind, name, nil)
	task := SSNTask{
		Description: "My SSN is 111222333",
	}
	if _, err := client.Put(ctx, ssnKey, &task); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}
}

func TestInspectGCS(t *testing.T) {
	tc := testutil.SystemTest(t)
	writeTestGCSFiles(t, tc.ProjectID)
	tests := []struct {
		fileName string
		want     string
	}{
		{
			fileName: ssnFileName,
			want:     "Created job",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.fileName, func(t *testing.T) {
			t.Parallel()
			testutil.Retry(t, 5, 15*time.Second, func(r *testutil.R) {
				u := uuid.New().String()[:8]
				buf := new(bytes.Buffer)
				if err := inspectGCSFile(buf, tc.ProjectID, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, topicName+u, subscriptionName+u, bucketName, test.fileName); err != nil {
					r.Errorf("inspectGCSFile(%s) got err: %v", test.fileName, err)
					return
				}
				if got := buf.String(); !strings.Contains(got, test.want) {
					r.Errorf("inspectGCSFile(%s) = %q, want %q substring", test.fileName, got, test.want)
				}
			})
		})
	}
}

func writeTestGCSFiles(t *testing.T, projectID string) {
	t.Helper()
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		switch err {
		case storage.ErrObjectNotExist:
			if err := bucket.Create(ctx, projectID, nil); err != nil {
				t.Fatalf("bucket.Create: %v", err)
			}
		default:
			t.Fatalf("error getting bucket attrs: %v", err)
		}
	}
	if err := writeObject(ctx, bucket, ssnFileName, "My SSN is 111222333"); err != nil {
		t.Fatalf("writeObject: %v", err)
	}
}

func writeObject(ctx context.Context, bucket *storage.BucketHandle, fileName, content string) error {
	obj := bucket.Object(fileName)
	_, err := obj.Attrs(ctx)
	if err != nil {
		switch err {
		case storage.ErrObjectNotExist:
			w := obj.NewWriter(ctx)
			w.Write([]byte(content))
			if err := w.Close(); err != nil {
				return err
			}
		default:
			return err
		}
	}
	return nil
}

func TestInspectString(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	if err := inspectString(buf, tc.ProjectID, "I'm Gary and my email is gary@example.com"); err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectString got %q, want %q", got, want)
	}
}

func TestInspectTextFile(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	if err := inspectTextFile(buf, tc.ProjectID, "testdata/test.txt"); err != nil {
		t.Errorf("TestInspectTextFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("inspectTextFile got %q, want %q", got, want)
	}
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectTextFile got %q, want %q", got, want)
	}
}

type Item struct {
	Description string
}

const (
	harmfulTable = "harmful"
	bqDatasetID  = "golang_samples_dlp"
)

func mustCreateBigqueryTestFiles(t *testing.T, projectID, datasetID string) {
	t.Helper()

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()
	d := client.Dataset(datasetID)
	if _, err := d.Metadata(ctx); err != nil {
		if err := d.Create(ctx, &bigquery.DatasetMetadata{}); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}
	schema, err := bigquery.InferSchema(Item{})
	if err != nil {
		t.Fatalf("InferSchema: %v", err)
	}
	if err := uploadBigQuery(ctx, d, schema, harmfulTable, "My SSN is 111222333"); err != nil {
		t.Fatalf("uploadBigQuery: %v", err)
	}
}

func uploadBigQuery(ctx context.Context, d *bigquery.Dataset, schema bigquery.Schema, table, content string) error {
	t := d.Table(table)
	if _, err := t.Metadata(ctx); err == nil {
		return nil
	}
	if err := t.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
		return err
	}
	source := bigquery.NewReaderSource(strings.NewReader(content))
	l := t.LoaderFrom(source)
	job, err := l.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	return status.Err()
}

func TestInspectBigquery(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	mustCreateBigqueryTestFiles(t, tc.ProjectID, bqDatasetID)

	tests := []struct {
		table string
		want  string
	}{
		{
			table: harmfulTable,
			want:  "Created job",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.table, func(t *testing.T) {
			t.Parallel()
			u := uuid.New().String()[:8]
			buf := new(bytes.Buffer)
			if err := inspectBigquery(buf, tc.ProjectID, []string{"US_SOCIAL_SECURITY_NUMBER"}, []string{}, []string{}, topicName+u, subscriptionName+u, tc.ProjectID, bqDatasetID, test.table); err != nil {
				t.Errorf("inspectBigquery(%s) got err: %v", test.table, err)
			}
			if got := buf.String(); !strings.Contains(got, test.want) {
				t.Errorf("inspectBigquery(%s) = %q, want %q substring", test.table, got, test.want)
			}
		})
	}
}

func TestInspectTable(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	if err := inspectTable(&buf, tc.ProjectID); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Infotype Name: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("InspectTable got %q, want %q", got, want)
	}
	if want := "Likelihood: VERY_LIKELY"; !strings.Contains(got, want) {
		t.Errorf("InspectTable got %q, want %q", got, want)
	}
}

func TestInspectStringWithExclusionRegex(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	if err := inspectStringWithExclusionRegex(&buf, tc.ProjectID, "Some email addresses: gary@example.com, bob@example.org", ".+@example.com"); err != nil {
		t.Errorf("inspectStringWithExclusionRegex: %v", err)
	}

	got := buf.String()

	if want := "Quote: bob@example.org"; !strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionRegex got %q, want %q", got, want)
	}
	if want := "Quote: gary@example.com"; strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionRegex got %q, want %q", got, want)
	}
}

func TestInspectStringCustomExcludingSubstring(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringCustomExcludingSubstring(&buf, tc.ProjectID, "Name: Doe, John. Name: Example, Jimmy", "[A-Z][a-z]{1,15}, [A-Z][a-z]{1,15}", []string{"Jimmy"}); err != nil {
		t.Fatal(err)
	}

	got := buf.String()

	if want := "Infotype Name: CUSTOM_NAME_DETECTOR"; !strings.Contains(got, want) {
		t.Errorf("inspectStringCustomExcludingSubstring got %q, want %q", got, want)
	}
	if want := "Quote: Doe, John"; !strings.Contains(got, want) {
		t.Errorf("inspectStringCustomExcludingSubstring got %q, want %q", got, want)
	}
	if want := "Jimmy"; strings.Contains(got, want) {
		t.Errorf("inspectStringCustomExcludingSubstring got %q, want %q", got, want)
	}
}

func TestInspectStringMultipleRules(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringMultipleRules(&buf, tc.ProjectID, "patient: Jane Doe"); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Infotype Name: PERSON_NAME"; !strings.Contains(got, want) {
		t.Errorf("inspectStringMultipleRules got %q, want %q", got, want)
	}
}

func TestInspectWithHotWordRules(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectWithHotWordRules(&buf, tc.ProjectID, "Patient's MRN 444-5-22222 and just a number 333-2-33333"); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "InfoType Name: C_MRN"; !strings.Contains(got, want) {
		t.Errorf("inspectWithHotWordRules got %q, want %q", got, want)
	}
	if want := "Findings: 2"; !strings.Contains(got, want) {
		t.Errorf("inspectWithHotWordRules got %q, want %q", got, want)
	}
}

func TestInspectPhoneNumber(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectPhoneNumber(&buf, tc.ProjectID, "I'm Gary and my phone number is (415) 555-0890"); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("inspectPhoneNumber got %q, want %q", got, want)
	}
}

func TestInspectStringWithoutOverlap(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringWithoutOverlap(&buf, tc.ProjectID, "example.com is a domain, james@example.org is an email."); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Infotype Name: DOMAIN_NAME"; !strings.Contains(got, want) {
		t.Errorf("inspectStringWithoutOverlap got %q, want %q", got, want)
	}
	if want := "Quote: example.com"; !strings.Contains(got, want) {
		t.Errorf("inspectStringWithoutOverlap got %q, want %q", got, want)
	}
	if want := "Quote: example.org"; strings.Contains(got, want) {
		t.Errorf("inspectStringWithoutOverlap got %q, want %q", got, want)
	}
}

func TestInspectStringCustomHotWord(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringCustomHotWord(&buf, tc.ProjectID, "patient name: John Doe", "patient", "PERSON_NAME"); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Infotype Name: PERSON_NAME"; !strings.Contains(got, want) {
		t.Errorf("inspectStringCustomHotWord got %q, want %q", got, want)
	}
}

func TestInspectStringWithExclusionDictSubstring(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringWithExclusionDictSubstring(&buf, tc.ProjectID, "Some email addresses: gary@example.com, TEST@example.com", []string{"TEST"}); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Infotype Name: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionDictSubstring got %q, want %q", got, want)
	}
	if want := "Infotype Name: DOMAIN_NAME"; !strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionDictSubstring got %q, want %q", got, want)
	}
	if want := "Quote: TEST"; strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionDictSubstring got %q, want %q", got, want)
	}
}

func TestInspectStringOmitOverlap(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringOmitOverlap(&buf, tc.ProjectID, "gary@example.com"); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Infotype Name: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectStringOmitOverlap got %q, want %q", got, want)
	}
}

func TestInspectStringCustomOmitOverlap(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectStringCustomHotWord(&buf, tc.ProjectID, "patient name: John Doe", "patient", "PERSON_NAME"); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Infotype Name: PERSON_NAME"; !strings.Contains(got, want) {
		t.Errorf("inspectStringCustomOmitOverlap got %q, want %q", got, want)
	}

	if want := "Quote: John Doe"; !strings.Contains(got, want) {
		t.Errorf("inspectStringCustomOmitOverlap got %q, want %q", got, want)
	}
	if want := "Quote: Larry Page"; strings.Contains(got, want) {
		t.Errorf("inspectStringCustomOmitOverlap got %q, want %q", got, want)
	}
}

func TestInspectWithCustomRegex(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	if err := inspectWithCustomRegex(&buf, tc.ProjectID, "Patients MRN 444-5-22222", "[1-9]{3}-[1-9]{1}-[1-9]{5}", "C_MRN"); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Infotype Name: C_MRN"; !strings.Contains(got, want) {
		t.Errorf("inspectWithCustomRegex got %q, want %q", got, want)
	}
	if want := "Likelihood: POSSIBLE"; !strings.Contains(got, want) {
		t.Errorf("inspectWithCustomRegex got %q, want %q", got, want)
	}
}

func TestInspectStringWithExclusionDictionary(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	if err := inspectStringWithExclusionDictionary(&buf, tc.ProjectID, "Some email addresses: gary@example.com, example@example.com", []string{"example@example.com"}); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Infotype Name: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectStringWithExclusionDictionary got %q, want %q", got, want)
	}
}

func TestInspectImageFile(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	pathToImage := "testdata/test.png"
	if err := inspectImageFile(&buf, tc.ProjectID, pathToImage); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("TestInspectImageFile got %q, want %q", got, want)
	}
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("TestInspectImageFile got %q, want %q", got, want)
	}
}

func TestInspectImageFileAllInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)
	inputPath := "testdata/image.jpg"

	var buf bytes.Buffer
	if err := inspectImageFileAllInfoTypes(&buf, tc.ProjectID, inputPath); err != nil {
		t.Errorf("inspectImageFileAllInfoTypes: %v", err)
	}
	got := buf.String()
	if want := "Info type: DATE"; !strings.Contains(got, want) {
		t.Errorf("inspectImageFileAllInfoTypes got %q, want %q", got, want)
	}
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("inspectImageFileAllInfoTypes got %q, want %q", got, want)
	}
	if want := "Info type: US_SOCIAL_SECURITY_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("inspectImageFileAllInfoTypes got %q, want %q", got, want)
	}
}

func TestInspectImageFileListedInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	pathToImage := "testdata/sensitive-data-image.jpg"

	if err := inspectImageFileListedInfoTypes(&buf, tc.ProjectID, pathToImage); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("inspectImageFileListedInfoTypes got %q, want %q", got, want)
	}
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("inspectImageFileListedInfoTypes got %q, want %q", got, want)
	}
	if want := "Info type: US_SOCIAL_SECURITY_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("inspectImageFileListedInfoTypes got %q, want %q", got, want)
	}
}

func TestInspectGcsFileWithSampling(t *testing.T) {
	tc := testutil.SystemTest(t)
	topicID := "go-lang-dlp-test-bigquery-with-sampling-topic"
	subscriptionID := "go-lang-dlp-test-bigquery-with-sampling-subscription"
	ctx := context.Background()
	sc, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer sc.Close()

	bucketnameForInspectGCSFileWithSampling, err := testutil.CreateTestBucket(ctx, t, sc, tc.ProjectID, "dlp-test-inspect-prefix")
	if err != nil {
		t.Fatal(err)
	}
	GCSUri := "gs://" + bucketnameForInspectGCSFileWithSampling + "/"

	var buf bytes.Buffer
	if err := inspectGcsFileWithSampling(&buf, tc.ProjectID, GCSUri, topicID, subscriptionID); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Job Created"; !strings.Contains(got, want) {
		t.Errorf("inspectGcsFileWithSampling got %q, want %q", got, want)
	}
	err = testutil.DeleteBucketIfExists(ctx, sc, bucketnameForInspectGCSFileWithSampling)
	if err != nil {
		t.Fatal(err)
	}

}

func TestInspectBigQueryTableWithSampling(t *testing.T) {
	tc := testutil.SystemTest(t)

	topicID := "go-lang-dlp-test-bigquery-with-sampling-topic"
	subscriptionID := "go-lang-dlp-test-bigquery-with-sampling-subscription"

	var buf bytes.Buffer
	if err := inspectBigQueryTableWithSampling(&buf, tc.ProjectID, topicID, subscriptionID); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Job Created"; !strings.Contains(got, want) {
		t.Errorf("InspectBigQueryTableWithSampling got %q, want %q", got, want)
	}
	if want := "Found"; !strings.Contains(got, want) {
		t.Errorf("InspectBigQueryTableWithSampling got %q, want %q", got, want)
	}

}

func TestInspectAugmentInfoTypes(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	textToInspect := "The patient's name is Quasimodo"
	wordList := []string{"quasimodo"}

	if err := inspectAugmentInfoTypes(&buf, tc.ProjectID, textToInspect, wordList); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Qoute: Quasimodo"; !strings.Contains(got, want) {
		t.Errorf("TestInspectAugmentInfoTypes got %q, want %q", got, want)
	}
	if want := "Info type: PERSON_NAME"; !strings.Contains(got, want) {
		t.Errorf("TestInspectAugmentInfoTypes got %q, want %q", got, want)
	}
}

func TestInspectTableWithCustomHotword(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	hotwordRegexPattern := "(Fake Social Security Number)"
	if err := inspectTableWithCustomHotword(&buf, tc.ProjectID, hotwordRegexPattern); err != nil {
		t.Fatal(err)
	}
	got := buf.String()

	if want := "Quote: 222-22-2222"; !strings.Contains(got, want) {
		t.Errorf("TestInspectTableWithCustomHotword got %q, want %q", got, want)
	}
	if want := "Infotype Name: US_SOCIAL_SECURITY_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("TestInspectTableWithCustomHotword got %q, want %q", got, want)
	}
	if want := "Quote: 111-11-1111"; strings.Contains(got, want) {
		t.Errorf("TestInspectTableWithCustomHotword got %q, want %q", got, want)
	}
}

func TestInspectDataStoreSendToScc(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	u := uuid.New().String()[:8]
	datastoreNamespace := fmt.Sprint("golang-samples" + u)
	datastoreKind := "task"

	if err := inspectDataStoreSendToScc(&buf, tc.ProjectID, datastoreNamespace, datastoreKind); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Job created successfully:"; !strings.Contains(got, want) {
		t.Errorf("InspectBigQuerySendToScc got %q, want %q", got, want)
	}
}

func TestInspectGCSFileSendToScc(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	ctx := context.Background()
	sc, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer sc.Close()

	// Creates a bucket using a function available in testutil.
	bucketNameForInspectGCSSendToScc, err := testutil.CreateTestBucket(ctx, t, sc, tc.ProjectID, "dlp-test-inspect-prefix")
	if err != nil {
		t.Fatal(err)
	}

	// Uploads a file on created bucket.
	filePathtoGCS(t, tc.ProjectID, bucketNameForInspectGCSSendToScc, dirPathForInspectGCSSendToScc)

	gcsPath := fmt.Sprint("gs://" + bucketNameForInspectGCSSendToScc + "/" + dirPathForInspectGCSSendToScc + "/test.txt")

	if err := inspectGCSFileSendToScc(&buf, tc.ProjectID, gcsPath); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Job created successfully:"; !strings.Contains(got, want) {
		t.Errorf("TestInspectGCSFileSendToScc got %q, want %q", got, want)
	}

	// Delete a bucket that has just been created.
	err = testutil.DeleteBucketIfExists(ctx, sc, bucketNameForInspectGCSSendToScc)
	if err != nil {
		t.Fatal(err)
	}
}

// filePathtoGCS uploads a file test.txt in given path from the testdata directory.
func filePathtoGCS(t *testing.T, projectID, bucketNameForInspectGCSSendToScc, dirPathForInspectGCSSendToScc string) error {
	t.Helper()
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Check if the bucket already exists.
	bucketExists := false
	_, err = client.Bucket(bucketNameForInspectGCSSendToScc).Attrs(ctx)
	if err == nil {
		bucketExists = true
	}

	// If the bucket doesn't exist, create it.
	if !bucketExists {
		if err := client.Bucket(bucketNameForInspectGCSSendToScc).Create(ctx, projectID, &storage.BucketAttrs{
			StorageClass: "STANDARD",
			Location:     "us-central1",
		}); err != nil {
			return err
		}
		fmt.Printf("Bucket '%s' created successfully.\n", bucketNameForInspectGCSSendToScc)
	} else {
		fmt.Printf("Bucket '%s' already exists.\n", bucketNameForInspectGCSSendToScc)
	}

	// Check if the directory already exists in the bucket.
	dirExists := false
	query := &storage.Query{Prefix: dirPathForInspectGCSSendToScc}
	it := client.Bucket(bucketNameForInspectGCSSendToScc).Objects(ctx, query)
	_, err = it.Next()
	if err == nil {
		dirExists = true
	}

	// If the directory doesn't exist, create it.
	if !dirExists {
		obj := client.Bucket(bucketNameForInspectGCSSendToScc).Object(dirPathForInspectGCSSendToScc)
		if _, err := obj.NewWriter(ctx).Write([]byte("")); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		fmt.Printf("Directory '%s' created successfully in bucket '%s'.\n", dirPathForInspectGCSSendToScc, bucketNameForInspectGCSSendToScc)
	} else {
		fmt.Printf("Directory '%s' already exists in bucket '%s'.\n", dirPathForInspectGCSSendToScc, bucketNameForInspectGCSSendToScc)
	}

	// file upload code

	// Open local file.
	file, err := ioutil.ReadFile(filePathToUpload)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return err
	}

	// Get a reference to the bucket
	bucket := client.Bucket(bucketNameForInspectGCSSendToScc)

	// Upload the file
	object := bucket.Object(inspectsGCSTestFileName)
	writer := object.NewWriter(ctx)
	_, err = writer.Write(file)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close writer: %v", err)
		return err
	}
	fmt.Printf("File uploaded successfully: %v\n", inspectsGCSTestFileName)

	// Check if the file exists in the bucket
	_, err = bucket.Object(inspectsGCSTestFileName).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			fmt.Printf("File %v does not exist in bucket %v\n", inspectsGCSTestFileName, bucketNameForInspectGCSSendToScc)
		} else {
			log.Fatalf("Failed to check file existence: %v", err)
		}
	} else {
		fmt.Printf("File %v exists in bucket %v\n", inspectsGCSTestFileName, bucketNameForInspectGCSSendToScc)
	}

	fmt.Println("filePathtoGCS function is executed-------")
	return nil
}

var projectID, jobTriggerForInspectSample string

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	projectID = tc.ProjectID
	if !ok {
		log.Fatal("couldn't initialize test")
		return
	}
	xyz, err := createJobTriggerForInspectDataToHybridJobTrigger(tc.ProjectID)
	jobTriggerForInspectSample = xyz
	if err != nil {
		log.Fatal("couldn't initialize test")
		return
	}
	m.Run()
	deleteActiveJob(tc.ProjectID, jobTriggerForInspectSample)
	deleteJobTriggerForInspectDataToHybridJobTrigger(tc.ProjectID, jobTriggerForInspectSample)
}

func TestInspectDataToHybridJobTrigger(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	trigger := jobTriggerForInspectSample
	fmt.Print("Name:" + trigger)
	if err := inspectDataToHybridJobTrigger(&buf, tc.ProjectID, "My email is test@example.org and my name is Gary.", trigger); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "successfully inspected data using hybrid job trigger"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}
	if want := "Findings"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}
	if want := "Job State: ACTIVE"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}
	if want := "EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}
	if want := "PERSON_NAME"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}

}

func deleteActiveJob(project, trigger string) error {

	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	req := &dlppb.ListDlpJobsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", project),
		Filter: fmt.Sprintf("trigger_name=%s", trigger),
	}

	it := client.ListDlpJobs(ctx, req)
	var jobIds []string
	for {
		j, err := it.Next()
		jobIds = append(jobIds, j.GetName())
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Next: %v", err)
		}
		fmt.Printf("Job %v status: %v\n", j.GetName(), j.GetState())
	}
	for _, v := range jobIds {
		req := &dlppb.DeleteDlpJobRequest{
			Name: v,
		}
		if err = client.DeleteDlpJob(ctx, req); err != nil {
			fmt.Printf("DeleteDlpJob: %v", err)
			return err
		}
		fmt.Printf("\nSuccessfully deleted job %v\n", v)
	}
	fmt.Print("Deleted Job Successfully !!!")
	return nil
}

// helpers for inspect hybrid job
func createJobTriggerForInspectDataToHybridJobTrigger(projectID string) (string, error) {

	log.Printf("[START] createJobTriggerForInspectDataToHybridJobTrigger: projectID %v and ", projectID)
	// Set up the client.
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Define the job trigger.
	hybridOptions := &dlppb.HybridOptions{
		Labels: map[string]string{
			"env": "prod",
		},
	}

	storageConfig := &dlppb.StorageConfig_HybridOptions{
		HybridOptions: hybridOptions,
	}
	infoTypes := []*dlppb.InfoType{
		{Name: "PERSON_NAME"},
		{Name: "EMAIL_ADDRESS"},
	}

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: infoTypes,
	}

	inspectJobConfig := &dlppb.InspectJobConfig{
		StorageConfig: &dlppb.StorageConfig{
			Type: storageConfig,
		},
		InspectConfig: inspectConfig,
	}

	trigger := &dlppb.JobTrigger_Trigger{
		Trigger: &dlppb.JobTrigger_Trigger_Manual{},
	}

	jobTrigger := &dlppb.JobTrigger{
		Triggers: []*dlppb.JobTrigger_Trigger{
			trigger,
		},
		Job: &dlppb.JobTrigger_InspectJob{
			InspectJob: inspectJobConfig,
		},
	}

	u := uuid.New().String()[:8]
	createDlpJobRequest := &dlppb.CreateJobTriggerRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/global", projectID),
		JobTrigger: jobTrigger,
		TriggerId:  jobTriggerIdPrefix + u,
	}

	resp, err := client.CreateJobTrigger(ctx, createDlpJobRequest)
	if err != nil {
		return "", err
	}
	log.Printf("[END] createJobTriggerForInspectDataToHybridJobTrigger: trigger.Name %v", resp.Name)
	return resp.Name, nil
}

func deleteJobTriggerForInspectDataToHybridJobTrigger(projectID, jobTriggerName string) error {

	log.Printf("\n[START] deleteJobTriggerForInspectDataToHybridJobTrigger")
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &dlppb.DeleteJobTriggerRequest{
		Name: jobTriggerName,
	}

	err = client.DeleteJobTrigger(ctx, req)
	if err != nil {
		return err
	}
	log.Print("[END] deleteJobTriggerForInspectDataToHybridJobTrigger")
	return nil
}
