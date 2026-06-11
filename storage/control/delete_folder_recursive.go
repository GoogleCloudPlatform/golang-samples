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

// [START storage_control_delete_folder_recursive]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// deleteFolderRecursive deletes a folder recursively.
// This operation is only applicable to a hierarchical namespace enabled bucket.
func deleteFolderRecursive(w io.Writer, bucketName, folderName string) error {
	// bucketName := "bucket-name"
	// folderName := "folder-name"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	// Set project to "_" to signify globally scoped bucket
	folderResourceName := fmt.Sprintf("projects/_/buckets/%v/folders/%v", bucketName, folderName)

	req := &controlpb.DeleteFolderRecursiveRequest{
		Name: folderResourceName,
	}

	// Execute the deleteFolderRecursive method synchronously (waiting for the LRO to complete).
	op, err := client.DeleteFolderRecursive(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteFolderRecursive(%q): %w", folderResourceName, err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Deleted folder: %v\n", folderResourceName)
	return nil
}

// [END storage_control_delete_folder_recursive]
