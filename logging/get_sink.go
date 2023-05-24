// Copyright 2023 Google LLC
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

package snippets

// [START logging_get_sink]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/logging/logadmin"
)

// getSink retrieves the metadata for a Cloud Logging Sink.
func getSink(w io.Writer, projectID, sinkName string) error {
	ctx := context.Background()

	client, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	sink, err := client.Sink(ctx, sinkName)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%v\n", sink)
	return nil
}

// [END logging_get_sink]
