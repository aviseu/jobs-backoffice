variable "project_id" {
  type = string
}

variable "topic_name" {
  type = string
}

variable "subscription_name" {
  type = string
}

variable "subscription_push_endpoint" {
  type = string
}

variable "subscription_push_service_account" {
  type = string
}

variable "message_retention_duration" {
  type    = string
  default = "10s"
}
