variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "load_balancer_name" {
  type = string
}

variable "address_name" {
  type = string
}

variable "backends" {
  type = map(string)
}

variable "default_backend" {
  type = string
}

variable "routes" {
  type = map(object({
    domain      = string
    certificate = string
    paths       = map(string)
  }))
}
