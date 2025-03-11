variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "trigger_region" {
  type = string
}

variable "job_name" {
  type = string
}

variable "container_image" {
  type = string
}

variable "container_image_tag" {
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

variable "task_timeout_seconds" {
  type    = string
  default = "600s"
}
