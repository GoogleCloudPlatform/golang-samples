// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package contextcaching

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// TestContextCaching tests createContextCache, getContextCache, useContextCache,
// updateContextCache, deleteContextCache.
func TestContextCaching(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	location := "us-central1"
	modelName := "gemini-1.5-pro-001"

	// 1) Create a cached content. The generated content name will be used in steps 2, 3, 4.
	contentName, err := createContextCache(buf, tc.ProjectID, location, modelName)
	if err != nil {
		t.Fatalf("createContextCache: %v", err.Error())
	}

	// 2) Retrieve the cached content, by its name.
	buf.Reset()
	err = getContextCache(buf, contentName, tc.ProjectID, location)
	if err != nil {
		t.Errorf("getContextCache: %v", err.Error())
	}

	// 3) Retrieve the list of cached contents
	buf.Reset()
	err = listContextCaches(buf, tc.ProjectID, location)
	if err != nil {
		t.Errorf("listContextCache: %v", err.Error())
	}

	// 4) Use the cached content, by calling the model with a prompt
	buf.Reset()
	err = useContextCache(buf, contentName, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("useContextCache: %v", err.Error())
	}

	// 5) Update the TTL of the cached content
	exp1, err := readExpiration(contentName, tc.ProjectID, location)
	if err != nil {
		t.Errorf("readExpiration original value: %v", err.Error())
	}
	buf.Reset()
	err = updateContextCache(buf, contentName, tc.ProjectID, location)
	if err != nil {
		t.Errorf("updateContextCache: %v", err.Error())
	}
	exp2, err := readExpiration(contentName, tc.ProjectID, location)
	if err != nil {
		t.Errorf("readExpiration updated value: %v", err.Error())
	}
	// We've created the cached content with a TTL of "60 minutes",
	// then updated it with a TTL of "2 hours".
	// The new expiration time should be slightly more than 1 hour after
	// the original expiration time.
	if delta := exp2.Sub(exp1); delta < 1*time.Hour {
		t.Errorf("want expiration time at least 1 hour greater than original value %v, got %v (diff=%v)", exp1, exp2, delta)
	}

	// 6) Delete the cached content
	buf.Reset()
	err = deleteContextCache(buf, contentName, tc.ProjectID, location)
	if err != nil {
		t.Errorf("deleteContextCache: %v", err.Error())
	}
	// The cached content must not exist anymore
	buf.Reset()
	err = getContextCache(buf, contentName, tc.ProjectID, location)
	if err == nil {
		t.Errorf("No error when retrieving deleted cached content: %s", buf.Bytes())
	}
}

// readExpiration is a helper that retrieves a cached content from the service, and
// return its expiration time.
// The retrieved cached content has a populated expiration time, but no TTL.
func readExpiration(contentName string, projectID, location string) (time.Time, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	cachedContent, err := client.GetCachedContent(ctx, contentName)
	if err != nil {
		return time.Time{}, fmt.Errorf("GetCachedContent: %w", err)
	}
	return cachedContent.Expiration.ExpireTime, nil
}
