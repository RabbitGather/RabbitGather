#!/bin/bash

export PROJECT_ID="rabbitgather" # update with your value
export SERVICE_ACCOUNT_NAME="github-service-account" # update with your value

gcloud iam service-accounts create "${SERVICE_ACCOUNT_NAME}" \
  --project "${PROJECT_ID}"

gcloud services enable iamcredentials.googleapis.com \
  --project "${PROJECT_ID}"

export WORKLOAD_IDENTITY_POOL_NAME="github-workload-identity-pool"

gcloud iam workload-identity-pools create "${WORKLOAD_IDENTITY_POOL_NAME}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --display-name="${WORKLOAD_IDENTITY_POOL_NAME}"

gcloud iam workload-identity-pools describe "${WORKLOAD_IDENTITY_POOL_NAME}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --format="value(name)"

# projects/517575915083/locations/global/workloadIdentityPools/github-workload-identity-pool
export WORKLOAD_IDENTITY_POOL_ID=$(gcloud iam workload-identity-pools describe "${WORKLOAD_IDENTITY_POOL_NAME}" \
                                     --project="${PROJECT_ID}" \
                                     --location="global" \
                                     --format="value(name)")


export WORKLOAD_IDENTITY_POOL_PROVIDER="github-provider"

gcloud iam workload-identity-pools providers create-oidc "${WORKLOAD_IDENTITY_POOL_PROVIDER}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="${WORKLOAD_IDENTITY_POOL_PROVIDER}" \
  --display-name="${WORKLOAD_IDENTITY_POOL_PROVIDER_DISPLAY_NAME}" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"

# TODO(developer): Update this value to your GitHub repository.
export REPO="meowalien/RabbitGather"

gcloud iam service-accounts add-iam-policy-binding "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WORKLOAD_IDENTITY_POOL_ID}/attribute.repository/${REPO}"

# projects/517575915083/locations/global/workloadIdentityPools/github-workload-identity-pool/providers/github-provider
export WORKLOAD_IDENTITY_PROVIDER_FULL_NAME=$(gcloud iam workload-identity-pools providers describe "${WORKLOAD_IDENTITY_POOL_PROVIDER}" \
                                                  --project="${PROJECT_ID}" \
                                                  --location="global" \
                                                  --workload-identity-pool="${WORKLOAD_IDENTITY_POOL_NAME}" \
                                                  --format="value(name)")


                                     # value from above

# gcloud iam service-accounts create "my-service-account" \
#   --project "${PROJECT_ID}"
#
#
#
# gcloud iam workload-identity-pools providers create-oidc "github-action-provider" \
#   --project="rabbitgather" \
#   --location="global" \
#   --workload-identity-pool="github-pool" \
#   --display-name="github-action-provider" \
#   --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
#   --issuer-uri="https://token.actions.githubusercontent.com"
#
# # github-action-service-account
#
#
# # TODO(developer): Update this value to your GitHub repository.
# export REPO="meowalien/RabbitGather" # e.g. "google/chrome"
# export WORKLOAD_IDENTITY_POOL_ID="projects/517575915083/locations/global/workloadIdentityPools/github-pool"
# gcloud iam service-accounts add-iam-policy-binding "github-action-service-account@${PROJECT_ID}.iam.gserviceaccount.com" \
# --project="${PROJECT_ID}" \
# --role="roles/iam.workloadIdentityUser" \
# --member="principalSet://iam.googleapis.com/${WORKLOAD_IDENTITY_POOL_ID}/attribute.repository/${REPO}"
#
# export PROJECT_ID="rabbitgather" # update with your value
#
# gcloud iam workload-identity-pools providers describe "github-action-provider" \
#   --project="${PROJECT_ID}" \
#   --location="global" \
#   --workload-identity-pool="github-pool" \
#   --format="value(name)"
#
# export workload_identity_provider="projects/517575915083/locations/global/workloadIdentityPools/github-pool/providers/github-action-provider"
