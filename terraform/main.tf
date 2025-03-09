module "frontend" {
  service_name            = "frontend"
  source                  = "./modules/cloud_run_service"
  project_id              = var.project_id
  region                  = var.region
  container_image         = var.frontend_container_image
  container_image_tag     = var.image_tag
  container_port          = 80
  request_max_concurrency = 5
  request_timeout_seconds = "5s"
  is_public               = true
}

module "api" {
  service_name            = "api"
  source                  = "./modules/cloud_run_service"
  project_id              = var.project_id
  region                  = var.region
  container_image         = var.api_container_image
  container_image_tag     = var.image_tag
  container_port          = 80
  request_max_concurrency = 5
  request_timeout_seconds = "5s"
  is_public               = true
  environment_variables = {
    "API_ADDR"               = "0.0.0.0:80"
    "API_CORS"               = "0"
    "DB_MAXOPENCONNS"        = "5"
    "DB_MAXIDLECONNS"        = "5"
    "PUBSUB_PROJECT_ID"      = "aviseu-jobs"
    "PUBSUB_IMPORT_TOPIC_ID" = "imports"
  }
  sql_instances = [
    "aviseu-jobs:europe-west4:jobs"
  ]
  secrets = {
    "DB_DSN" : "jobs-dsn"
  }
  service_account_roles = [
    "roles/cloudsql.client",
    "roles/pubsub.publisher",
    "roles/secretmanager.secretAccessor"
  ]
}

module "import" {
  service_name            = "import"
  source                  = "./modules/cloud_run_service"
  project_id              = var.project_id
  region                  = var.region
  container_image         = var.import_container_image
  container_image_tag     = var.image_tag
  container_port          = 80
  request_max_concurrency = 1
  request_timeout_seconds = "60s"
  is_public               = false
  environment_variables = {
    "IMPORT_ADDR"                       = "0.0.0.0:80"
    "DB_MAXOPENCONNS"                   = "5"
    "DB_MAXIDLECONNS"                   = "5"
    "PUBSUB_PROJECT_ID"                 = "aviseu-jobs"
    "PUBSUB_IMPORT_TOPIC_ID"            = "imports"
    "IMPORT_MAX_CONNECTIONS"            = "1"
    "JOB_WORKERS"                       = "2"
    "JOB_BUFFER"                        = "10"
    "GATEWAY_IMPORT_RESULT_BUFFER_SIZE" = "10"
    "GATEWAY_IMPORT_RESULT_WORKERS"     = "2"

  }
  sql_instances = [
    "aviseu-jobs:europe-west4:jobs"
  ]
  secrets = {
    "DB_DSN" : "jobs-dsn"
  }
  service_account_roles = [
    "roles/cloudsql.client",
    "roles/secretmanager.secretAccessor"
  ]
  pubsub_triggers = {
    "imports" : "/import"
  }
}

module "schedule" {
  job_name             = "schedule"
  source               = "./modules/cloud_run_job"
  project_id           = var.project_id
  region               = var.region
  trigger_region       = var.trigger_region
  container_image      = var.schedule_container_image
  container_image_tag  = var.image_tag
  task_timeout_seconds = "600s"
  environment_variables = {
    "PUBSUB_PROJECT_ID"      = "aviseu-jobs"
    "PUBSUB_IMPORT_TOPIC_ID" = "imports"

  }
  sql_instances = [
    "aviseu-jobs:europe-west4:jobs"
  ]
  secrets = {
    "DB_DSN" : "jobs-dsn"
  }
  service_account_roles = [
    "roles/cloudsql.client",
    "roles/pubsub.publisher",
    "roles/secretmanager.secretAccessor"
  ]
}
