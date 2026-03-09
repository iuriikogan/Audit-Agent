# cloud_sql.tf - Defines the Cloud SQL (PostgreSQL) instance

resource "google_sql_database_instance" "main" {
  name             = "cra-db-instance"
  database_version = var.db_version
  region           = var.region
  settings {
    tier = var.db_tier
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.vpc.id
    }
  }
  deletion_protection = false # Set to true for production environments

  depends_on = [google_project_service.apis, google_service_networking_connection.private_vpc_connection]
}

resource "google_sql_database" "main" {
  name     = "cra_db"
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "main" {
  name     = var.db_user
  instance = google_sql_database_instance.main.name
  password = var.db_password
}
