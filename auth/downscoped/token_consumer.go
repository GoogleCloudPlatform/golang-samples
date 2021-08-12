// Copyright 2021 Google LLC
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

package downscopedoverview

// [START downscoping_token_consumer]

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// A token consumer should define their own tokenSource. In the Token() method,
// it should send a query to a token broker requesting a downscoped token.
// The token broker holds the root credential that is used to generate the
// downscoped token.
type localTokenSource struct {
	requestedObject string
	brokerURL       string
}

func (localTokenSource) Token() (*oauth2.Token, error) {
	var remoteToken oauth2.Token
	// Retrieve remoteToken, an oauth2.Token, from token broker.
	return &remoteToken, nil
}

// getObjectContents will read the contents of an object in Google Storage
// named "myFile.txt", contained in the bucket "foo"
func getObjectContents() ([]byte, error) {
	ctx := context.Background()
	thisTokenSource := localTokenSource{
		requestedObject: "//storage.googleapis.com/projects/_/buckets/foo",
		brokerURL:       "yourURL.com/internal/broker",
	}

	// Wrap the TokenSource in an oauth2.ReuseTokenSource to enable automatic refreshing.
	refreshableTS := oauth2.ReuseTokenSource(nil, thisTokenSource)
	// You can now use the token source to access Google Cloud Storage resources as follows.
	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(refreshableTS))
	defer storageClient.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to create the storage client: %v", err)
	}
	bkt := storageClient.Bucket("foo")
	obj := bkt.Object("myFile.txt")
	rc, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the object: %v", err)
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("could not read the object's contents: %v", err)
	}
	return data, err
}

// [END downscoping_token_consumer]
