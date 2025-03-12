resource "google_service_account" "service_account" {
  account_id   = "backoffice-${var.service_name}"
  display_name = "Backoffice ${var.service_name}"
}

resource "google_project_iam_member" "service_account_iam" {
  depends_on = [google_service_account.service_account]
  for_each   = toset(var.service_account_roles)

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_cloud_run_v2_service" "service" {
  depends_on = [google_service_account.service_account, google_project_iam_member.service_account_iam]

  project  = var.project_id
  name     = "backoffice-${var.service_name}"
  location = var.region

  deletion_protection = false

  template {
    containers {
      image = "${var.container_image}:${var.container_image_tag}"
      ports {
        container_port = var.container_port
      }

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
          memory = "128Mi"
        }
        cpu_idle          = true
        startup_cpu_boost = true
      }

      startup_probe {
        tcp_socket {
          port = var.container_port
        }
        initial_delay_seconds = 0
        period_seconds        = 60
        timeout_seconds       = 1
        failure_threshold     = 3
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

    max_instance_request_concurrency = var.request_max_concurrency
    timeout                          = var.request_timeout_seconds

    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    service_account = google_service_account.service_account.email
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  count       = var.is_public ? 1 : 0
  location    = google_cloud_run_v2_service.service.location
  project     = google_cloud_run_v2_service.service.project
  service     = google_cloud_run_v2_service.service.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

resource "google_service_account" "service_account_triggers" {
  count = length(var.pubsub_triggers) > 0 ? 1 : 0

  account_id   = "backoffice-${var.service_name}-triggers"
  display_name = "Backoffice ${var.service_name} Triggers"
}

resource "google_project_iam_member" "service_account_triggers_iam" {
  depends_on = [google_service_account.service_account_triggers]
  count      = length(var.pubsub_triggers) > 0 ? 1 : 0

  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.service_account_triggers[0].email}"
}

resource "google_eventarc_trigger" "primary" {
  for_each = var.pubsub_triggers

  name            = "${var.service_name}-${each.key}-trigger"
  location        = var.region
  service_account = google_service_account.service_account_triggers[0].email

  matching_criteria {
    attribute = "type"
    value     = "google.cloud.pubsub.topic.v1.messagePublished"
  }

  destination {
    cloud_run_service {
      service = google_cloud_run_v2_service.service.name
      region  = var.region
      path    = each.value.path
    }
  }

  transport {
    pubsub {
      topic = each.value.topic
    }
  }
}

resource "google_compute_region_network_endpoint_group" "group" {
  name                  = "${google_cloud_run_v2_service.service.name}-neg"
  project               = var.project_id
  region                = var.region
  network_endpoint_type = "SERVERLESS"

  cloud_run {
    service = google_cloud_run_v2_service.service.name
  }
}

resource "google_compute_region_backend_service" "backend" {
  region                = var.region
  project               = var.project_id
  name                  = "${google_cloud_run_v2_service.service.name}-backend"
  protocol              = "HTTPS"
  port_name             = "http"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  timeout_sec           = 30

  backend {
    group           = google_compute_region_network_endpoint_group.group.self_link
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }

  connection_draining_timeout_sec = 0
}
