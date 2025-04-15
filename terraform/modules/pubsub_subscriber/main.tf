resource "google_pubsub_subscription" "subscription" {
  project = var.project_id
  name    = var.subscription_name
  topic   = var.topic_name

  ack_deadline_seconds = 60
  message_retention_duration = var.message_retention_duration

  push_config {
    push_endpoint = var.subscription_push_endpoint

    attributes = {
      x-goog-version = "v1"
    }
  }
}
