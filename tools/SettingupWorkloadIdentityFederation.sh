#!/bin/bash


gcloud iam service-accounts create "my-service-account" \
  --project "${PROJECT_ID}"



gcloud iam workload-identity-pools providers create-oidc "github-action-provider" \
  --project="rabbitgather" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --display-name="github-action-provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"

# github-action-service-account


# TODO(developer): Update this value to your GitHub repository.
export REPO="username/name" # e.g. "google/chrome"
export WORKLOAD_IDENTITY_POOL_ID="projects/517575915083/locations/global/workloadIdentityPools/github-pool"
gcloud iam service-accounts add-iam-policy-binding "github-action-service-account@${PROJECT_ID}.iam.gserviceaccount.com" \
--project="${PROJECT_ID}" \
--role="roles/iam.workloadIdentityUser" \
--member="principalSet://iam.googleapis.com/${WORKLOAD_IDENTITY_POOL_ID}/attribute.repository/${REPO}"

export PROJECT_ID="rabbitgather" # update with your value

gcloud iam workload-identity-pools providers describe "github-action-provider" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --format="value(name)"

export workload_identity_provider="projects/517575915083/locations/global/workloadIdentityPools/github-pool/providers/github-action-provider"
