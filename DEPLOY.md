# Deployment Instructions

This document provides detailed instructions for deploying the Multi-Agent CRA System locally, to Google Kubernetes Engine (GKE) using Terraform, and to Cloud Run using the provided shell script.

## 1. Local Deployment

Run the application locally for development and testing.

### Prerequisites
*   **Go**: Version 1.25 or higher ([Download](https://go.dev/dl/)).
*   **Google Gemini API Key**: A valid API Key.

### Steps
1.  **Clone the repository** (if not already done):
    ```bash
    git clone <repository-url>
    cd multi-agent-cra
    ```

2.  **Set Environment Variables**:
    ```bash
    # Linux/macOS
    export GEMINI_API_KEY="your_actual_api_key_here"

    # Windows (PowerShell)
    $env:GEMINI_API_KEY="your_actual_api_key_here"
    ```

3.  **Run the Application**:
    ```bash
    go run cmd/main.go
    ```
    *Note: The current local execution runs all agents within a single process via the coordinator.*

---

## 2. Google Kubernetes Engine (GKE) Deployment

This method uses **Terraform** to provision a GKE Autopilot cluster, secure secrets, and deploy the agents as separate Kubernetes workloads.

### Prerequisites
*   **Google Cloud Project**: With billing enabled.
*   **Terraform**: Installed.
*   **gcloud CLI**: Installed and authenticated (`gcloud auth login`, `gcloud auth application-default login`).
*   **Docker**: For building the image.

### Step 1: Build and Push Docker Image
Before running Terraform, the container image must exist in a registry (e.g., Google Artifact Registry or Container Registry).

1.  **Set Variables**:
    ```bash
    export PROJECT_ID="your-project-id"
    export IMAGE_NAME="gcr.io/${PROJECT_ID}/agent-cra:latest"
    ```

2.  **Build and Push**:
    ```bash
    # Enable Container Registry API if needed, or use Artifact Registry
    gcloud services enable containerregistry.googleapis.com

    # Build
    docker build -t $IMAGE_NAME .

    # Configure Docker to push to GCR
    gcloud auth configure-docker

    # Push
    docker push $IMAGE_NAME
    ```

### Step 2: Deploy Infrastructure with Terraform
1.  **Navigate to the Terraform directory**:
    ```bash
    cd terraform
    ```

2.  **Create a `terraform.tfvars` file**:
    Create a file named `terraform.tfvars` with your specific configuration. **Do not commit this file.**
    ```hcl
    project_id       = "your-project-id"
    region           = "us-central1"
    cluster_name     = "agent-engine-cluster"
    image_repository = "gcr.io/your-project-id/agent-cra:latest" # Must match the image pushed in Step 1
    gemini_api_key   = "your-actual-gemini-api-key"
    ```

3.  **Initialize and Apply**:
    ```bash
    terraform init
    terraform apply
    ```
    *Confirm the action by typing `yes` when prompted.*

    **What this does:**
    *   Creates a VPC Network and Subnet.
    *   Provisions a GKE Autopilot Cluster.
    *   Creates a Secret in Google Secret Manager for the API Key.
    *   Sets up Workload Identity (IAM binding between K8s Service Accounts and Google Service Accounts).
    *   Deploys 4 microservices (`agent-classifier`, `agent-auditor`, `agent-vuln`, `agent-reporter`).

### Step 3: Verify Deployment
1.  **Get Cluster Credentials**:
    ```bash
    gcloud container clusters get-credentials agent-engine-cluster --region us-central1
    ```

2.  **Check Pods**:
    ```bash
    kubectl get pods
    ```
    You should see pods for each agent (classifier, auditor, vuln, reporter) running.

---

## 3. Cloud Run Deployment

This method uses the `deploy.sh` script to deploy the agents as serverless Cloud Run services.

### Prerequisites
*   **Google Cloud SDK**: `gcloud` installed and authenticated.
*   **Project ID**: Set your active project (`gcloud config set project YOUR_PROJECT_ID`).

### Step 1: Create the API Key Secret
The deployment script expects a secret named `gemini-api-key` to exist in Secret Manager.

```bash
# Replace YOUR_API_KEY with your actual key
echo -n "YOUR_API_KEY" | gcloud secrets create gemini-api-key --data-file=-
```

### Step 2: Run the Deployment Script
1.  **Make the script executable**:
    ```bash
    chmod +x deploy.sh
    ```

2.  **Run the script**:
    ```bash
    ./deploy.sh
    ```

    **What this script does:**
    *   Enables necessary Google Cloud APIs.
    *   Creates a dedicated Service Account.
    *   Creates an Artifact Registry repository.
    *   Builds the Docker image using Cloud Build (no local Docker required).
    *   Deploys 4 Cloud Run services, injecting the API Key secret and setting the `AGENT_ROLE` environment variable.

### Step 3: Verify
The script will output the URLs of the deployed services. You can also list them:
```bash
gcloud run services list
```
