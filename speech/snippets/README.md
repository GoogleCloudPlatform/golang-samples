# Cloud Speech API Go samples

This directory contains [Cloud Speech API](https://cloud.google.com/speech/) Go snippets.

## Snippets

### Auto Punctuation

For more details, see the [Getting Punctuation](https://cloud.google.com/speech-to-text/docs/automatic-punctuation) tutorial in the docs.

[Go Code](auto_punctuation.go)

### Speech Adaptation

For more details see the [Using Speech Adaptation](https://cloud.google.com/speech-to-text/docs/context-strength) tutorial in the docs.

[Go Code](context_classes.go)

### Enhanced Model

For more details, see the [Using Enhanced Models](https://cloud.google.com/speech-to-text/docs/enhanced-models) tutorial in the docs.

> **Caution**: If you attempt to use an enhanced model but your Google Cloud Project does not have data logging enabled, Speech-to-Text API sends a `400` HTTP code response with the status `INVALID_ARGUMENT`. You must [enable data logging](https://cloud.google.com/speech-to-text/docs/enable-data-logging) to use the enhanced speech recognition models.

[Go Code](enhanced_model.go)

### Model Selection

For more details, see the [Transcribing Video Files](https://cloud.google.com/speech-to-text/docs/video-model) tutorial in the docs.

### Multi-channel Transcription

For more details, see the [Transcribing Audio with Multiple Channels](https://cloud.google.com/speech-to-text/docs/multi-channel) tutorial in the docs.

[Go Code](multichannel.go)
