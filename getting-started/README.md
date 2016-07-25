# Getting started with Go on the Google Cloud Platform

This repository contains sample code for the [Go Getting Started on Google Cloud Platform][gogs]

Please refer to the guide for full instructions on how to run the samples.

## Checking out the code

    $ go get -d github.com/GoogleCloudPlatform/golang-samples/getting-started

## Install aedeploy

    $ go get google.golang.org/appengine/cmd/aedeploy

## Run and deploy "Hello world"

    $ cd $GOPATH/src/github.com/GoogleCloudPlatform/golang-samples/getting-started/helloworld
    $ go run helloworld.go
    $ aedeploy gcloud app deploy

## Run and deploy "Bookshelf"

    $ cd $GOPATH/src/github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf/app
    $ vim ../config.go
    $ go run *.go
    $ aedeploy gcloud app deploy

## Run and deploy "Bookshelf pub/sub worker"

    $ cd $GOPATH/src/github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf/pubsub_worker
    $ vim ../config.go
    $ go run *.go
    $ aedeploy gcloud app deploy

## Contributing

See [CONTRIBUTING.md](/CONTRIBUTING.md)

## License

The source code in this repository is available under the [Apache 2.0 license](/LICENSE).

[gogs]: https://cloud.google.com/go
