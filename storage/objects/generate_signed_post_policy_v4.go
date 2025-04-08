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
	"context"
	"fmt"
	"html/template"
	"io"
	"time"

	"cloud.google.com/go/storage"
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
func generateSignedPostPolicyV4(w io.Writer, bucket, object string) (*storage.PostPolicyV4, error) {
	// bucket := "bucket-name"
	// object := "object-name"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	metadata := map[string]string{
		"x-goog-meta-test": "data",
	}

	// Generating a signed POST policy requires credentials authorized to sign a URL.
	// You can pass these in through PostPolicyV4Options with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing
	// In this example, none of these options are used, which means the
	// GenerateSignedPostPolicyV4 function attempts to use the same authentication
	// that was used to instantiate	the Storage client. This authentication must
	// include a private key or have iam.serviceAccounts.signBlob permissions.
	opts := &storage.PostPolicyV4Options{
		Expires: time.Now().Add(10 * time.Minute),
		Fields: &storage.PolicyV4Fields{
			Metadata: metadata,
		},
	}

	policy, err := client.Bucket(bucket).GenerateSignedPostPolicyV4(object, opts)
	if err != nil {
		return nil, fmt.Errorf("storage.GenerateSignedPostPolicyV4: %w", err)
	}

	// Generate the form, using the data from the policy.
	if err = tmpl.Execute(w, policy); err != nil {
		return policy, fmt.Errorf("executing template: %w", err)
	}

	return policy, nil
}

// [END storage_generate_signed_post_policy_v4]
