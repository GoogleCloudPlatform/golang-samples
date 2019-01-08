# Contributing

1. Sign one of the contributor license agreements below.
1. [Install Go](https://golang.org/doc/install)
1. Get the package:

    `go get -d github.com/GoogleCloudPlatform/golang-samples`
1. Change into the checked out source:

    `cd $(go env GOPATH)/src/github.com/GoogleCloudPlatform/golang-samples`
1. Fork the repo.
1. Set your fork as a remote:

    `git remote add fork git@github.com:GITHUB_USERNAME/golang-samples.git`
1. Make changes, commit to your fork.
1. Send a pull request with your changes.

# Formatting

All code must be formatted with `gofmt` (with the latest Go version) and pass
`go vet`.

# Testing

Tests are required for all samples. When writing a pull request, be sure to
write and run the tests in any modified directories.

## Running system tests

1. Set the `GOLANG_SAMPLES_PROJECT_ID` environment variable to a suitable test project.
1. Ensure you are logged in using `gcloud login` or set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to the path of your credentials file.
   Tests are authenticated using [Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials).
1. Install the test dependencies:

    `go get -t -d github.com/GoogleCloudPlatform/golang-samples/...`
1. Run the tests:

    `go test github.com/GoogleCloudPlatform/golang-samples/...`

Note: You may want to `cd` to the directory you're modifying and run
`go test -v ./...` to avoid running every test in the repo.

## Contributor License Agreements

Before we can accept your pull requests you'll need to sign a Contributor
License Agreement (CLA):

- **If you are an individual writing original source code** and **you own the
intellectual property**, then you'll need to sign an [individual CLA][indvcla].
- **If you work for a company that wants to allow you to contribute your work**,
then you'll need to sign a [corporate CLA][corpcla].

You can sign these electronically (just scroll to the bottom). After that,
we'll be able to accept your pull requests.

[gcloudcli]: https://developers.google.com/cloud/sdk/gcloud/
[indvcla]: https://developers.google.com/open-source/cla/individual
[corpcla]: https://developers.google.com/open-source/cla/corporate
