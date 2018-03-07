<img src="https://avatars2.githubusercontent.com/u/2810941?v=3&s=96" alt="Google Cloud Platform logo" title="Google Cloud Platform" align="right" height="96" width="96"/>

# Google Cloud Data Loss Prevention (DLP) API: Go Samples

[![Open in Cloud Shell][shell_img]][shell_link]

The [Data Loss Prevention API](https://cloud.google.com/dlp/docs/) provides programmatic access to a powerful detection engine for personally identifiable information and other privacy-sensitive data in unstructured data streams.

## Table of Contents

* [Before you begin](#before-you-begin)
* [Samples](#samples)
  * [Inspect](#inspect)
  * [Redact](#redact)
  * [Metadata](#metadata)
  * [DeID](#deid)
  * [Risk Analysis](#risk-analysis)

## Before you begin

Before running the samples, make sure you've enabled the
[DLP API](https://console.developers.google.com/apis/api/dlp.googleapis.com/overview),
enabled the [PubSub API](https://console.developers.google.com/apis/api/pubsub.googleapis.com/overview)
(for the risk examples), and set up [authentication](https://cloud.google.com/docs/authentication/getting-started).

## Samples

### Inspect

View the [source code][inspect_0_code].

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/inspect.go,dlp/README.md)

__Usage:__
```bash
go build
./dlp -project <project> inspect <string>
```

```
Examples:
  ./dlp -project my-project inspect "My SSN is 111222333"

For more information, see https://cloud.google.com/dlp/docs. Full options are explained at
https://cloud.google.com/dlp/docs/reference/rest/v2beta1/content/inspect#InspectConfig
```

[inspect_0_docs]: https://cloud.google.com/dlp/docs
[inspect_0_code]: inspect.go

### Redact

View the [source code][redact_1_code].

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/redact.go,dlp/README.md)

__Usage:__
```bash
go build
./dlp -project <project> redact <string>
```

```
Examples:
  ./dlp -project my-project redact "My SSN is 111222333"

For more information, see https://cloud.google.com/dlp/docs. Full options are explained at
https://cloud.google.com/dlp/docs/reference/rest/v2beta1/content/redact.
```

[redact_1_docs]: https://cloud.google.com/dlp/docs
[redact_1_code]: redact.go

### Metadata

View the [source code][metadata_2_code].

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/metadata.go,dlp/README.md)

__Usage:__
```bash
go build
./dlp -project <project> infoTypes <filter>
```

```
Examples:
  ./dlp -project my-project infoTypes supported_by=INSPECT

For more information, see https://cloud.google.com/dlp/docs
```

[metadata_2_docs]: https://cloud.google.com/dlp/docs
[metadata_2_code]: metadata.go

### DeID

View the [source code][deid_3_code].

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/deid.go,dlp/README.md)

__Usage:__
```bash
go build
./dlp -project <project> mask <string>
./dlp -project <project> fpe <string> <wrappedKey> <keyName>
```

```
Examples:
  ./dlp -project my-project mask "My SSN is 372819127"
  ./dlp -project my-project fpe "My SSN is 372819127" <YOUR_ENCRYPTED_AES_256_KEY> <YOUR_KEY_NAME>

For more information, see https://cloud.google.com/dlp/docs.
```

[deid_3_docs]: https://cloud.google.com/dlp/docs
[deid_3_code]: deid.go

### Risk Analysis

View the [source code][risk_4_code].

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/risk.go,dlp/README.md)

__Usage:__
```bash
go build
./dlp -project <project> riskNumerical   <dataProject> <PubSubTopicName> <PubSubSubscriptionName> <datasetID> <tableID> <columnName>
./dlp -project <project> riskCategorical <dataProject> <PubSubTopicName> <PubSubSubscriptionName> <datasetID> <tableID> <columnName>
./dlp -project <project> riskKAnonymity  <dataProject> <PubSubTopicName> <PubSubSubscriptionName> <datasetID> <tableID> <commaSepColumnNames>
./dlp -project <project> riskLDiversity  <dataProject> <PubSubTopicName> <PubSubSubscriptionName> <datasetID> <tableID> <sensitiveColumnName> <commaSepColumnNames>
./dlp -project <project> riskKMap        <dataProject> <PubSubTopicName> <PubSubSubscriptionName> <datasetID> <tableID> <region> <columnName>
```

```
Examples:
  ./dlp -project my-project riskNumerical   bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
  ./dlp -project my-project riskCategorical bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
  ./dlp -project my-project riskKAnonymity  bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number,county
  ./dlp -project my-project riskLDiversity  bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 city state_number,county
  ./dlp -project my-project riskKMap        bigquery-public-data risk-topic risk-sub san_francisco bikeshare_trips USA zip_code

For more information, see https://cloud.google.com/dlp/docs.
```

[risk_4_docs]: https://cloud.google.com/dlp/docs
[risk_4_code]: risk.go

[shell_img]: http://gstatic.com/cloudssh/images/open-btn.png
[shell_link]: https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/README.md