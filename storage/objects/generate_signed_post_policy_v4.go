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

package objects

// [START storage_generate_signed_post_policy_v4]
import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

// generateSignedPostPolicyV4 generates a signed post policy.
func generateSignedPostPolicyV4(w io.Writer, bucket, object, serviceAccountJSON string) (*storage.PostPolicyV4, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	// serviceAccountJSON := "service_account.json"
	jsonKey, err := ioutil.ReadFile(serviceAccountJSON)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	metadata := map[string]string{
		"x-goog-meta-test": "data",
	}
	opts := &storage.PostPolicyV4Options{
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(10 * time.Minute),
		Fields:         &storage.PolicyV4Fields{Metadata:metadata},
	}

	policy, err := storage.GenerateSignedPostPolicyV4(bucket, object, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.GenerateSignedPostPolicyV4: %v", err)
	}

	// Create an HTML form with the provided policy.
	fmt.Fprintf(w, "<form action='%v' method='POST' enctype='multipart/form-data'>\n", policy.URL)

	// Include all fields returned in the HTML form as they're required.
	for k, v := range policy.Fields {
		fmt.Fprintf(w, "  <input name='%v' value='%v' type='hidden'/>\n", k, v)
	}

	fmt.Fprint(w, "  <input type='file' name='file'/><br />\n")
	fmt.Fprint(w,  "  <input type='submit' value='Upload File' name='submit'/><br />\n")
	fmt.Fprint(w, "</form>")

	return policy, nil
}
// [END storage_generate_signed_post_policy_v4]
