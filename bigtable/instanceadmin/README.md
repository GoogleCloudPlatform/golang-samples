# Cloud Bigtable Instance Admin in Go

This is a simple application that demonstrates using the [Google Cloud APIs Go
Client Library](https://github.com/GoogleCloudPlatform/google-cloud-go) to perform
Cloud Bigtable instance administration tasks, specifically creating an instance with optional tags.

## Prerequisites

1.  **Set up Cloud Console:**
    1.  Go to the [Cloud Console](https://console.cloud.google.com/) and create or select your project. You will need the project ID later.
    2.  Ensure billing is enabled for your project.
    3.  Enable the **Cloud Bigtable API** and the **Cloud Bigtable Admin API**. You can find these by navigating to **APIs & Services > Library**.

2.  **Set up gcloud:**
    1.  Install and initialize the Google Cloud SDK.
    2.  Update your components: `gcloud components update`
    3.  Authenticate: `gcloud auth login`
    4.  Set your default project: `gcloud config set project PROJECT_ID` (replace `PROJECT_ID` with your actual project ID).

3.  **(Optional) Create Tags:**
    *   The sample code includes an option to attach tags to the Bigtable instance upon creation.
    *   If you want to use this feature, you MUST create the tag keys and tag values in your GCP organization/project *before* running the code.
    *   For detailed instructions on creating and managing tags, please see the [Tags documentation](https://cloud.google.com/bigtable/docs/tags). The code has placeholder `tagKey` and `tagValue` variables that you will need to update with your actual tag IDs.

## Running

1.  Navigate to the `instanceadmin` directory.

2.  Execute the program using `go run`:

    ```bash
    go run main.go -project YOUR_PROJECT_ID -instance YOUR_INSTANCE_ID -zone YOUR_ZONE
    ```
    Substitute `YOUR_PROJECT_ID`, `YOUR_INSTANCE_ID`, and `YOUR_ZONE` with your desired values.

3.  **To Create an Instance WITH Tags:**
    *   Edit `main.go`.
    *   Uncomment the lines defining `tagKey` and `tagValue`, and update them with your pre-existing tag key and value IDs.
    *   Uncomment the `Tags: map[string]string{tagKey: tagValue},` line within the `bigtable.InstanceConf` struct.
    *   Run the `go run` command as shown above.

## Cleaning up

To avoid incurring extra charges to your Google Cloud Platform account, delete the resources created for this sample. The primary resource to delete is the Cloud Bigtable instance.

**Option 1: Using Cloud Console**

1.  Go to the Bigtable instances page in the [Cloud Console](https://console.cloud.google.com/bigtable/instances).
2.  Select the instance you created.
3.  Click the **Delete Instance** button.
4.  Confirm the deletion by typing the instance ID.

**Option 2: Using gcloud**

1.  Run the following command, replacing `YOUR_INSTANCE_ID` and `YOUR_PROJECT_ID`:

    ```bash
    gcloud bigtable instances delete YOUR_INSTANCE_ID --project YOUR_PROJECT_ID
    ```
    You will be prompted to confirm the deletion.
