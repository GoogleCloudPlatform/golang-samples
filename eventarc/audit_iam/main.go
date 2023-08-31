// Copyright 2023 Google LLC
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

// [START eventarc_audit_iam_handler]

// Processes CloudEvents containing Cloud Audit Logs for IAM
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cloudevent "github.com/cloudevents/sdk-go/v2"
	"github.com/googleapis/google-cloudevents-go/cloud/auditdata"
	"google.golang.org/protobuf/encoding/protojson"
)

func HandleCloudEvent(w http.ResponseWriter, r *http.Request) {
	// Transform the HTTP request into a CloudEvent
	event, err := cloudevent.NewEventFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to create CloudEvent from request.")
		log.Fatal("cloudevent.NewEventFromHTTPRequest:", err)
	}

	// Extract the LogEntryData from the CloudEvent
	var logentry auditdata.LogEntryData
	// AuditLog objects include a `@type` annotation, which errors when using
	// `protojson.Unmarshal`. UnmarshalOptions prevents this error.
	umo := &protojson.UnmarshalOptions{DiscardUnknown: true}
	err = umo.Unmarshal(event.Data(), &logentry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to parse Audit Log")
		log.Fatal("protojson.Unmarshal:", err)
	}

	// Extract relevant fields from the audit log entry.
	// Identify the user that requested key creation
	actor := logentry.ProtoPayload.AuthenticationInfo.PrincipalEmail

	// Extract the resource name from the CreateServiceAccountKey request
	// For details of this type, see https://cloud.google.com/iam/docs/reference/rpc/google.iam.admin.v1#createserviceaccountkeyrequest
	principal := logentry.ProtoPayload.GetRequest().AsMap()["name"]

	// The response is of type google.iam.admin.v1.ServiceAccountKey,
	// which is described at https://cloud.google.com/iam/docs/reference/rpc/google.iam.admin.v1#google.iam.admin.v1.ServiceAccountKey
	// This key path can be used with gcloud to disable/delete the key:
	// e.g. gcloud iam service-accounts keys disable ${keypath}
	keypath := logentry.ProtoPayload.GetResponse().AsMap()["name"]

	s := fmt.Sprintf("New Service Account Key created for %s by %s: %v", principal, actor, keypath)
	log.Printf(s)
	fmt.Fprintln(w, s)
}

// [END eventarc_audit_iam_handler]
// [START eventarc_audit_iam_server]

func main() {
	// disable leading timestamp, since it is automatic with Cloud Logging.
	log.SetFlags(0)

	http.HandleFunc("/", HandleCloudEvent)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// [END eventarc_audit_iam_server]
