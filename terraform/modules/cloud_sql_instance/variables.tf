variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "instance_name" {
  type = string
}

variable "disk_size" {
  type    = number
  default = 10
}

variable "tier" {
  type    = string
  default = "db-f1-micro"
}

variable "database_version" {
  type    = string
  default = "POSTGRES_17"
}

variable "deletion_protection" {
  type    = bool
  default = true
}
