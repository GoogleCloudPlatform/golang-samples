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
go build
cat ../testdata/audio.raw | livecaption
```

## Capturing audio from the mic

Alternatively, `gst-launch` can be used to capture audio from the mic. For example:

```bash
gst-launch-1.0 -v pulsesrc ! audioconvert ! audioresample ! audio/x-raw,channels=1,rate=16000 ! filesink location=/dev/stdout | livecaption
```

In order to discover your recording device you may use the `gst-device-monitor-1.0` command line tool. For example:

```bash
$ gst-device-monitor-1.0
Probing devices...


Device found:

	name  : Built-in Output
	class : Audio/Sink
	caps  : audio/x-raw, format=(string)F32LE, layout=(string)interleaved, rate=(int)44100, channels=(int)2, channel-mask=(bitmask)0x0000000000000003;
	        audio/x-raw, format=(string){ S8, U8, S16LE, S16BE, U16LE, U16BE, S24_32LE, S24_32BE, U24_32LE, U24_32BE, S32LE, S32BE, U32LE, U32BE, S24LE, S24BE, U24LE, U24BE, S20LE, S20BE, U20LE, U20BE, S18LE, S18BE, U18LE, U18BE, F32LE, F32BE, F64LE, F64BE }, layout=(string)interleaved, rate=(int)[ 1, 2147483647 ], channels=(int)2, channel-mask=(bitmask)0x0000000000000003;
	        audio/x-raw, format=(string){ S8, U8, S16LE, S16BE, U16LE, U16BE, S24_32LE, S24_32BE, U24_32LE, U24_32BE, S32LE, S32BE, U32LE, U32BE, S24LE, S24BE, U24LE, U24BE, S20LE, S20BE, U20LE, U20BE, S18LE, S18BE, U18LE, U18BE, F32LE, F32BE, F64LE, F64BE }, layout=(string)interleaved, rate=(int)[ 1, 2147483647 ], channels=(int)1;
	gst-launch-1.0 ... ! osxaudiosink device=46


Device found:

	name  : Built-in Microph
	class : Audio/Source
	caps  : audio/x-raw, format=(string)F32LE, layout=(string)interleaved, rate=(int)44100, channels=(int)2, channel-mask=(bitmask)0x0000000000000003;
	        audio/x-raw, format=(string){ S8, U8, S16LE, S16BE, U16LE, U16BE, S24_32LE, S24_32BE, U24_32LE, U24_32BE, S32LE, S32BE, U32LE, U32BE, S24LE, S24BE, U24LE, U24BE, S20LE, S20BE, U20LE, U20BE, S18LE, S18BE, U18LE, U18BE, F32LE, F32BE, F64LE, F64BE }, layout=(string)interleaved, rate=(int)44100, channels=(int)2, channel-mask=(bitmask)0x0000000000000003;
	        audio/x-raw, format=(string){ S8, U8, S16LE, S16BE, U16LE, U16BE, S24_32LE, S24_32BE, U24_32LE, U24_32BE, S32LE, S32BE, U32LE, U32BE, S24LE, S24BE, U24LE, U24BE, S20LE, S20BE, U20LE, U20BE, S18LE, S18BE, U18LE, U18BE, F32LE, F32BE, F64LE, F64BE }, layout=(string)interleaved, rate=(int)44100, channels=(int)1;
	gst-launch-1.0 osxaudiosrc device=39 ! ...
```

In the above example the recording device (`Built-In Microphone`) is `osxaudiosrc device=39`, so in order to run the example you would need to adapt the command-line accordingly:

```bash
gst-launch-1.0 -v osxaudiosrc device=39 ! audioconvert ! audioresample ! audio/x-raw,channels=1,rate=16000 ! filesink location=/dev/stdout | livecaption
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