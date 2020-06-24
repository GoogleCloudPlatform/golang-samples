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
	"html/template"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

// form is a template for an HTML form that will use the data from the signed
// post policy.
var form = `<form action="{{ .URL }}" method="POST" enctype="multipart/form-data">
	{{- range $name, $value := .Fields }}
	<input name="{{ $name }}" value="{{ $value }}" type="hidden"/>
	{{- end }}
	<input type="file" name="file"/><br />
	<input type="submit" value="Upload File" name="submit"/><br />
</form>`

var tmpl = template.Must(template.New("policyV4").Parse(form))

// generateSignedPostPolicyV4 generates a signed post policy.
func generateSignedPostPolicyV4(w io.Writer, bucket, object, serviceAccountJSONPath string) (*storage.PostPolicyV4, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	// serviceAccountJSONPath := "service_account.json"
	jsonKey, err := ioutil.ReadFile(serviceAccountJSONPath)
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
		Fields: &storage.PolicyV4Fields{
			Metadata: metadata,
		},
	}

	policy, err := storage.GenerateSignedPostPolicyV4(bucket, object, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.GenerateSignedPostPolicyV4: %v", err)
	}

	// Generate the form, using the data from the policy.
	if err = tmpl.Execute(w, policy); err != nil {
		return policy, fmt.Errorf("executing template: %v", err)
	}

	return policy, nil
}

// [END storage_generate_signed_post_policy_v4]
