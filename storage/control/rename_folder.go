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

// [START storage_control_rename_folder]
import (
	"context"
	"fmt"
	"io"
	"time"

	control "cloud.google.com/go/storage/control/apiv2"
	"cloud.google.com/go/storage/control/apiv2/controlpb"
)

// renameFolder changes the name of an existing folder.
func renameFolder(w io.Writer, bucket, src, dst string) error {
	// bucket := "bucket-name"
	// src := "original-folder-name"
	// dst := "new-folder-name"

	ctx := context.Background()
	client, err := control.NewStorageControlClient(ctx)
	if err != nil {
		return fmt.Errorf("NewStorageControlClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// Construct source folder path including the bucket name.
	srcPath := fmt.Sprintf("projects/_/buckets/%v/folders/%v", bucket, src)

	req := &controlpb.RenameFolderRequest{
		Name:                srcPath,
		DestinationFolderId: dst,
	}
	op, err := client.RenameFolder(ctx, req)
	if err != nil {
		return fmt.Errorf("RenameFolder(%q): %w", srcPath, err)
	}

	// Wait for long-running operation to complete.
	f, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for RenameFolder: %w", err)
	}

	fmt.Fprintf(w, "folder %v moved to new path %v", srcPath, f.Name)
	return nil
}

// [END storage_control_rename_folder]
