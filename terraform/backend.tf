terraform {
  backend "gcs" {
    bucket = "aviseu-jobs-terraform-state"
    prefix = "terraform/state"
  }
}
