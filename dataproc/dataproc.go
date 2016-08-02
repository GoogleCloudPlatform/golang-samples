// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This script creates a Google Cloud Dataproc cluster
// submits a job to the cluster and then deletes the cluster
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	dataproc "google.golang.org/api/dataproc/v1"
	storage "google.golang.org/api/storage/v1"
)

const (
	// Filename of example job (Python) file uploaded to Cloud Storage
	exampleJobName = "test-file"

	// Allows for full access to Google Cloud Platform products
	scope = "https://www.googleapis.com/auth/cloud-platform"

	// Compute URI base
	computeUriBase = "https://www.googleapis.com/compute/v1/projects/%s/zones/%s"
)

var (
	bucketName  = flag.String("bucket", "", "GCS Bucket for storage of job file")
	clusterName = flag.String("cluster-name", "", "Name of cluster")
	projectID   = flag.String("project", "", "Your cloud project ID.")
	pysparkFile = flag.String("pyspark-file", "", "PySpark Job file")
	regionID    = flag.String("region", "global", "Your cloud project region.")
	zoneID      = flag.String("zone", "", "Cloud Platform zone")
)

// Get a Cloud Dataproc service object
func GetDataprocService() (service *dataproc.Service) {
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Dataproc service
	service, err = dataproc.New(client)
	if err != nil {
		log.Fatal(err)
	}

	return service
}

// Get a Cloud Storage service object
func GetStorageService() (service *storage.Service) {
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new storage service
	service, err = storage.New(client)
	if err != nil {
		log.Fatal(err)
	}

	return service
}

// Create a Cloud Dataproc cluster
func CreateCluster(service *dataproc.Service, name string, region string,
	project string, zone string) (response string, err error) {
	// Create a gceConfig object for the cluster
	gceConfig := dataproc.GceClusterConfig{
		ZoneUri: fmt.Sprintf(computeUriBase, project, zone),
	}

	// Create a cluserConfig for the cluster
	clusterConfig := dataproc.ClusterConfig{
		GceClusterConfig: &gceConfig,
	}

	// Create a cluster object
	cluster := dataproc.Cluster{
		ClusterName: name,
		ProjectId:   project,
		Config:      &clusterConfig,
	}

	// Create the cluster
	res, err := service.Projects.Regions.Clusters.Create(project, region, &cluster).Do()

	// Handle creation errors
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", res), err
}

// Delete the Cloud Dataproc cluster
func DeleteCluster(service *dataproc.Service, project string, region string, name string) (response string, err error) {
	// Delete the cluster
	res, err := service.Projects.Regions.Clusters.Delete(project, region, name).Do()

	// Handle deletion errors
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", res), err
}

// Get metadata about a cluster
func GetClusterDataByName(service *dataproc.Service, name string) (clusterData []string) {
	clusters, _ := ListClusters(service)
	for _, item := range clusters {
		if item[0] == name {
			return item
		}
	}

	return nil
}

// Get the link to a job's output
func GetJobOutput(storageService *storage.Service, projectId string, clusterId string, outputBucket string, jobId string) {
	objectName := fmt.Sprintf("google-cloud-dataproc-metainfo/%s/jobs/%s/driveroutput.000000000", clusterId, jobId)

	// URL encode the object name
	u, err := url.Parse(objectName)
	if err != nil {
		log.Fatal(err)
	}
	encodedObjectName := u.EscapedPath()

	if res, err := storageService.Objects.Get(outputBucket, encodedObjectName).Do(); err == nil {
		fmt.Println(fmt.Printf("The media download link for the output is:\n\n%s.\n\n", res.MediaLink))
		fmt.Printf("The media download link for %v/%v is %v.\n", outputBucket, res.Name, res.MediaLink)
	} else {
		log.Fatal(fmt.Sprintf("Failed to get %s/%s: %s.", outputBucket, encodedObjectName, err))
	}
}

// Get the status of a job
func GetJobStatus(service *dataproc.Service, jobId string, project string, region string) (status string, err error) {
	//Get the Job's status
	res, err := service.Projects.Regions.Jobs.Get(project, region, jobId).Do()

	return res.Status.State, err
}

func ListClusters(service *dataproc.Service) (clusters [][]string, err error) {
	// List all clusters in a project for a given region
	res, err := service.Projects.Regions.Clusters.List(*projectID, *regionID).Do()

	for _, item := range res.Clusters {
		clusterDetails := []string{
			item.ClusterName,
			item.ClusterUuid,
			item.Status.State,
			item.Config.ConfigBucket,
		}
		clusters = append(clusters, clusterDetails)
	}

	return clusters, err
}

func SubmitJob(service *dataproc.Service, stroageService *storage.Service, filepath string, project string, bucket string, cluster string) (jobId string, err error) {
	// Error if the file does not exist
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatal("File does not exist")
	}

	// Upload the file to GCS
	job_filename_parts := strings.Split(filepath, "/")
	filename := job_filename_parts[len(job_filename_parts)-1]
	object := &storage.Object{Name: filename}
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening %q: %v", filepath, err)
	}
	_, uploadError := stroageService.Objects.Insert(bucket, object).Media(file).Do()
	if uploadError != nil {
		log.Fatal("Error uploading file to GCS", filepath, err)
	}

	// Submit the PySpark job
	placement := dataproc.JobPlacement{
		ClusterName: cluster,
	}
	pySparkJob := dataproc.PySparkJob{
		MainPythonFileUri: "gs://" + bucket + "/" + filename,
	}
	jobID := "test-job-" + fmt.Sprintf("%v", time.Now().Unix())
	jobReference := dataproc.JobReference{
		JobId:     jobID,
		ProjectId: project,
	}
	job := dataproc.Job{
		Placement:  &placement,
		PysparkJob: &pySparkJob,
		Reference:  &jobReference,
	}
	jobRequest := dataproc.SubmitJobRequest{
		Job: &job,
	}

	_, jobError := service.Projects.Regions.Jobs.Submit(*projectID, "global", &jobRequest).Do()
	return jobID, jobError
}

func WaitForCluster(service *dataproc.Service, name string) (running bool, err error) {
	for running == false {
		clusterData := GetClusterDataByName(service, name)
		if clusterData[2] == "RUNNING" {
			running = true
		}

		// Sleep for one second
		time.Sleep(1000 * time.Millisecond)
	}
	return running, err
}

// Wait for a job to complete
func WaitForJob(service *dataproc.Service, jobId string, project string, region string) (finished bool, err error) {
	for finished == false {
		jobStatus, err := GetJobStatus(service, jobId, project, region)
		if err != nil {
			log.Fatal(err)
		}
		if jobStatus == "DONE" {
			finished = true
		}

		// Sleep for one second
		time.Sleep(1000 * time.Millisecond)
	}

	return finished, err
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Check for required flags (golang flags does not support?)
	if *bucketName == "" || *clusterName == "" || *projectID == "" || *pysparkFile == "" || *regionID == "" || *zoneID == "" {
		log.Fatal("Incorrect arguments specified see 'go-dataproc -help' for help")
	}

	// Create a new Dataproc service
	service := GetDataprocService()

	// Create a new storage service
	storageService := GetStorageService()

	// Create a cluster
	_, create_error := CreateCluster(service, *clusterName, *regionID, *projectID, *zoneID)
	if create_error != nil {
		log.Fatal(create_error)
	}
	fmt.Println("Cluster created")

	// Wait on the cluster to become active
	fmt.Println("Waiting for cluster to be ready")
	_, clusterWaitError := WaitForCluster(service, *clusterName)
	if clusterWaitError != nil {
		log.Fatal(clusterWaitError)
	}

	// Submit a job to the cluster
	fmt.Println("Submitting job")
	jobId, job_error := SubmitJob(service, storageService, *pysparkFile, *projectID, *bucketName, *clusterName)
	if job_error != nil {
		log.Fatal(job_error)
	}

	// Wait for the job to complete
	fmt.Println("Waiting for job to complete")
	WaitForJob(service, jobId, *projectID, *regionID)
	fmt.Println("Job is finished")

	// Get the job outputBucket
	clusterData := GetClusterDataByName(service, *clusterName)
	GetJobOutput(storageService, *projectID, clusterData[1], clusterData[3], jobId)

	// Delete the cluster
	_, delete_error := DeleteCluster(service, *projectID, *regionID, *clusterName)
	if delete_error != nil {
		log.Fatal(delete_error)
	}
	fmt.Println("Cluster deleted")
	fmt.Printf("Done...")
}

