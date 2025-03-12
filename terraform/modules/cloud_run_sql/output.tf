output "dsn" {
  value = "postgres://${var.user}:${random_password.password.result}@${google_cloud_run_v2_service.postgres.uri}/${var.database_name}}"
}

output "connection_name" {
  value = ""
}
