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

// [START auth_downscoping_token_consumer]

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/google/downscope"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// A token consumer should define their own tokenSource. In the Token() method,
// it should send a query to a token broker requesting a downscoped token.
// The token broker holds the root credential that is used to generate the
// downscoped token.
type localTokenSource struct {
	ctx        context.Context
	bucketName string
	brokerURL  string
}

func (lts localTokenSource) Token() (*oauth2.Token, error) {
	var remoteToken *oauth2.Token
	// Usually you would now retrieve remoteToken, an oauth2.Token, from token broker.
	// This snippet performs the same functionality locally.
	accessBoundary := []downscope.AccessBoundaryRule{
		{
			AvailableResource:    "//storage.googleapis.com/projects/_/buckets/" + lts.bucketName,
			AvailablePermissions: []string{"inRole:roles/storage.objectViewer"},
		},
	}
	rootSource, err := google.DefaultTokenSource(lts.ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("failed to generate rootSource: %w", err)
	}
	dts, err := downscope.NewTokenSource(lts.ctx, downscope.DownscopingConfig{RootSource: rootSource, Rules: accessBoundary})
	if err != nil {
		return nil, fmt.Errorf("failed to generate downscoped token source: %w", err)
	}
	// Token() uses the previously declared TokenSource to generate a downscoped token.
	remoteToken, err = dts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return remoteToken, nil
}

// getObjectContents will read the contents of an object in Google Storage
// named objectName, contained in the bucket "bucketName".
func getObjectContents(output io.Writer, bucketName string, objectName string) error {
	// bucketName := "foo"
	// prefix := "profile-picture-"

	ctx := context.Background()

	thisTokenSource := localTokenSource{
		ctx:        ctx,
		bucketName: bucketName,
		brokerURL:  "yourURL.com/internal/broker",
	}

	// Wrap the TokenSource in an oauth2.ReuseTokenSource to enable automatic refreshing.
	refreshableTS := oauth2.ReuseTokenSource(nil, thisTokenSource)
	// You can now use the token source to access Google Cloud Storage resources as follows.
	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(refreshableTS))
	if err != nil {
		return fmt.Errorf("failed to create the storage client: %w", err)
	}
	defer storageClient.Close()
	bkt := storageClient.Bucket(bucketName)
	obj := bkt.Object(objectName)
	rc, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve the object: %w", err)
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("could not read the object's contents: %w", err)
	}
	// Data now contains the contents of the requested object.
	output.Write(data)
	return nil
}

// [END auth_downscoping_token_consumer]
