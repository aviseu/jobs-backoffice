variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "service_name" {
  type = string
}

variable "container_port" {
  type = string
}

variable "container_image" {
  type = string
}

variable "container_image_tag" {
  type = string
}

variable "request_max_concurrency" {
  type = number
}

variable "request_timeout_seconds" {
  type = string
}

variable "environment_variables" {
  type    = map(string)
  default = {}
}

variable "secrets" {
  type    = map(string)
  default = {}
}

variable "sql_instances" {
  type    = list(string)
  default = []
}

variable "service_account_roles" {
  type    = list(string)
  default = []
}

variable "pubsub_triggers" {
  type = map(object({
    topic = string
    path  = string
  }))
  default = {}
}

variable "is_public" {
  type = bool
}
