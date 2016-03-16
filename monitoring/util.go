// Copyright 2015 Google, Inc
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

// This example shows how to create a custom metric and write TimeSeries value
// to it. It writes a GAUGE measurement, which is a measure of value at a
// specific point in time. This means the startTime and endTime of the interval
// are the same. To make it easier to see the output, a random value is written.
// When reading the TimeSeries back, a window of the last 5 minutes is used.
// See README.md for instructions on how to run.

package main

import (
	"encoding/json"
)

// printResource prints out our API response objects as JSON.
func formatResource(resource interface{}) []byte {
	b, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}
