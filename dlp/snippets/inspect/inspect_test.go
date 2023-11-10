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

func createBigQueryDataSetId(projectID string) error {

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	meta := &bigquery.DatasetMetadata{
		Location: "US", // See https://cloud.google.com/bigquery/docs/locations
	}

	if err := client.Dataset(dataSetID).Create(ctx, meta); err != nil {
		return err
	}

	return nil
}

func createTableInsideDataset(projectID, dataSetID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	sampleSchema := bigquery.Schema{
		{Name: "user_id", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
		{Name: "title", Type: bigquery.StringFieldType},
		{Name: "score", Type: bigquery.StringFieldType},
	}

	metaData := &bigquery.TableMetadata{
		Schema:         sampleSchema,
		ExpirationTime: time.Now().AddDate(1, 0, 0), // Table will be automatically deleted in 1 year.
	}

	tableRef := client.Dataset(dataSetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		log.Printf("[INFO] createBigQueryDataSetId Error while table creation: %v", err)
		return err
	}

	duration := time.Duration(90) * time.Second
	time.Sleep(duration)

	inserter := client.Dataset(dataSetID).Table(tableID).Inserter()
	items := []*BigQueryTableItem{
		// Item implements the ValueSaver interface.
		{UserId: "602-61-8588", Age: 32, Title: "Biostatistician III", Score: "A"},
		{UserId: "618-96-2322", Age: 69, Title: "Programmer I", Score: "C"},
		{UserId: "618-96-2322", Age: 69, Title: "Executive Secretary", Score: "C"},
	}
	if err := inserter.Put(ctx, items); err != nil {
		return err
	}

	return nil
}

type BigQueryTableItem struct {
	UserId string
	Age    int
	Title  string
	Score  string
}

func (i *BigQueryTableItem) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"user_id": i.UserId,
		"age":     i.Age,
		"title":   i.Title,
		"score":   i.Score,
	}, bigquery.NoDedupeID, nil
}

func deleteBigQueryAssets(projectID string) error {

	log.Printf("[START] deleteBigQueryAssets: projectID %v and ", projectID)
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	log.Printf("[INFO] deleteBigQueryAssets: delete dataset err %v", err)

	if err := client.Dataset("dlp_test_dataset").DeleteWithContents(ctx); err != nil {
		log.Printf("[INFO] deleteBigQueryAssets: delete dataset err %v", err)
		return err
	}

	duration := time.Duration(30) * time.Second
	time.Sleep(duration)

	log.Printf("[END] deleteBigQueryAssets:")
	return nil
}

func deleteJob(projectID, jobName string) error {
	ctx := context.Background()

	log.Printf("[START] deleteJob: projectID %v", projectID)
	// delete job
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Printf("[INFO] deleteJob:: error %v", err)
		return err
	}
	log.Printf("[INFO] deleteJob:: error %v", err)

	req := &dlppb.DeleteDlpJobRequest{
		Name: jobName,
	}
	for {
		ct, cancel := context.WithTimeout(ctx, 300000)
		defer cancel()
		abc, err := client.GetDlpJob(ct, &dlppb.GetDlpJobRequest{
			Name: jobName,
		})
		if err != nil {
			log.Printf("[INFO] deleteJob:: error %v", err)
			return err
		}
		if abc.State == dlppb.DlpJob_DONE {
			log.Printf("[INFO] deleteJob:: job done")
			break
		} else if abc.State == dlppb.DlpJob_FAILED {
			log.Printf("[INFO] deleteJob:: job failed")
			return err
		} else {
			log.Printf("[INFO] deleteJob:: job continue")
			continue
		}
	}
	err = client.DeleteDlpJob(ctx, req)
	if err != nil {
		log.Printf("[INFO] deleteJob:: error %v", err)
		return err
	}

	log.Printf("[END] deleteJob")
	return nil
}

var (
	projectID                  string
	jobTriggerForInspectSample string
	bucketExpiryAge            = time.Minute * 2
	testPrefix                 = "dlp-test-inspect-prefix"
)

func createStoredInfoTypeForTesting(t *testing.T, projectID, outputPath string) (string, error) {
	t.Helper()
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()
	u := uuid.New().String()[:8]
	displayName := "stored-info-type-for-inspect-test" + u
	description := "Dictionary of GitHub usernames used in commits"

	cloudStoragePath := &dlppb.CloudStoragePath{
		Path: outputPath,
	}

	bigQueryField := &dlppb.BigQueryField{
		Table: &dlppb.BigQueryTable{
			ProjectId: "bigquery-public-data",
			DatasetId: "samples",
			TableId:   "github_nested",
		},
		Field: &dlppb.FieldId{
			Name: "actor",
		},
	}

	largeCustomDictionaryConfig := &dlppb.LargeCustomDictionaryConfig{
		OutputPath: cloudStoragePath,
		Source: &dlppb.LargeCustomDictionaryConfig_BigQueryField{
			BigQueryField: bigQueryField,
		},
	}

	storedInfoTypeConfig := &dlppb.StoredInfoTypeConfig{
		DisplayName: displayName,
		Description: description,
		Type: &dlppb.StoredInfoTypeConfig_LargeCustomDictionary{
			LargeCustomDictionary: largeCustomDictionaryConfig,
		},
	}

	req := &dlppb.CreateStoredInfoTypeRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		Config:           storedInfoTypeConfig,
		StoredInfoTypeId: "go-sample-test-stored-infoType" + u,
	}
	resp, err := client.CreateStoredInfoType(ctx, req)
	if err != nil {
		return "nil", err
	}

	return resp.Name, nil
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
		log.Printf("[INFO] [filePathtoGCS] Bucket '%s' created successfully.\n", bucketNameForInspectGCSSendToScc)
	} else {
		log.Printf("[INFO] [filePathtoGCS] Bucket '%s' already exists.\n", bucketNameForInspectGCSSendToScc)
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
		log.Printf("[INFO] [filePathtoGCS] Directory '%s' created successfully in bucket '%s'.\n", dirPathForInspectGCSSendToScc, bucketNameForInspectGCSSendToScc)
	} else {
		log.Printf("[INFO] [filePathtoGCS] Directory '%s' already exists in bucket '%s'.\n", dirPathForInspectGCSSendToScc, bucketNameForInspectGCSSendToScc)
	}

	// file upload code

	// Open local file.
	file, err := ioutil.ReadFile(filePathToUpload)
	if err != nil {
		log.Fatalf("[INFO] [filePathtoGCS] Failed to read file: %v", err)
		return err
	}

	// Get a reference to the bucket
	bucket := client.Bucket(bucketNameForInspectGCSSendToScc)

	// Upload the file
	object := bucket.Object(inspectsGCSTestFileName)
	writer := object.NewWriter(ctx)
	_, err = writer.Write(file)
	if err != nil {
		log.Fatalf("[INFO] [filePathtoGCS] Failed to write file: %v", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("[INFO] [filePathtoGCS] Failed to close writer: %v", err)
		return err
	}
	log.Printf("[INFO] [filePathtoGCS] File uploaded successfully: %v\n", inspectsGCSTestFileName)

	// Check if the file exists in the bucket
	_, err = bucket.Object(inspectsGCSTestFileName).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			log.Printf("[INFO] [filePathtoGCS] File %v does not exist in bucket %v\n", inspectsGCSTestFileName, bucketNameForInspectGCSSendToScc)
		} else {
			log.Fatalf("[INFO] [filePathtoGCS] Failed to check file existence: %v", err)
		}
	} else {
		log.Printf("[INFO] [filePathtoGCS] File %v exists in bucket %v\n", inspectsGCSTestFileName, bucketNameForInspectGCSSendToScc)
	}

	log.Println("[INFO] [filePathtoGCS] filePathtoGCS function is executed-------")
	return nil
}

func deleteStoredInfoTypeAfterTest(t *testing.T, name string) error {
	t.Helper()
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &dlppb.DeleteStoredInfoTypeRequest{
		Name: name,
	}
	err = client.DeleteStoredInfoType(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

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
	createBigQueryDataSetId(tc.ProjectID)
	createTableInsideDataset(tc.ProjectID, dataSetID)

	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer c.Close()
	m.Run()
	deleteBigQueryAssets(tc.ProjectID)
	deleteActiveJob(tc.ProjectID, jobTriggerForInspectSample)
	deleteJobTriggerForInspectDataToHybridJobTrigger(tc.ProjectID, jobTriggerForInspectSample)
	if err := testutil.DeleteExpiredBuckets(c, tc.ProjectID, testPrefix, bucketExpiryAge); err != nil {
		// Don't fail the test if cleanup fails
		log.Printf("[INFO] [TestMain] Post-test cleanup failed: %v", err)
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
			log.Printf("[INFO] [deleteActiveJob] Next: %v", err)
		}
		log.Printf("[INFO] [deleteActiveJob] Job %v status: %v\n", j.GetName(), j.GetState())
	}
	for _, v := range jobIds {
		req := &dlppb.DeleteDlpJobRequest{
			Name: v,
		}
		if err = client.DeleteDlpJob(ctx, req); err != nil {
			log.Printf("[INFO] [deleteActiveJob] DeleteDlpJob: %v", err)
			return err
		}
		log.Printf("\n[INFO] [deleteActiveJob] Successfully deleted job %v\n", v)
	}
	log.Println("[INFO] [deleteActiveJob] Deleted Job Successfully !!!")
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
