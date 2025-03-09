variable "project_id" {
  description = "The Google Cloud project ID"
  type        = string
}

variable "region" {
  description = "The region to deploy Cloud Run services"
  type        = string
  default     = "europe-west4"
}

variable "trigger_region" {
  description = "The region to trigger Cloud Run jobs"
  type        = string
  default     = "europe-west3"
}

variable "frontend_container_image" {
  description = "The frontend server image"
  type        = string
}

variable "api_container_image" {
  description = "The api server image"
  type        = string
}

variable "schedule_container_image" {
  description = "The schedule job image"
  type        = string
}

variable "import_container_image" {
  description = "The import server image"
  type        = string
}

variable "image_tag" {
  description = "The image tag"
  type        = string
}

variable "deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = false
}
