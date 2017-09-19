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
go build && ./analyze <command> <text>
```

Where `command` is `entities`, `sentiment`, or `syntax`.

For example:

```bash
go build && ./analyze entities "Renee French designed the Go gopher."
```

Prints something like this:

```
entities: <
  name: "Renee French"
  type: PERSON
  salience: 0.4693242
  mentions: <
    text: <
      content: "Renee French"
    >
  >
>
entities: <
  name: "Go"
  type: ORGANIZATION
  metadata: <
    key: "wikipedia_url"
    value: "http://en.wikipedia.org/wiki/Go_(programming_language)"
  >
  salience: 0.34126133
  mentions: <
    text: <
      content: "Go"
      begin_offset: 26
    >
  >
>
language: "en"
```
