<img src="https://avatars2.githubusercontent.com/u/2810941?v=3&s=96" alt="Google Cloud Platform logo" title="Google Cloud Platform" align="right" height="96" width="96"/>

# Google Cloud Data Loss Prevention (DLP) API: Go Samples

[![Open in Cloud Shell][shell_img]][shell_link]

The [Data Loss Prevention API](https://cloud.google.com/dlp/docs/) provides programmatic access to a powerful detection engine for personally identifiable information and other privacy-sensitive data in unstructured data streams.

## Table of Contents

* [Before you begin](#before-you-begin)
* [Samples](#samples)
  * [De-identify](#de-identify)
  * [Inspect](#inspect)
  * [Metadata](#metadata)
  * [Redact](#redact)
  * [Risk Analysis](#risk-analysis)
  * [Templates](#templates)
  * [Triggers](#triggers)

## Before you begin

Before running the samples, make sure you've:

1.  Enabled the [DLP API](https://console.developers.google.com/apis/api/dlp.googleapis.com/overview).
1.  Enabled the [PubSub API](https://console.developers.google.com/apis/api/pubsub.googleapis.com/overview).
1.  Set up [authentication](https://cloud.google.com/docs/authentication/getting-started).
1.  Installed the sample application by running:
    ```bash
    go get -u github.com/GoogleCloudPlatform/golang-samples/dlp
    ```

## Samples

Usage: `./dlp -project <my-project> [options] subcommand [args]`

```
Options:
  -bytesType value
    	Bytes type of input file for inspectFile and redactImage [IMAGE_SVG, TEXT_UTF8, BYTES_TYPE_UNSPECIFIED, IMAGE_JPEG, IMAGE_BMP, IMAGE_PNG] (default BYTES_TYPE_UNSPECIFIED)
  -includeQuote
    	Include a quote of findings for inspect* (default false)
  -infoTypes string
    	Info types to inspect*, redactImage, createTrigger, and createInspectTemplate (default "PHONE_NUMBER,EMAIL_ADDRESS,CREDIT_CARD_NUMBER,US_SOCIAL_SECURITY_NUMBER")
  -languageCode string
    	Language code for infoTypes (default "en-US")
  -maxFindings int
    	Number of results for inspect*, createTrigger, and createInspectTemplate (default 0 (no limit))
  -minLikelihood value
    	Minimum likelihood value for inspect*, redactImage, createTrigger, and createInspectTemplate [LIKELY, VERY_LIKELY, LIKELIHOOD_UNSPECIFIED, VERY_UNLIKELY, UNLIKELY, POSSIBLE] (default LIKELIHOOD_UNSPECIFIED)
  -project string
    	GCloud project ID (required)
```

Subcommands and their args are described below.

### De-identify

View the [source code](deid.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/deid.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
go build
./dlp -project <project> [options] dateShift <string>
./dlp -project <project> [options] fpe <string> <wrappedKeyFileName> <cryptoKeyname> <surrogateInfoType>
./dlp -project <project> [options] mask <string>
./dlp -project <project> [options] reidentifyFPE <string> <wrappedKeyFileName> <cryptoKeyname> <surrogateInfoType>
```

__Examples:__
```
./dlp -project my-project dateShift "My birthday is January 1, 1970"
./dlp -project my-project mask "My SSN is 111222333"
ENC=$(./dlp -project my-project fpe "My SSN is 111222333" key.enc projects/my-project/locations/global/keyRings/my-key-ring/cryptoKeys/my-key randomstring)
./dlp -project my-project reidentifyFPE "$ENC" key.enc projects/my-project/locations/global/keyRings/my-key-ring/cryptoKeys/my-key randomstring
```

For more information, see https://cloud.google.com/dlp/docs.

### Inspect

View the [source code](inspect.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/inspect.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
go build
./dlp -project <project> [options] inspect <string>
./dlp -project <project> [options] inspectBigquery <pubSubTopic> <pubSubSub> <dataProject> <datasetID> <tableID>
./dlp -project <project> [options] inspectDatastore <pubSubTopic> <pubSubSub> <dataProject> <namespaceID> <kind>
./dlp -project <project> [options] inspectFile <filename>
./dlp -project <project> [options] inspectGCSFile <pubSubTopic> <pubSubSub> <bucketName> <fileName>
```

__Examples:__
```
./dlp -project my-project inspect "My SSN is 111222333 and my phone number is (123) 456-7890"
./dlp -project my-project inspectBigquery inspect-topic inspect-sub dataProject datasetID tableID
./dlp -project my-project inspectDatastore inspect-topic inspect-sub my-data-project my-namespace my-kind
./dlp -project my-project inspectFile my-file
./dlp -project my-project inspectGCSFile inspect-topic inspect-sub my-bucket my-file
```

For more information, see https://cloud.google.com/dlp/docs.

### Jobs

View the [source code](jobs.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp_snippets/jobs.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
go build
./dlp -project <project> [options] deleteJob <jobID>
./dlp -project <project> [options] listJobs <filter> <jobType>
```

__Examples:__
```
./dlp -project my-project inspect "My SSN is 111222333 and my phone number is (123) 456-7890"
    ./dlp -project my-project deleteJob /projects/my-project/dlpJobs/my-job
    ./dlp -project my-project listJobs "" ""

```

For more information, see https://cloud.google.com/dlp/docs.

### Metadata

View the [source code](metadata.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/metadata.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
go build
./dlp -project <project> [options] infoTypes <filter>
```

__Examples:__
```
./dlp -project my-project infoTypes supported_by=INSPECT
```

For more information, see https://cloud.google.com/dlp/docs.

### Redact

View the [source code](redact.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/redact.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
go build
./dlp -project <project> [options] redactImage <inputPath> <outputPath>
```

__Examples:__
```
./dlp -project my-project -bytesType IMAGE_PNG redactImage input.png output.png
```

For more information, see https://cloud.google.com/dlp/docs.

### Risk Analysis

View the [source code](risk.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/risk.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
go build
./dlp -project <project> [options] riskCategorical <dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <columnName>
./dlp -project <project> [options] riskKAnonymity  <dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <column,names>
./dlp -project <project> [options] riskKMap        <dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <region> <column,names>
./dlp -project <project> [options] riskLDiversity  <dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <sensitiveAttribute> <column,names>
./dlp -project <project> [options] riskNumerical   <dataProject> <pubSubTopic> <pubSubSub> <datasetID> <tableID> <columnName>
```

__Examples:__
```
./dlp -project my-project riskNumerical   bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
./dlp -project my-project riskCategorical bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number
./dlp -project my-project riskKAnonymity  bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 state_number,county
./dlp -project my-project riskLDiversity  bigquery-public-data risk-topic risk-sub nhtsa_traffic_fatalities accident_2015 city state_number,county
./dlp -project my-project riskKMap        bigquery-public-data risk-topic risk-sub san_francisco bikeshare_trips USA zip_code
```

For more information, see https://cloud.google.com/dlp/docs.

### Templates

View the [source code](templates.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/templates.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
go build
./dlp -project <project> [options] createInspectTemplate <templateID> <displayName> <description>
./dlp -project <project> [options] deleteInspectTemplate <fullTemplateID>
./dlp -project <project> [options] listInspectTemplates
```

__Examples:__
```bash
./dlp -project my-project createInspectTemplate my-template "My Template" "My template description"
./dlp -project my-project deleteInspectTemplate projects/my-project/inspectTemplates/my-template
./dlp -project my-project listInspectTemplates
```

For more information, see https://cloud.google.com/dlp/docs.


### Triggers
View the [source code](triggers.go).

[![Open in Cloud Shell][shell_img]](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/triggers.go,dlp/dlp_snippets/README.md)

__Usage:__
```bash
./dlp -project <project> [options] createTrigger <triggerID> <displayName> <description> <bucketName>
./dlp -project <project> [options] deleteTrigger <fullTriggerID>
./dlp -project <project> [options] listTriggers
```

__Examples:__
```bash
./dlp -project my-project createTrigger my-trigger "My Trigger" "My trigger description" my-bucket
./dlp -project my-project deleteTrigger projects/my-project/jobTriggers/my-trigger
./dlp -project my-project listTriggers
```

[shell_img]: http://gstatic.com/cloudssh/images/open-btn.png
[shell_link]: https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dlp/dlp_snippets/README.md