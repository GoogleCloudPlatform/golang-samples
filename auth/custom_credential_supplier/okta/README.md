# Running the Custom Okta Credential Supplier Sample

This document provides instructions on how to run the custom Okta credential supplier sample.

## 1. Okta Configuration

Before running the sample, you need to configure an Okta application for Machine-to-Machine (M2M) communication.

### Create an M2M Application in Okta

1.  Log in to your Okta developer console.
2.  Navigate to **Applications** > **Applications** and click **Create App Integration**.
3.  Select **API Services** as the sign-on method and click **Next**.
4.  Give your application a name and click **Save**.

### Obtain Okta Credentials

Once the application is created, you will find the following information in the **General** tab:

*   **Okta Domain**: Your Okta developer domain (e.g., `dev-123456.okta.com`).
*   **Client ID**: The client ID for your application.
*   **Client Secret**: The client secret for your application.

You will need these values to configure the sample.

## 2. GCP Configuration

You need to configure a Workload Identity Pool in GCP to trust the Okta application.

### Set up Workload Identity Federation

1.  In the Google Cloud Console, navigate to **IAM & Admin** > **Workload Identity Federation**.
2.  Click **Create Pool** to create a new Workload Identity Pool.
3.  Add a new **OIDC provider** to the pool.
4.  Configure the provider with your Okta domain as the issuer URL.
5.  Map the Okta `sub` (subject) assertion to a GCP principal.

For detailed instructions, refer to the [Workload Identity Federation documentation](https://cloud.google.com/iam/docs/workload-identity-federation).

### GCS Bucket

Ensure you have a GCS bucket that the authenticated user will have access to. You will need the name of this bucket to run the sample.

## 3. Running the Script

To run the sample, set the following environment variables:

```bash
export OKTA_DOMAIN="your-okta-domain"
export OKTA_CLIENT_ID="your-okta-client-id"
export OKTA_CLIENT_SECRET="your-okta-client-secret"
export GCP_WORKLOAD_IDENTITY_POOL="your-gcp-workload-identity-pool"
export GCP_WORKLOAD_IDENTITY_PROVIDER="your-gcp-workload-identity-provider"
export GCS_BUCKET_NAME="your-gcs-bucket-name"
export GOOGLE_CLOUD_PROJECT="your-google-cloud-project"

go run .
```

The script will then authenticate with Okta, exchange the Okta token for a GCP token, and use the GCP token to list the objects in the specified GCS bucket.
