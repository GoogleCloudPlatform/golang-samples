# Google Cloud Natural Language API Go example

## Authentication

* Follow the [instructions][project] to set up your project and enable the Cloud Natural Language API.
* From the Cloud Console, create a service account,
  download its json credentials file, then set the 
  `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

  ```bash
  export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-project-credentials.json
  ```

[cloud-console]: https://console.cloud.google.com
[project]: https://cloud.google.com/natural-language/docs/getting-started#set_up_your_project

## Run the sample

```bash
go build && ./analyze_v2 <command> <text>
```

Where `command` is `entities`, `sentiment`, or `classify`.

For example:

```bash
go build && ./analyze_v2 entities "Renee French designed the Go gopher."
```

Prints something like this:

```
entities: <
  name: "Go"
  type: OTHER
  mentions: <
    text: <
      content: "Go"
      begin_offset: 26
    >
    type: PROPER
    probability: 0.444
  >
>
entities: <
  name: "Renee French"
  type: PERSON
  mentions: <
    text: <
      content: "Renee French"
    >
    type: PROPER
    probability: 0.927
  >
>
entities: <
  name: "gopher"
  type: CONSUMER_GOOD
  mentions: <
    text: <
      content: "gopher"
      begin_offset: 29
    >
    type: PROPER
    probability: 0.441
  >
>
language_code: "en"
language_supported: true
```
