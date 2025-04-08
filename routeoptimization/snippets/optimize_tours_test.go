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

package snippets

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestOptimizeTours(t *testing.T) {
	st := testutil.SystemTest(t)
	got, err := optimizeTours(st.ProjectID)
	if err != nil {
		t.Fatalf("optimizeTours: %v", err)
	}
	if got.GetMetrics() == nil {
		t.Fatalf("optimizeTours response %v has no metrics field", got)
	}

	if count := got.GetMetrics().GetAggregatedRouteMetrics().PerformedShipmentCount; count != 1 {
		t.Fatalf("response performed_shipment_count: got %d want 1", count)
	}
}
