// Copyright 2025 Google LLC
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var (
	region       string
	projectID    string
	serviceName  string
	revisionName string
	instanceId   string
	bucketName   string
	ctx          context.Context = context.Background()
	gcs          *storage.Client

	readinessProbeConfig ReadinessProbeConfig
	readinessEnabled     bool
	isHealthy            bool

	cachedInstances []InstanceView        = []InstanceView{}
	cachedRegions   map[string]RegionView = make(map[string]RegionView)

	tmpl = template.Must(template.New("layout.html").ParseFiles("layout.html"))
)

type RegionView struct {
	NumHealthy int
	Total      int
}
type InstanceView struct {
	InstanceId       string
	RevisionName     string
	Region           string
	HealthStr        string
	ReadinessEnabled bool
}

type InstanceMetadata struct {
	RevisionName     string
	Region           string
	ReadinessEnabled bool
}

// This struct mimics the v1 Service resource returned from the Cloud Run Admin API.
// We are interested in the readiness probe fields.
// https://cloud.google.com/run/docs/reference/rest/v1/namespaces.services#Service
type ReadinessProbeConfig struct {
	TimeoutSeconds   int `json:"timeoutSeconds"`
	PeriodSeconds    int `json:"periodSeconds"`
	SuccessThreshold int `json:"successThreshold"`
	FailureThreshold int `json:"failureThreshold"`
	HttpGetAction    struct {
		Path string `json:"path"`
		Port int    `json:"port"`
	} `json:"httpGet"`
}

type Service struct {
	Spec struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Name           string                `json:"name"`
					Image          string                `json:"image"`
					ReadinessProbe *ReadinessProbeConfig `json:"readinessProbe"`
				} `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

func init() {
	var err error

	var longRegion string
	if longRegion, err = queryMetadataServer("/computeMetadata/v1/instance/region"); err != nil {
		log.Fatal(err)
	}
	// region is of the format projects/12345/regions/us-central1
	regionSlice := strings.Split(longRegion, "/")
	region = regionSlice[len(regionSlice)-1]

	if projectID, err = queryMetadataServer("/computeMetadata/v1/project/project-id"); err != nil {
		log.Fatal(err)
	}

	serviceName = os.Getenv("K_SERVICE")
	revisionName = os.Getenv("K_REVISION")

	if instanceId, err = queryMetadataServer("/computeMetadata/v1/instance/id"); err != nil {
		log.Fatal(err)
	}

	bucketName = projectID + "-" + serviceName
	gcs, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	accessToken, err := getAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	apiURL := fmt.Sprintf(
		"https://%s-run.googleapis.com/apis/serving.knative.dev/v1/namespaces/%s/services/%s", region, projectID, serviceName)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code not ok %v", resp)
	}

	var serviceConfig Service
	if err = json.Unmarshal(body, &serviceConfig); err != nil {
		log.Fatal(err)
	}

	if serviceConfig.Spec.Template.Spec.Containers[0].ReadinessProbe != nil {
		readinessEnabled = true
		readinessProbeConfig = *serviceConfig.Spec.Template.Spec.Containers[0].ReadinessProbe
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := createRemoteConfigIfNone(); err != nil {
		log.Fatal(err)
	}

	if err := cache(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			if err := refreshReadinessConfig(); err != nil {
				log.Print(err)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := cleanUpStaleInstances(); err != nil {
				log.Print(err)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			if err := cache(); err != nil {
				log.Print(err)
			}

		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", rootRequestHandler)
	http.HandleFunc("/are_you_ready", areYouReadyHandler)
	http.HandleFunc("/set_readiness", setReadinessHandler)

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func cache() error {
	var sortedInstances []InstanceView
	var sortedString []string

	var regions = make(map[string]RegionView)

	ids, err := listInstances()
	if err != nil {
		return err
	}

	for _, id := range ids {
		meta, err := readMeta(id)
		if err != nil {
			continue
		}
		h, err := readHealth(id)
		if err != nil {
			continue
		}

		r, ok := regions[meta.Region]
		if !ok {
			numHealthy := 0
			if h && meta.ReadinessEnabled {
				numHealthy = 1
			}
			regions[meta.Region] = RegionView{
				NumHealthy: numHealthy,
				Total:      1,
			}
		} else {
			r.Total += 1
			if h && meta.ReadinessEnabled {
				r.NumHealthy += 1
			}
			regions[meta.Region] = r
		}

		idx, _ := slices.BinarySearch(sortedString, meta.Region+meta.RevisionName+id)
		sortedInstances = slices.Insert(sortedInstances, idx, InstanceView{
			InstanceId:       id,
			RevisionName:     meta.RevisionName,
			Region:           meta.Region,
			HealthStr:        getHealthStr(meta.ReadinessEnabled, h),
			ReadinessEnabled: meta.ReadinessEnabled,
		})
		sortedString = slices.Insert(sortedString, idx, meta.Region+meta.RevisionName+id)
	}

	cachedInstances = sortedInstances
	cachedRegions = regions

	return nil
}

func getHealthStr(enabled, healthy bool) string {
	if !enabled {
		return "NOT ENABLED ⚠️"
	}
	if healthy {
		return "HEALTHY ✅"
	}
	return "UNHEALTHY ❌"
}

func rootRequestHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, map[string]any{
		"Region":               region,
		"ServiceName":          serviceName,
		"ReadinessEnabled":     readinessEnabled,
		"Revision":             revisionName,
		"Project":              projectID,
		"Instances":            cachedInstances,
		"IsHealthy":            isHealthy,
		"ReadinessProbeConfig": readinessProbeConfig,
		"InstanceId":           instanceId,
		"Regions":              cachedRegions,
		"Bucket":               bucketName,
		"HealthStr":            getHealthStr(readinessEnabled, isHealthy),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func areYouReadyHandler(w http.ResponseWriter, r *http.Request) {
	if !readinessEnabled {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "NOT ENABLED")
		return
	}

	if isHealthy {
		fmt.Fprint(w, "HEALTHY")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "UNHEALTHY")
	}
}

func setReadinessHandler(w http.ResponseWriter, r *http.Request) {
	reqInstanceId := r.FormValue("instance_id")
	if reqInstanceId != "" {
		h, err := readHealth(reqInstanceId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = writeHealth(!h, reqInstanceId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	reqRegion := r.FormValue("region")
	if reqRegion != "" {
		reqHealthy := r.FormValue("is_healthy")
		newHealth := false
		if reqHealthy == "true" {
			newHealth = true
		}
		ids, err := listInstances()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, id := range ids {
			m, err := readMeta(id)
			if err != nil {
				continue
			}
			if m.Region == reqRegion {
				if err = writeHealth(newHealth, id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	var err error
	if isHealthy, err = readHealth(instanceId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	time.Sleep(2 * time.Second)

	if err := cache(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func refreshReadinessConfig() error {
	var err error
	isHealthy, err = readHealth(instanceId)
	if err != nil {
		return err
	}
	return writeHeartbeat()
}

func cleanUpStaleInstances() error {
	ids, err := listInstances()
	if err != nil {
		return err
	}
	for _, id := range ids {
		t, err := readHeartbeat(id)
		if err != nil {
			continue
		}
		if time.Since(*t) > 20*time.Second {
			deleteInstance(id)
		}
	}
	return nil
}

func listInstances() ([]string, error) {
	it := gcs.Bucket(bucketName).Objects(ctx, nil)

	instances := make(map[string]struct{})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var prefix string
		if strings.HasPrefix(attrs.Name, "meta-") {
			prefix = "meta-"
		} else if strings.HasPrefix(attrs.Name, "health-") {
			prefix = "health-"
		} else if strings.HasPrefix(attrs.Name, "heartbeat-") {
			prefix = "heartbeat-"
		} else {
			continue
		}
		id := strings.TrimPrefix(attrs.Name, prefix)
		instances[id] = struct{}{}
	}

	ids := make([]string, 0, len(instances))

	for id := range instances {
		ids = append(ids, id)
	}
	return ids, nil
}

func deleteInstance(instanceId string) {
	bucket := gcs.Bucket(bucketName)

	objMeta := bucket.Object("meta-" + instanceId)
	objMeta.Delete(ctx)

	objHeartbeat := bucket.Object("heartbeat-" + instanceId)
	objHeartbeat.Delete(ctx)

	objHealth := bucket.Object("health-" + instanceId)
	objHealth.Delete(ctx)
}

func createRemoteConfigIfNone() error {
	bucket := gcs.Bucket(bucketName)
	if _, err := bucket.Attrs(ctx); err != nil { // Bucket does not exist, create it.
		if err := bucket.Create(ctx, projectID, nil); err != nil {
			return err
		}
	}
	path := "meta-" + instanceId
	objMeta := bucket.Object(path)
	if _, err := objMeta.Attrs(ctx); err != nil { // Object does not exist, create it.
		if err = writeMetadata(newMeta()); err != nil {
			return err
		}
		isHealthy = true
		if err = writeHealth(isHealthy, instanceId); err != nil {
			return err
		}
		if err = writeHeartbeat(); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func writeHeartbeat() error {
	path := "heartbeat-" + instanceId
	w := gcs.Bucket(bucketName).Object(path).NewWriter(ctx)
	defer w.Close()
	_, err := w.Write([]byte(time.Now().UTC().Format(time.RFC3339)))
	return err
}

func readHeartbeat(instanceId string) (*time.Time, error) {
	path := "heartbeat-" + instanceId
	r, err := gcs.Bucket(bucketName).Object(path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func readHealth(instanceId string) (bool, error) {
	path := "health-" + instanceId
	r, err := gcs.Bucket(bucketName).Object(path).NewReader(ctx)
	if err != nil {
		return false, err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return false, err
	}
	return string(data) == "true", nil
}

func writeHealth(newHealth bool, instanceId string) error {
	newHealthStr := "false"
	if newHealth {
		newHealthStr = "true"
	}
	path := "health-" + instanceId
	w := gcs.Bucket(bucketName).Object(path).NewWriter(ctx)
	defer w.Close()
	_, err := w.Write([]byte(newHealthStr))
	return err
}

func writeMetadata(m *InstanceMetadata) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	path := "meta-" + instanceId
	w := gcs.Bucket(bucketName).Object(path).NewWriter(ctx)
	defer w.Close()
	_, err = w.Write(jsonData)
	return err
}

func readMeta(instanceId string) (*InstanceMetadata, error) {
	path := "meta-" + instanceId
	r, err := gcs.Bucket(bucketName).Object(path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var m InstanceMetadata
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func newMeta() *InstanceMetadata {
	return &InstanceMetadata{
		RevisionName:     revisionName,
		Region:           region,
		ReadinessEnabled: readinessEnabled,
	}
}

func getAccessToken() (string, error) {
	type accessTokenStruct struct {
		AccessToken string `json:"access_token"`
	}
	var token accessTokenStruct
	str, err := queryMetadataServer("/computeMetadata/v1/instance/service-accounts/default/token")
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(str), &token)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func queryMetadataServer(path string) (string, error) {
	metadataServerURL := "http://metadata.google.internal"
	req, err := http.NewRequest(http.MethodGet, metadataServerURL+path, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
