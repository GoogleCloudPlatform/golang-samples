<img src="https://avatars2.githubusercontent.com/u/2810941?v=3&s=96" alt="Google
Cloud Platform logo" title="Google Cloud Platform" align="right" height="96"
width="96"/>

# Google Cloud Job Discovery API Samples

Cloud Job Discovery is part of Google for Jobs - a Google-wide commitment to help
people find jobs more easily. Job Discovery provides plug and play access to 
Google’s search and machine learning capabilities, enabling the entire recruiting
ecosystem - company career sites, job boards, applicant tracking systems, and
staffing agencies to improve job site engagement and candidate conversion.

## Installation
```shell
go get -u google.golang.org/api/jobs/v2
```

## Prerequisite
1.  **Enable APIs** - [Enable the Job Discovery API](https://console.cloud.google.com/flows/enableapi?apiid=jobs.googleapis.com)
    and create a new project or select an existing project.
2.  **Activate your Credentials** - If you do not already have an active set of credentials, create and download a [JSON Service Account key](https://pantheon.corp.google.com/apis/credentials/serviceaccountkey). Set the environment variable `GOOGLE_APPLICATION_CREDENTIALS` as the path to the downloaded JSON file.
```
export GOOGLE_APPLICATION_CREDENTIALS=/PATH/TO/YOUR/key.json
```

## Run the Samples
For example, to run the quickstart:
```shell
go run quickstart/main.go
```
