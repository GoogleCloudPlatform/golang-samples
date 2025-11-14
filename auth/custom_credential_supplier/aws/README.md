# Running the Custom Credential Supplier Sample

If you want to use AWS security credentials that cannot be retrieved using methods supported natively by the [google-cloud-go/auth](https://github.com/vverman/google-cloud-go/tree/main/auth) library, a custom AwsSecurityCredentialsProvider implementation may be specified when creating an AWS client. The supplier must return valid, unexpired AWS security credentials when called by the GCP credential. Currently, using ADC with your AWS workloads is only supported with EC2. An example of a good use case for using a custom credential suppliers is when your workloads are running in other AWS environments, such as ECS, EKS, Fargate, etc.


This document provides instructions on how to run the custom credential supplier sample in different environments.

## Running Locally

To run the sample on your local system, you need to configure your AWS and GCP credentials as environment variables.

```bash
export AWS_ACCESS_KEY_ID="YOUR_AWS_ACCESS_KEY_ID"
export AWS_SECRET_ACCESS_KEY="YOUR_AWS_SECRET_ACCESS_KEY"
export GCP_WORKLOAD_AUDIENCE="YOUR_GCP_WORKLOAD_AUDIENCE"
export GCS_BUCKET_NAME="YOUR_GCS_BUCKET_NAME"

# Optional: If you want to use service account impersonation
export GCP_SERVICE_ACCOUNT_IMPERSONATION_URL="YOUR_GCP_SERVICE_ACCOUNT_IMPERSONATION_URL"

go run .
```

## Running in a Containerized Environment (EKS)

This section provides a brief overview of how to run the sample in an Amazon EKS cluster.

### 1. EKS Cluster Setup

First, you need an EKS cluster. You can create one using `eksctl` or the AWS Management Console. For detailed instructions, refer to the [Amazon EKS documentation](https://docs.aws.amazon.com/eks/latest/userguide/create-cluster.html).

### 2. Configure IAM Roles for Service Accounts (IRSA)

IRSA allows you to associate an IAM role with a Kubernetes service account. This provides a secure way for your pods to access AWS services.

- Create an IAM OIDC provider for your cluster.
- Create an IAM role and policy that grants the necessary AWS permissions.
- Associate the IAM role with a Kubernetes service account.

For detailed steps, see the [IAM Roles for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) documentation.

### 3. Configure GCP to Trust the AWS Role

You need to configure your GCP project to trust the AWS IAM role you created. This is done by creating a Workload Identity Pool and Provider in GCP.

- Create a Workload Identity Pool.
- Create a Workload Identity Provider that trusts the AWS role ARN.
- Grant the GCP service account the necessary permissions.

### 4. Containerize and Package the Application

Build a Docker image of the Go application and push it to a container registry (e.g., Amazon ECR) that your EKS cluster can access.

```Dockerfile
FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "./aws"]
```

### 5. Deploy to EKS

Create a Kubernetes deployment manifest (`pod.yaml`) to deploy your application to the EKS cluster.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: custom-credential-pod
spec:
  serviceAccountName: your-k8s-service-account # The service account associated with the IAM role
  containers:
  - name: gcp-auth-sample
    image: your-container-image:latest # Your image from ECR
    env:
    - name: AWS_REGION
      value: "your-aws-region"
    - name: GCP_WORKLOAD_AUDIENCE
      value: "your-gcp-workload-audience"
    - name: GOOGLE_CLOUD_PROJECT
      value: "your-google-cloud-project"
    # Optional: If you want to use service account impersonation
    # - name: GCP_SERVICE_ACCOUNT_IMPERSONATION_URL
    #   value: "your-gcp-service-account-impersonation-url"
    - name: GCS_BUCKET_NAME
      value: "your-gcs-bucket-name"
```

Deploy the pod:

```bash
kubectl apply -f pod.yaml
```

### 6. Clean Up

To clean up the resources, delete the EKS cluster and any other AWS and GCP resources you created.

```bash
eksctl delete cluster --name your-cluster-name
```
