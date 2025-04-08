# Google Cloud Speech API Go example

## Authentication

* Create a project with the [Google Cloud Console][cloud-console], and enable
  the [Speech API][speech-api].
* From the Cloud Console, create a service account,
  download its json credentials file, then set the 
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

  ```bash
  export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
  ```

[cloud-console]: https://console.cloud.google.com
[speech-api]: https://console.cloud.google.com/apis/api/speech.googleapis.com/overview?project=_
[adc]: https://cloud.google.com/docs/authentication#developer_workflow

## Run the sample

Before running any example you must first install the Speech API client:

```bash
go get -u cloud.google.com/go/speech/apiv1
```

To run the example with a local file:

```bash
go run caption.go ../testdata/audio.raw
```

To run the example with a GCS file:

```bash
go run caption.go gs://...
```
