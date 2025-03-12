output "dsn" {
  value = "postgres://${var.user}:${random_password.password.result}@${trimprefix(google_cloud_run_v2_service.postgres.uri, "https://")}/${var.database_name}}"
}

output "connection_name" {
  value = ""
}
