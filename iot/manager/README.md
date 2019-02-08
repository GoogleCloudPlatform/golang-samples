# Google Cloud IoT Core Device Manager Sample

## Authentication

* Create a project with the [Google Cloud Console][cloud-console], and enable
  the [Cloud IoT Core API][cloud-iot-api].
* From the Cloud Console, create a service account,
  download its json credentials file, then set the `GCLOUD_PROJECT` and
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

  ```bash
  export GCLOUD_PROJECT=your-project-id
  export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
  ```

[cloud-console]: https://console.cloud.google.com
[cloud-iot-api]: https://console.developers.google.com/iot

## Run the sample

To build and run the sample:

```bash
go run manager.go
```

Before you can connect devices, you will need to generate public / private key
credentials to register and connect.

The following commands generate keys suitable for use on a single device.

```bash
openssl req -x509 -newkey rsa:2048 -days 3650 -keyout rsa_private.pem -nodes -out \
    rsa_cert.pem -subj "/CN=unused"
openssl ecparam -genkey -name prime256v1 -noout -out ec_private.pem
openssl ec -in ec_private.pem -pubout -out ec_public.pem
```
