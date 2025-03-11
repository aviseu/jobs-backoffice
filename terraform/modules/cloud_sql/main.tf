resource "google_sql_database_instance" "instance" {
  project = var.project_id
  name    = var.instance_name
  region  = var.region

  database_version    = var.database_version
  deletion_protection = var.deletion_protection

  instance_type = "CLOUD_SQL_INSTANCE"

  settings {
    tier              = var.tier
    availability_type = "ZONAL"
    disk_size         = var.disk_size
    edition           = "ENTERPRISE"

    backup_configuration {
      enabled = false
    }
  }
}

resource "google_sql_user" "user" {
  name        = var.user
  instance    = google_sql_database_instance.instance.name
  password_wo = random_password.password.result
}

resource "google_sql_database" "database" {
  depends_on = [google_sql_user.user]
  name       = var.database_name
  instance   = google_sql_database_instance.instance.name
}

resource "random_password" "password" {
  length  = 32
  special = false
}
