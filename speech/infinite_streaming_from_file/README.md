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

Before running any example you must first install the Speech API client and the Color Pen:

```bash
go get -u cloud.google.com/go/speech/apiv1
go get -u google.golang.org/genproto/googleapis/cloud/speech/v1
go get -u github.com/fatih/color
```

To run the example with a local file:

```bash
go build
cat ../testdata/audio.raw | infinite_audio
```

## Capturing audio from the mic

Alternatively, `SoX` can be used to capture audio from the mic. For example:

```bash
go build
sox -d -e signed -b 16 -c 1 -q mic.wav | infinite_mic
```

## Content Limits

The Speech API contains the following limits on the size of content (and are subject to change):

| Content Limit	| Audio Length |
| ------------- | ------------ |
| Synchronous Requests | ~1 Minute |
| Asynchronous Requests	| ~180 Minutes |
| Streaming Requests | ~1 Minute |

Please note that each `StreamingRecognize` session is considered a single request even though it includes multiple frames of `StreamingRecognizeRequest` audio within the stream.

For more information, please refer to https://cloud.google.com/speech/limits#content.
