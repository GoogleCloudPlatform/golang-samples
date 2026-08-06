// Copyright 2026 Google LLC
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

package control

// [START storage_control_resume_anywhere_cache]
import (
	"context"
	"fmt"
	"io"
	"time"

	storagecontrol "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// resumeAnywhereCache resumes an anywhere cache.
func resumeAnywhereCache(w io.Writer, cacheName string) error {
	// cacheName := "projects/_/buckets/bucket-name/anywhereCaches/us-central1-a"

	ctx := context.Background()
	client, err := storagecontrol.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("storagecontrol.NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	req := &controlpb.ResumeAnywhereCacheRequest{
		Name: cacheName,
	}

	anywhereCache, err := client.ResumeAnywhereCache(ctx, req)
	if err != nil {
		return fmt.Errorf("ResumeAnywhereCache(%q): %w", cacheName, err)
	}

	fmt.Fprintf(w, "Resumed anywhere cache: %s\n", anywhereCache.GetName())
	return nil
}

// [END storage_control_resume_anywhere_cache]
