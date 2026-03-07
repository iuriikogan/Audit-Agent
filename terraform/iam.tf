# ------------------------------------------------------------------------------
# 1. Scope Classifier (Reader/Predictor)
# ------------------------------------------------------------------------------
resource "google_service_account" "sa_classifier" {
  account_id   = "sa-classifier"
  display_name = "Agent: Scope Classifier"
}

# Can use Vertex AI to generate content
resource "google_project_iam_member" "classifier_vertex" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.sa_classifier.email}"
}

# ------------------------------------------------------------------------------
# 2. Regulatory Auditor (Checker/Reader)
# ------------------------------------------------------------------------------
resource "google_service_account" "sa_auditor" {
  account_id   = "sa-auditor"
  display_name = "Agent: Regulatory Auditor"
}

# Can use Vertex AI to generate content (verify)
resource "google_project_iam_member" "auditor_vertex" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.sa_auditor.email}"
}

# ------------------------------------------------------------------------------
# 3. Vulnerability Watchdog (Reader/External Query)
# ------------------------------------------------------------------------------
resource "google_service_account" "sa_vuln" {
  account_id   = "sa-vuln"
  display_name = "Agent: Vulnerability Watchdog"
}

# Can use Vertex AI
resource "google_project_iam_member" "vuln_vertex" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.sa_vuln.email}"
}

# ------------------------------------------------------------------------------
# 4. Compliance Reporter (Writer)
# ------------------------------------------------------------------------------
resource "google_service_account" "sa_reporter" {
  account_id   = "sa-reporter"
  display_name = "Agent: Compliance Reporter"
}

# Can use Vertex AI
resource "google_project_iam_member" "reporter_vertex" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.sa_reporter.email}"
}

# Grant Writer permission ONLY to the Reporter (e.g., to GCS bucket for reports)
# (Assuming a bucket 'cra-reports-${project_id}' exists or will be created)
# resource "google_storage_bucket_iam_member" "reporter_gcs_write" {
#   bucket = "cra-reports-${var.project_id}"
#   role   = "roles/storage.objectCreator"
#   member = "serviceAccount:${google_service_account.sa_reporter.email}"
# }

# ------------------------------------------------------------------------------
# 5. Workload Identity Bindings
# ------------------------------------------------------------------------------

# Classifier
resource "google_service_account_iam_member" "wi_classifier" {
  service_account_id = google_service_account.sa_classifier.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[default/ksa-classifier]"
}

# Auditor
resource "google_service_account_iam_member" "wi_auditor" {
  service_account_id = google_service_account.sa_auditor.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[default/ksa-auditor]"
}

# Vuln Watchdog
resource "google_service_account_iam_member" "wi_vuln" {
  service_account_id = google_service_account.sa_vuln.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[default/ksa-vuln]"
}

# Reporter
resource "google_service_account_iam_member" "wi_reporter" {
  service_account_id = google_service_account.sa_reporter.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[default/ksa-reporter]"
}