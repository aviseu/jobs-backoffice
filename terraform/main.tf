module "database" {
  source          = "github.com/aviseu/terraform//modules/cloud_sql_database?ref=v1.2.3"
  instance_name   = "jobs-db"
  connection_name = "aviseu-jobs:europe-west4:jobs-db"
  database_name   = "backoffice"
  user            = "backoffice"
}

module "dsn" {
  depends_on = [module.database]

  source      = "github.com/aviseu/terraform//modules/secret?ref=v1.2.3"
  project_id  = "aviseu-jobs"
  secret_name = "backoffice-dsn"
  secret_data = module.database.dsn
}

module "importsTopic" {
  source     = "github.com/aviseu/terraform//modules/pubsub_topic?ref=v1.2.3"
  project_id = "aviseu-jobs"
  topic_name = "imports"
}

module "frontend" {
  service_name            = "frontend"
  source                  = "github.com/aviseu/terraform//modules/cloud_run_service?ref=v1.2.3"
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
  source                  = "github.com/aviseu/terraform//modules/cloud_run_service?ref=v1.2.3"
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
    "DB_MAXOPENCONNS"        = "5"
    "DB_MAXIDLECONNS"        = "5"
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
  source                  = "github.com/aviseu/terraform//modules/cloud_run_service?ref=v1.2.3"
  project_id              = "aviseu-jobs"
  region                  = "europe-west4"
  container_image         = "europe-west4-docker.pkg.dev/aviseu-jobs/jobs/jobs-backoffice-import"
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

  sql_instances = length(module.database.connection_name) > 0 ? [
    module.database.connection_name
  ] : []

  secrets = {
    "DB_DSN" : module.dsn.secret_id
  }

  service_account_roles = [
    "roles/cloudsql.client",
    "roles/secretmanager.secretAccessor",
    "roles/run.invoker"
  ]
}

module "importsSubscription" {
  source     = "github.com/aviseu/terraform//modules/pubsub_subscriber?ref=v1.2.3"
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
  source               = "github.com/aviseu/terraform//modules/cloud_run_job?ref=v1.2.3"
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
