resource "google_sql_user" "user" {
  name        = var.user
  instance    = var.instance_name
  password_wo = random_password.password.result
  deletion_policy = "ABANDON"
}

resource "google_sql_database" "database" {
  depends_on = [google_sql_user.user]
  name       = var.database_name
  instance   = var.instance_name
  deletion_policy = "ABANDON"
}

resource "random_password" "password" {
  length  = 32
  special = false
}
