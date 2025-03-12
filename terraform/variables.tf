variable "image_tag" {
  description = "The image tag"
  type        = string
}

variable "deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = false
}
