# Google Cloud Vision API Product Search Go example

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

## Run the samples

To build and run the samples, you need to provide your `PROJECT_ID` and a `LOCATION`, which is a GCP compute zone such as `us-west1`.

```bash
go run delete_product.go PROJECT_ID LOCATION PRODUCT_ID
```

```bash
go run list_products.go PROJECT_ID LOCATION
```

```bash
go run get_product.go PROJECT_ID LOCATION PRODUCT_ID
```

```bash
go run create_product.go PROJECT_ID LOCATION PRODUCT_ID PRODUCT_DISPLAY_NAME PRODUCT_CATEGORY
```
