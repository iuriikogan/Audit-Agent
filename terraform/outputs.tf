output "project_id" {
  description = "Google Cloud Project ID"
  value       = var.project_id
}

output "region" {
  description = "Google Cloud Region"
  value       = var.region
}

output "artifact_registry_repo" {
  description = "Artifact Registry Docker Repository ID"
  value       = google_artifact_registry_repository.cra_repo.repository_id
}

output "pubsub_scan_topic" {
  description = "Pub/Sub Topic for Scan Requests"
  value       = google_pubsub_topic.scan_requests.name
}


output "dashboard_public_url" {
  description = "The public IP address to access the dashboard"
  value       = "http://${google_compute_global_address.default.address}"
}
