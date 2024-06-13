// Copyright 2024 Google LLC
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

// [START storage_control_delete_folder]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// deleteFolder deletes the folder with the given name.
func deleteFolder(w io.Writer, bucket, folder string) error {
	// bucket := "bucket-name"
	// folder := "folder-name"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// Construct folder path including the bucket name.
	folderPath := fmt.Sprintf("projects/_/buckets/%v/folders/%v", bucket, folder)

	req := &controlpb.DeleteFolderRequest{
		Name: folderPath,
	}
	if err := client.DeleteFolder(ctx, req); err != nil {
		return fmt.Errorf("DeleteFolder(%q): %w", folderPath, err)
	}

	fmt.Fprintf(w, "deleted folder %q", folderPath)
	return nil
}

// [END storage_control_delete_folder]
