# Google Cloud Vision API Go example

## Authentication

* Create a project with the [Google Cloud Console][cloud-console], and enable
  the [Vision API][vision-api].
* From the Cloud Console, create a service account,
  download its json credentials file, then set the 
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

  ```bash
  export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
  ```

[cloud-console]: https://console.cloud.google.com
[vision-api]: https://console.cloud.google.com/apis/api/vision.googleapis.com/overview?project=_
[adc]: https://cloud.google.com/docs/authentication#developer_workflow

## Run the sample

To build and run the sample:

```bash
go run detect.go main.go ../testdata/cat.jpg

go run detect.go main.go gs://...

go run detect.go main.go https://...
```

## Modifiying the source code

Do not edit the `detect.go` file directly. In order to modify it edit the code at `vision/detect/generated/sample-template.go` and run `go generate` at the command line:

```bash
cd vision/detect
go generate
```