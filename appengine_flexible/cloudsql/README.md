## Setup

Before you can run or deploy the sample, you will need to create a Cloud SQL database:

1. Create a Cloud SQL instance. You can do this from the [Google Developers
    Console](https://console.developers.google.com) or via the
    [Cloud SDK](https://cloud.google.com/sdk). To create it via the SDK,
    use the following command:

        $ gcloud sql instances create [your-instance-name] \
            --assign-ip \
            --authorized-networks 0.0.0.0/0 \
            --tier D0

2. Create a new user and database for the application. The easiest way to do
    this is via the [Google Developers Console](https://console.developers.google.com/project/_/sql/instances/).
    Alternatively, you can use MySQL tools such as the command line client or
    workbench.

3. Update the connection string in `app.yaml` with the username, password,
    hostname and database name of the Cloud SQL instance you just created.

## Running locally

To run locally, set the environment variables before running the sample:

    $ export MYSQL_CONNECTION=user:password@tcp([host]:3306)/dbname
    $ go run cloudsql.go
