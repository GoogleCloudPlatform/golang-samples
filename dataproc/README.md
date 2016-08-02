# Google Cloud Dataproc Go sample

This directory contains a Go sample for
[Google Cloud Dataproc](https://cloud.google.com/dataproc).
This sample does the following to show an example use of Cloud Dataproc:

1. Create a cluster
1. Submit a PySpark job to the cluster
1. Get a link to the job output once the job completes
1. Delete the cluster

## Using this sample
To use this sample, you will need:

* An active Google Cloud Platform account
* A project within your account
* A [Google Cloud Storage](https://cloud.google.com/storage) bucket

Once you have those, you can run this sample by doing the following. Please
note, this example will use the included PySpark example `pyspark_sort.py` in
this repository:

1. Set your [GOPATH](https://github.com/golang/go/wiki/GOPATH)
1. Get this repository by running
`go get -d github.com/GoogleCloudPlatform/golang-samples`
4. Change directories to `$GOPATH/src/github.com/evilsoapbox/golang-samples`
5. Install dependancies
    go get -d golang.org/x/net/context
    go get -d golang.org/x/oauth2/google
    go get -d google.golang.org/api/dataproc/v1
    go get -d google.golang.org/api/storage/v1
6. Run the sample script:
`go run dataproc.go -bucket YOUR_BUCKET -pyspark-file pyspark_sort.py
-project YOUR_PROJECT -zone YOUR_ZONE -cluster-name YOUR_CLUSTER_NAME`

In the command above, you will need to specify your own values for several
options, all of which are required:

* `-bucket` - Your Google Cloud Storage bucket for holding the PySpark file
* `-project` - Your Google Cloud Platform project ID
* `-zone` - The zone you wish to use, such as `us-central1-f`
* `-cluster-name` - The name your cluster will be given
