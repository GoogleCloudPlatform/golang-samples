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
go run analyze.go <command> <text>
```

Where `command` is `entities`, `sentiment`, or `syntax`.

For example:

```bash
go run analyze.go entities "Renee French designed the Go gopher."
```

Prints something like this:

```json
{
    "entities": [
        {
            "mentions": [
                {
                    "text": {
                        "content": "Renee French"
                    }
                }
            ],
            "name": "Renee French",
            "salience": 0.51491237,
            "type": "PERSON"
        },
        {
            "mentions": [
                {
                    "text": {
                        "beginOffset": 26,
                        "content": "Go"
                    }
                }
            ],
            "metadata": {
                "wikipedia_url": "http://en.wikipedia.org/wiki/Go_(programming_language)"
            },
            "name": "Go",
            "salience": 0.31268582,
            "type": "OTHER"
        }
    ],
    "language": "en"
}
```
