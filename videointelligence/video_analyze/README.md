# Google Cloud Video Intelligence API Go example

## Authentication

* Create a project with the [Google Cloud Console][cloud-console], and enable
  the [Video Intelligence API][video-api].
* From the Cloud Console, create a service account,
  download its json credentials file, then set the 
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

  ```bash
  export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
  ```

[cloud-console]: https://console.cloud.google.com
[video-api]: https://console.cloud.google.com/apis/api/videointelligence.googleapis.com/overview?project=_
[adc]: https://cloud.google.com/docs/authentication#developer_workflow

## Run the sample

To build and run the sample:

```bash
go build -o video_analyze

gcloud storage cp gs://cloud-samples-data/video/cat.mp4 .
./video_analyze cat.mp4

./video_analyze gs://cloud-samples-data/video/cat.mp4
```
