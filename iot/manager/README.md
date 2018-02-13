# Google Cloud IoT Core Device Manager Sample

## Authentication

* Create a project with the [Google Cloud Console][cloud-console], and enable
  the [Cloud IoT Core API][].
* From the Cloud Console, create a service account,
  download its json credentials file, then set the `GCLOUD_PROJECT` and
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

  ```bash
  export GCLOUD_PROJECT=your-project-id
  export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
  ```

[cloud-console]: https://console.cloud.google.com
[vision-api]: https://console.cloud.google.com/apis/api/cloudiot.googleapis.com/overview?project=_
[adc]: https://cloud.google.com/docs/authentication#developer_workflow

## Run the sample

To build and run the sample:

```bash
go run manager.go
```

Before you can connect devices, you will need to generate public / private key
credentials to register and connect. The `generate_keys.sh` script is included
to demonstrate how to do this for a single device.
