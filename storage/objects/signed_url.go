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

import (
	"fmt"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
)

func signedURL(bucket, object string) error {
	// Download a p12 service account private key from the Google Developers Console.
	// And convert it to PEM by running the command below:
	//	$ openssl pkcs12 -in key.p12 -passin pass:notasecret -out my-private-key.pem -nodes
	pkey, err := ioutil.ReadFile("my-private-key.pem")
	if err != nil {
		return err
	}
	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: "xxx@developer.gserviceaccount.com",
		PrivateKey:     pkey,
		Method:         "GET",
		Expires:        time.Now().Add(48 * time.Hour),
	})
	if err != nil {
		return err
	}
	fmt.Println(url)
	return nil
}
