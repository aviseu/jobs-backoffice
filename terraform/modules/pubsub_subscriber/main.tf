resource "google_pubsub_subscription" "subscription" {
  project = var.project_id
  name    = var.subscription_name
  topic   = var.topic_name

  ack_deadline_seconds = 60
  message_retention_duration = var.message_retention_duration

  push_config {
    push_endpoint = var.subscription_push_endpoint

    oidc_token {
      service_account_email = var.subscription_push_service_account
    }

    attributes = {
      x-goog-version = "v1"
    }
  }

  expiration_policy {
    ttl = ""
  }
}
