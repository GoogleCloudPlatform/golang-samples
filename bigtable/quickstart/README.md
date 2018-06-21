# Cloud Bigtable Quickstart in Go

This application demonstrates using the [Google Cloud APIs Go
Client Library](https://github.com/GoogleCloudPlatform/google-cloud-go) to connect
to a [Cloud Bigtable](https://cloud.google.com/bigtable/) instance and read a row from a table.

## Before you begin

1. [Select or create][projects] a Cloud Platform project.

1. Enable [billing][billing] for your project.

1. Enable the [Cloud Bigtable API][enable_api].

    Note: The quickstart performs an operation on an existing table (created below).
    If you require your code to create instances or tables,
    the [Admin API](https://console.cloud.google.com/flows/enableapi?apiid=bigtableadmin.googleapis.com)
    must be enabled as well.

1. [Set up authentication with a service account][auth] so you can access the API from your local workstation.

1. Follow the instructions in the [user documentation](https://cloud.google.com/bigtable/docs/creating-instance) to
create a Cloud Bigtable instance (if necessary).

1. Follow the [cbt tutorial](https://cloud.google.com/bigtable/docs/quickstart-cbt) to install the
`cbt` command line tool.
Here are the `cbt` commands to create a table, column family and add some data:
```
   cbt createtable my-table
   cbt createfamily my-table cf1
   cbt set my-table r1 cf1:c1=test-value
```

[projects]: https://console.cloud.google.com/project
[billing]: https://support.google.com/cloud/answer/6293499#enable-billing
[enable_api]: https://console.cloud.google.com/flows/enableapi?apiid=bigtable.googleapis.com
[auth]: https://cloud.google.com/docs/authentication/getting-started


## Running the quickstart

The [quickstart](main.go) sample shows how to read rows from a table with the Bigtable client library.

Run the quickstart to read the row you just wrote using `cbt`:
```
   `go run main.go -project PROJECT_ID -instance INSTANCE_ID -table my-table
```
Expected output similar to:
```
2018/06/15 18:50:52 Getting a single row by row key:
2018/06/15 18:50:54 Row key: r1
2018/06/15 18:50:54 Data: test-value
```

## Cleaning up

To avoid incurring extra charges to your Google Cloud Platform account, remove
the resources created for this sample.

1.  Go to the [Cloud Bigtable instance page](https://console.cloud.google.com/project/_/bigtable/instances) in the Cloud Console.

1.  Click on the instance name.

1.  Click **Delete instance**.

    ![Delete](https://cloud.google.com/bigtable/img/delete-quickstart-instance.png)

1. Type the instance ID, then click **Delete** to delete the instance.
