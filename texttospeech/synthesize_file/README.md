# Google Cloud Text-to-Speech API Go example

## Authentication

* Create a project with the [Google Cloud Console][cloud-console], and enable
  the [Text-to-Speech API][text-to-speech-api].
* From the Cloud Console, create a service account,
  download its json credentials file, then set the 
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

  ```bash
  export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
  ```

[cloud-console]: https://console.cloud.google.com
[text-to-speech-api]: https://console.cloud.google.com/apis/api/texttospeech.googleapis.com/overview?project=_
[adc]: https://cloud.google.com/docs/authentication#developer_workflow

## Run the sample

Before running any example you must first install the Text-to-Speech API client:

```bash
go get -u cloud.google.com/go/texttospeech/apiv1
```

To run the example with a local text file:

```bash
go run synthesize_file.go --text ../resources/hello.txt
```

To run the example with a local SSML file:

```bash
go run synthesize_file.go --ssml ../resources/hello.ssml
```

