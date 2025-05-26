module "database" {
  source          = "github.com/aviseu/terraform//modules/gcp_cloud_sql_database"
  instance_name   = "jobs-db"
  connection_name = "aviseu-jobs:europe-west4:jobs-db"
  database_name   = "backoffice"
  user            = "backoffice"
}

module "dsn" {
  depends_on = [module.database]

  source      = "github.com/aviseu/terraform//modules/gcp_secret"
  project_id  = "aviseu-jobs"
  secret_name = "backoffice-dsn"
  secret_data = module.database.dsn
}

module "importsTopic" {
  source     = "github.com/aviseu/terraform//modules/gcp_pubsub_topic"
  project_id = "aviseu-jobs"
  topic_name = "imports"
}

module "jobsTopic" {
  source     = "github.com/aviseu/terraform//modules/gcp_pubsub_topic"
  project_id = "aviseu-jobs"
  topic_name = "jobs"
}

module "frontend" {
  service_name            = "frontend"
  source                  = "github.com/aviseu/terraform//modules/gcp_cloud_run_service"
  project_id              = "aviseu-jobs"
  region                  = "europe-west4"
  container_image         = "europe-west4-docker.pkg.dev/aviseu-jobs/jobs/jobs-backoffice-frontend"
  container_image_tag     = var.image_tag
  container_port          = 80
  request_max_concurrency = 5
  request_timeout_seconds = "5s"
  is_public               = true
}

module "api" {
  depends_on = [module.database, module.dsn]

  service_name            = "api"
  source                  = "github.com/aviseu/terraform//modules/gcp_cloud_run_service"
  project_id              = "aviseu-jobs"
  region                  = "europe-west4"
  container_image         = "europe-west4-docker.pkg.dev/aviseu-jobs/jobs/jobs-backoffice-api"
  container_image_tag     = var.image_tag
  container_port          = 80
  request_max_concurrency = 5
  request_timeout_seconds = "5s"
  is_public               = true

  environment_variables = {
    "API_ADDR"               = "0.0.0.0:80"
    "API_CORS"               = "0"
    "DB_MAXOPENCONNS"        = "20"
    "DB_MAXIDLECONNS"        = "20"
    "PUBSUB_PROJECT_ID"      = "aviseu-jobs"
    "PUBSUB_IMPORT_TOPIC_ID" = module.importsTopic.topic_name
  }

  sql_instances = length(module.database.connection_name) > 0 ? [
    module.database.connection_name
  ] : []

  secrets = {
    "DB_DSN" : module.dsn.secret_id
  }

  service_account_roles = [
    "roles/cloudsql.client",
    "roles/pubsub.publisher",
    "roles/secretmanager.secretAccessor"
  ]
}

module "import" {
  depends_on = [module.database, module.dsn, module.importsTopic]

  service_name            = "import"
  source                  = "github.com/aviseu/terraform//modules/gcp_cloud_run_service"
  project_id              = "aviseu-jobs"
  region                  = "europe-west4"
  container_image         = "europe-west4-docker.pkg.dev/aviseu-jobs/jobs/jobs-backoffice-import"
  container_image_tag     = var.image_tag
  container_port          = 80
  request_max_concurrency = 1
  request_timeout_seconds = "60s"
  is_public               = false

  environment_variables = {
    "IMPORT_ADDR"                        = "0.0.0.0:80"
    "DB_MAXOPENCONNS"                    = "20"
    "DB_MAXIDLECONNS"                    = "20"
    "PUBSUB_PROJECT_ID"                  = "aviseu-jobs"
    "PUBSUB_JOB_TOPIC_ID"                = "jobs"
    "IMPORT_MAX_CONNECTIONS"             = "1"
    "GATEWAY_IMPORT_METRIC_BUFFER_SIZE"  = "10"
    "GATEWAY_IMPORT_METRIC_WORKERS"      = "2"
    "GATEWAY_IMPORT_PUBLISH_BUFFER_SIZE" = "10"
    "GATEWAY_IMPORT_PUBLISH_WORKERS"     = "2"
    "GATEWAY_IMPORT_JOB_BUFFER_SIZE"     = "10"
    "GATEWAY_IMPORT_JOB_WORKERS"         = "2"
  }

  sql_instances = length(module.database.connection_name) > 0 ? [
    module.database.connection_name
  ] : []

  secrets = {
    "DB_DSN" : module.dsn.secret_id
  }

  service_account_roles = [
    "roles/cloudsql.client",
    "roles/secretmanager.secretAccessor",
    "roles/run.invoker",
    "roles/pubsub.publisher"
  ]
}

module "importsSubscription" {
  source     = "github.com/aviseu/terraform//modules/gcp_pubsub_subscriber"
  project_id = "aviseu-jobs"
  topic_name = module.importsTopic.topic_name
  subscription_name = "imports-subscription"
  subscription_push_endpoint = "${module.import.service_url}/import"
  message_retention_duration = "600s"
  subscription_push_service_account = module.import.service_account_email
}

module "schedule" {
  depends_on = [module.database, module.dsn, module.importsTopic]

  job_name             = "schedule"
  source               = "github.com/aviseu/terraform//modules/gcp_cloud_run_job"
  project_id           = "aviseu-jobs"
  region               = "europe-west4"
  trigger_region       = "europe-west3"
  container_image      = "europe-west4-docker.pkg.dev/aviseu-jobs/jobs/jobs-backoffice-schedule"
  container_image_tag  = var.image_tag
  task_timeout_seconds = "600s"

  environment_variables = {
    "PUBSUB_PROJECT_ID"      = "aviseu-jobs"
    "PUBSUB_IMPORT_TOPIC_ID" = module.importsTopic.topic_name
  }

  sql_instances = length(module.database.connection_name) > 0 ? [
    module.database.connection_name
  ] : []

  secrets = {
    "DB_DSN" : module.dsn.secret_id
  }

  service_account_roles = [
    "roles/cloudsql.client",
    "roles/pubsub.publisher",
    "roles/secretmanager.secretAccessor"
  ]
}
