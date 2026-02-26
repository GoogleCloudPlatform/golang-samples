# Running the Custom AWS Credential Supplier Sample (Go)

This sample demonstrates how to use a custom AWS security credential supplier to authenticate with Google Cloud using AWS as an external identity provider. It uses the **AWS SDK for Go v2** to fetch credentials from sources like Amazon Elastic Kubernetes Service (EKS) with IAM Roles for Service Accounts (IRSA), Elastic Container Service (ECS), or Fargate.

## Prerequisites

*   An AWS account.
*   A Google Cloud project with the IAM API enabled.
*   A Google Cloud Storage bucket.
*   **Go 1.21** or later installed.

If you want to use AWS security credentials that cannot be retrieved using methods supported natively by the Google Auth library, a custom `AwsSecurityCredentialsProvider` implementation may be specified.

## Running Locally

For local development, you can provide credentials and configuration in a JSON file.

### 1. Install Dependencies

Initialize the module and download required packages:

```bash
go mod tidy
```

### 2. Configure Credentials

1.  Copy the example secrets file to a new file named `custom-credentials-aws-secrets.json` in the project root:
    ```bash
    cp custom-credentials-aws-secrets.json.example custom-credentials-aws-secrets.json
    ```
2.  Open `custom-credentials-aws-secrets.json` and fill in the required values for your AWS and Google Cloud configuration. Do not check your `custom-credentials-aws-secrets.json` file into version control.

**Note:** Do not check your secrets file into version control. 

### 3. Run the Application

Execute the Go program:

```bash
go run .
```

The application will detect the `custom-credentials-aws-secrets.json` file, use the AWS SDK to resolve credentials, exchange them for a Google Cloud token, and retrieve metadata for your GCS bucket using the Google Cloud Storage Client Library.

## Running in a Containerized Environment (EKS)

This section provides a brief overview of how to run the sample in an Amazon EKS cluster.

### 1. EKS Cluster Setup

First, you need an EKS cluster. You can create one using `eksctl` or the AWS Management Console. For detailed instructions, refer to the [Amazon EKS documentation](https://docs.aws.amazon.com/eks/latest/userguide/create-cluster.html).

### 2. Configure IAM Roles for Service Accounts (IRSA)

IRSA allows you to associate an IAM role with a Kubernetes service account. This provides a secure way for your pods to access AWS services without hardcoding long-lived credentials.

Run the following command to create the IAM role and bind it to a Kubernetes Service Account:

```bash
eksctl create iamserviceaccount \
  --name your-k8s-service-account \
  --namespace default \
  --cluster your-cluster-name \
  --region your-aws-region \
  --role-name your-role-name \
  --attach-policy-arn arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess \
  --approve
```

> **Note**: The `--attach-policy-arn` flag is used here to demonstrate attaching permissions. Update this with the specific AWS policy ARN your application requires.

For detailed steps, see the [IAM Roles for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) documentation.

### 3. Configure Google Cloud to Trust the AWS Role

You need to configure your Google Cloud project to trust the AWS IAM role you created.

1.  **Create a Workload Identity Pool and Provider** that trusts your AWS account.
2.  **Bind the AWS Role to a Google Cloud Service Account** (or grant permissions directly to the federated identity).

For detailed steps, see [Configuring Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation-with-other-clouds).

### 4. Containerize and Package the Application

Create a `Dockerfile` for the Go application. See the [`Dockerfile`](Dockerfile) file for an example.

Build and push the image to your registry (e.g., Amazon ECR):

```bash
docker build -t your-container-image:latest .
docker push your-container-image:latest
```

### 5. Deploy to EKS

Create a Kubernetes deployment manifest to deploy your application to the EKS cluster. See the [`pod.yaml`](pod.yaml) file for an example.

**Note:** The provided [`pod.yaml`](pod.yaml) is an example and may need to be modified for your specific needs.

Deploy the pod:

```bash
kubectl apply -f pod.yaml
```

### 6. Clean Up

To clean up the resources, delete the EKS cluster and any other AWS and GCP resources you created.

```bash
eksctl delete cluster --name your-cluster-name
```
