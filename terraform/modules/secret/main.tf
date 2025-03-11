resource "google_secret_manager_secret" "dsn" {
  project   = var.project_id
  secret_id = var.secret_name

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "version" {
  depends_on = [google_secret_manager_secret.dsn]

  enabled        = true
  secret         = google_secret_manager_secret.dsn.id
  secret_data_wo = var.secret_data
}
