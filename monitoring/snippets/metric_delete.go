// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
