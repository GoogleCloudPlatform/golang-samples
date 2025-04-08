# Cloud Run CI

This utility facilitates deploying temporary Cloud Run services for testing purposes.

It has a hard dependency on `gcloud`, the [Cloud SDK](https://cloud.google.com/sdk/).

Please install and authenticate gcloud before using cloudrunci in your test.

## Installation

```go
import "github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
```

## Configuration

Use the `GCLOUD_BIN` environment variable to override the gcloud path.
