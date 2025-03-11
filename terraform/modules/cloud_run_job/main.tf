resource "google_service_account" "service_account" {
  account_id   = "backoffice-${var.job_name}"
  display_name = "Backoffice ${var.job_name}"
}

resource "google_project_iam_member" "service_account_iam" {
  depends_on = [google_service_account.service_account]
  for_each   = toset(var.service_account_roles)

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_cloud_run_v2_job" "job" {
  depends_on = [google_service_account.service_account, google_project_iam_member.service_account_iam]

  project  = var.project_id
  name     = "backoffice-${var.job_name}"
  location = var.region

  deletion_protection = false

  template {
    parallelism = 1
    task_count  = 1

    template {
      max_retries     = 3
      timeout         = var.task_timeout_seconds
      service_account = google_service_account.service_account.email

      containers {
        image = "${var.container_image}:${var.container_image_tag}"

        dynamic "env" {
          for_each = var.environment_variables
          content {
            name  = env.key
            value = env.value
          }
        }

        dynamic "env" {
          for_each = var.secrets
          content {
            name = env.key
            value_source {
              secret_key_ref {
                secret  = env.value
                version = "latest"
              }
            }
          }
        }

        resources {
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }

        dynamic "volume_mounts" {
          for_each = length(var.sql_instances) > 0 ? [0] : []
          content {
            mount_path = "/cloudsql"
            name       = "cloudsql"
          }
        }
      }

      dynamic "volumes" {
        for_each = length(var.sql_instances) > 0 ? [0] : []
        content {
          name = "cloudsql"
          cloud_sql_instance {
            instances = var.sql_instances
          }
        }
      }
    }
  }
}

resource "google_service_account" "service_account_triggers" {
  account_id   = "backoffice-${var.job_name}-triggers"
  display_name = "Backoffice ${var.job_name} Triggers"
}

resource "google_project_iam_member" "service_account_triggers_iam" {
  depends_on = [google_service_account.service_account_triggers]

  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.service_account_triggers.email}"
}

resource "google_cloud_scheduler_job" "cron" {
  name             = "backoffice-${var.job_name}-cron"
  region           = var.trigger_region
  schedule         = "0 0 * * *"
  attempt_deadline = "60s"
  time_zone        = "CET"

  http_target {
    http_method = "POST"
    uri         = "https://${google_cloud_run_v2_job.job.location}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${google_cloud_run_v2_job.job.project}/jobs/${google_cloud_run_v2_job.job.name}:run"

    oauth_token {
      service_account_email = google_service_account.service_account_triggers.email
    }
  }
}
