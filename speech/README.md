# Cloud Speech API Go samples

This directory contains [Cloud Speech API](https://cloud.google.com/speech/) Go samples and utilities.

## Samples

### Caption

The `caption` command sends audio data to the Google Speech API and prints its transcript. It returns recognized text for short audio (less than ~1 minute). For long audio, see the Caption Async example.

For more details, see the [Synchronous Speech Recognition](https://cloud.google.com/speech/docs/sync-recognize) tutorial in the docs.

[Go Code](caption)

### Caption Async

The `captionasync` command sends audio data to the Google Speech API and prints its transcript. It uses the Asynchronous Speech Recognition to process audio that is longer than a minute and is stored in the Google Cloud Storage. For shorter audio, including audio stored locally (inline), Synchronous Speech Recognition is faster and simpler.

For more details, see the [Asynchronous Speech Recognition](https://cloud.google.com/speech/docs/async-recognize) tutorial in the docs.

[Go Code](captionasync)

### Enhanced Model

For more details, see the [Using Enhanced Models](https://cloud.google.com/speech-to-text/docs/enhanced-models) tutorial in the docs.

> **Caution**: If you attempt to use an enhanced model but your Google Cloud Project does not have data logging enabled, Speech-to-Text API sends a `400` HTTP code response with the status `INVALID_ARGUMENT`. You must [enable data logging](https://cloud.google.com/speech-to-text/docs/enable-data-logging) to use the enhanced speech recognition models.

[Go Code](enhanced_model)

### Live Caption

The `livecaption` command pipes the stdin audio data to Google Speech API and outputs the transcript. It uses Streaming Speech Recognition to process the stdin audio data and returns results in real time as the audio is processed.

Please note that the audio size limit for a single request is 1 minute. For periods longer than that, audio must be sent at a rate that approximates real time. Refer to the [Limits](https://cloud.google.com/speech/limits) page on the documentation for more details.

For more details, see the [Streaming Speech Recognition](https://cloud.google.com/speech/docs/streaming-recognize) tutorial in the docs.

[Go Code](livecaption)

### Word Offset

The `wordoffset` command sends audio data to the Google Speech API and prints the transcript along with word offset (timestamp) information.

For more details, see the [Time Offsets (Timestamps)](https://cloud.google.com/speech/docs/async-time-offsets) tutorial in the docs.

[Go Code](wordoffset)
