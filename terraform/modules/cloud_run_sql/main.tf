resource "random_password" "password" {
  length  = 32
  special = false
}

resource "google_cloud_run_v2_service" "postgres" {
  name     = var.instance_name
  location = var.region

  template {
    containers {
      image = "postgres:17"

      ports {
        container_port = 5432
      }

      env {
        name  = "POSTGRES_USER"
        value = var.user
      }

      env {
        name  = "POSTGRES_PASSWORD"
        value = random_password.password.result
      }

      env {
        name  = "POSTGRES_DB"
        value = var.database_name
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "128Mi"
        }
        cpu_idle          = true
        startup_cpu_boost = true
      }
    }

    max_instance_request_concurrency = 50
    timeout                          = "30s"
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}
