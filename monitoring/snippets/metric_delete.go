// Copyright 2019 Google LLC
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

// [START monitoring_delete_metric]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// deleteMetric deletes the given metric. name should be of the form
// "projects/PROJECT_ID/metricDescriptors/METRIC_TYPE".
func deleteMetric(w io.Writer, name string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	req := &monitoringpb.DeleteMetricDescriptorRequest{
		Name: name,
	}

	if err := c.DeleteMetricDescriptor(ctx, req); err != nil {
		return fmt.Errorf("could not delete metric: %v", err)
	}
	fmt.Fprintf(w, "Deleted metric: %q\n", name)
	return nil
}

// [END monitoring_delete_metric]
