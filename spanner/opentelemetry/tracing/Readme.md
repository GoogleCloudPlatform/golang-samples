# Cloud Spanner OpenTelemetry Traces

## Setup

This sample requires [Go](https://go.dev/doc/install) and [gcloud](https://cloud.google.com/sdk/docs/install).

Ensure that your Go runtime version is supported by the OpenTelemetry-Go compatibility policy before enabling OpenTelemetry instrumentation. 
Refer to compatibility here https://github.com/googleapis/google-cloud-go/blob/main/debug.md#opentelemetry

1.  **Follow the set-up instructions in [the documentation](https://cloud.google.com/go/docs/setup).**

2.  Enable APIs for your project.

    a. [Click here](https://console.cloud.google.com/flows/enableapi?apiid=spanner.googleapis.com&showconfirmation=true)
    to visit Cloud Platform Console and enable the Google Cloud Spanner API.

    b. [Click here](https://console.cloud.google.com/flows/enableapi?apiid=cloudtrace.googleapis.com&showconfirmation=true)
    to visit Cloud Platform Console and enable the Cloud Trace API.

3.  Create a Cloud Spanner instance and database via the Cloud Plaform Console's
    [Cloud Spanner section](http://console.cloud.google.com/spanner).

4.  Enable application default credentials by running the command `gcloud auth application-default login`.

## Run the Example

1. Set up database configuration in the `spanner_opentelemetry_tracing.go` file:
    ```
    var projectID = "projectID"
    var instanceID = "instanceID"
    var databaseID = "databaseID"
    ```

2. Configure trace data export. You can use either the OpenTelemetry [Collector](https://opentelemetry.io/docs/collector/quick-start/) with the OTLP Exporter or the Cloud Trace Exporter. By default, the Cloud Trace Exporter is used.

- To use OTLP Exporter, Set up the OpenTelemetry [Collector](https://opentelemetry.io/docs/collector/quick-start/) and update the OTLP endpoint in `spanner_opentelemetry_tracing.go` file
    ```
    var useCloudTraceExporter = true; // Replace to false for OTLP
    defaultOtlpEndpoint := "http://localhost:4317"; // Replace with your OTLP endpoint
    ```

3. Enable OpenTelemetry traces by setting environment variable.
    ```
    export GOOGLE_API_GO_EXPERIMENTAL_TELEMETRY_PLATFORM_TRACING="opentelemetry"
    ```

4. Then run the application from command line, after switching to this directory:
    ```
    go run spanner_opentelemetry_tracing.go
    ```

You should start seeing traces in Cloud Trace .
