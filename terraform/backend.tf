terraform {
  backend "gcs" {
    bucket = "aviseu-jobs-terraform-state"
    prefix = "backoffice/state"
  }
}
