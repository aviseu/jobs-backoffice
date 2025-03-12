output "dsn" {
  value = "postgres://${google_sql_user.user.name}:${random_password.password.result}@/${google_sql_database.database.name}?host=/cloudsql/${var.connection_name}"
}

output "connection_name" {
  value = var.connection_name
}
