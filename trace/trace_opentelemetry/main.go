// Copyright 2020 Google LLC
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

// [START trace_opentelemetry]

func main() {

    // Create exporter
    exporter, err := stackdriver.NewExporter(
        stackdriver.WithProjectID("PROJECT-ID"),
        )

    // Create trace provider with the exporter
    tp, err :=
    sdktrace.NewProvider(
    sdktrace.WithConfig(
    // AlwaysSample() is used here to make sure traces are available for
    // observation and analysis. In a production environment or high QPS
    // setup please use ProbabilitySampler set at the desired probability.
    // Example: DefaultSampler:sdktrace.ProbabilitySampler(0.0001)
    sdktrace.Config{DefaultSampler:sdktrace.AlwaysSample()}),
    sdktrace.WithSyncer(exporter))
    if err != nil {
            log.Fatal(err)
    }
	global.SetTraceProvider(tp)
	
// [END trace_opentelemetry]