#!/bin/bash
set -e

# Configuration
PROJECT_ID=${PROJECT_ID:-$(gcloud config get-value project)}
REGION=${REGION:-us-central1}
REPO_NAME="agent-cra-repo"
IMAGE_TAG="latest"
SERVICE_ACCOUNT_NAME="agent-cra-identity"

echo "==================================================="
echo "Deploying Multi-Agent System to Cloud Run"
echo "Project: $PROJECT_ID"
echo "Region:  $REGION"
echo "==================================================="

# 1. Enable APIs
echo "Enabling required APIs..."
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  secretmanager.googleapis.com \
  cloudbuild.googleapis.com

# 2. Create Service Account (Identity)
if ! gcloud iam service-accounts describe "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" &>/dev/null; then
    echo "Creating Service Account: $SERVICE_ACCOUNT_NAME..."
    gcloud iam service-accounts create "$SERVICE_ACCOUNT_NAME" \
        --display-name="Agent CRA Cloud Run Identity"
else
    echo "Service Account $SERVICE_ACCOUNT_NAME exists."
fi

# Grant Vertex AI User role
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/aiplatform.user"

# Grant Secret Accessor role
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor"

# 3. Create Artifact Registry
if ! gcloud artifacts repositories describe "$REPO_NAME" --location="$REGION" &>/dev/null; then
    echo "Creating Artifact Registry repository..."
    gcloud artifacts repositories create "$REPO_NAME" \
        --repository-format=docker \
        --location="$REGION" \
        --description="Docker repository for Agent CRA"
else
    echo "Artifact Registry $REPO_NAME exists."
fi

# 4. Build and Push Image
IMAGE_URI="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO_NAME}/agent-cra:${IMAGE_TAG}"
echo "Building and pushing container image to $IMAGE_URI..."
gcloud builds submit --tag "$IMAGE_URI" .

# 5. Deploy Agents as Micro-Services
# We deploy 4 instances of the same image, but with different arguments/env vars if supported by the app.
# Note: The current Go app in main.go runs ALL agents in one process via the Coordinator.
# To support true Cloud Run micro-services, the Go app needs to accept a flag to run *only* one agent.
# Assuming the previous refactor added args support (e.g., cmd/main.go uses os.Args or flags).

AGENTS=("classifier" "auditor" "vuln" "reporter")

for AGENT in "${AGENTS[@]}"; do
    SERVICE_NAME="agent-${AGENT}"
    echo "Deploying Cloud Run Service: $SERVICE_NAME..."

    gcloud run deploy "$SERVICE_NAME" \
        --image "$IMAGE_URI" \
        --region "$REGION" \
        --service-account "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
        --set-env-vars="PROJECT_ID=${PROJECT_ID},AGENT_ROLE=${AGENT}" \
        --set-secrets="GEMINI_API_KEY=gemini-api-key:latest" \
        --no-allow-unauthenticated \
        --args="--role=${AGENT}","--mode=server" # Passing role and server mode
done

echo "==================================================="
echo "Deployment Complete!"
echo "Services:"
gcloud run services list --format="table(SERVICE,REGION,URL,LAST_DEPLOYED_BY,LAST_DEPLOYED_TIME)"
echo "==================================================="
