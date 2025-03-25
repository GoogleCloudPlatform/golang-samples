# Google Parameter Manager

<a href="https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/go-docs-samples&page=editor&open_in_editor=parametermanager/README.md">
<img alt="Open in Cloud Shell" src ="http://gstatic.com/cloudssh/images/open-btn.png"></a>

Google [Parameter Manager](https://cloud.google.com/secret-manager/parameter-manager/docs/overview)
provides a centralized storage for all configuration parameters related to your workload deployments.
Parameters are variables, often in the form of key-value pairs, which customize how an application functions.
These sample Go applications demonstrate how to access
the Parameter Manager API using the Google Go API Client Libraries.

## Prerequisites

### Enable the API

You
must [enable the Parameter Manager API](https://console.cloud.google.com/apis/enableflow?apiid=parametermanager.googleapis.com)
for your project in order to use these samples

### Set Environment Variables

You must set your project ID in order to run the tests

```text
$ export GOOGLE_CLOUD_PROJECT=<your-project-id-here>
```

You must set your Location in order to run the regional tests

```text
$ export GOOGLE_CLOUD_PROJECT_LOCATION=<your-location-id-here>
```

### Grant Permissions

You must ensure that
the [user account or service account](https://cloud.google.com/iam/docs/service-accounts#differences_between_a_service_account_and_a_user_account)
you used to authorize your gcloud session has the proper permissions to edit Parameter Manager resources for your project.
In the Cloud Console under IAM, add the following roles to the project whose service account you're using to test:

* Parameter Manager Admin (`roles/parametermanager.admin`)
* Parameter Manager Parameter Accessor (`roles/parametermanager.parameterAccessor`)
* Parameter Manager Parameter Version Adder (`roles/parametermanager.parameterVersionAdder`)

To use the rendering of secret through parameter manager add the following role also:
*  Secret Manager Secret Accessor (`roles/secretmanager.secretAccessor`)

More information can be found in
the [Parameter Manager Docs](https://cloud.google.com/secret-manager/parameter-manager/docs/access-control)