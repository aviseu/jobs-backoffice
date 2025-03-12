module "database_instance" {
  source              = "github.com/aviseu/terraform//modules/cloud_sql_instance?ref=v1.0.0"
  project_id          = var.project_id
  region              = var.region
  instance_name       = "jobs"
  database_version    = "POSTGRES_17"
  tier                = "db-f1-micro"
  disk_size           = 10
  deletion_protection = var.deletion_protection
}

module "database" {
  source          = "github.com/aviseu/terraform//modules/cloud_sql_database?ref=v1.0.0"
  instance_name   = module.database_instance.name
  connection_name = module.database_instance.connection_name
  database_name   = "jobs"
  user            = "jobs"
}

module "dsn" {
  depends_on = [module.database]

  source      = "github.com/aviseu/terraform//modules/secret?ref=v1.0.0"
  project_id  = var.project_id
  secret_name = "jobs-dsn"
  secret_data = module.database.dsn
}

module "importsTopic" {
  source     = "github.com/aviseu/terraform//modules/pubsub?ref=v1.0.0"
  project_id = var.project_id
  topic_name = "imports"
}

module "frontend" {
  service_name            = "frontend"
  source                  = "github.com/aviseu/terraform//modules/cloud_run_service?ref=v1.0.0"
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
  depends_on = [module.database, module.dsn]

  service_name            = "api"
  source                  = "github.com/aviseu/terraform//modules/cloud_run_service?ref=v1.0.0"
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
  source                  = "github.com/aviseu/terraform//modules/cloud_run_service?ref=v1.0.0"
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

  sql_instances = length(module.database.connection_name) > 0 ? [
    module.database.connection_name
  ] : []

  secrets = {
    "DB_DSN" : module.dsn.secret_id
  }

  service_account_roles = [
    "roles/cloudsql.client",
    "roles/secretmanager.secretAccessor"
  ]

  pubsub_triggers = {
    "imports" : {
      topic : module.importsTopic.topic_name,
      path : "/import"
    }
  }
}

module "schedule" {
  depends_on = [module.database, module.dsn, module.importsTopic]

  job_name             = "schedule"
  source               = "github.com/aviseu/terraform//modules/cloud_run_job?ref=v1.0.0"
  project_id           = var.project_id
  region               = var.region
  trigger_region       = var.trigger_region
  container_image      = var.schedule_container_image
  container_image_tag  = var.image_tag
  task_timeout_seconds = "600s"

  environment_variables = {
    "PUBSUB_PROJECT_ID"      = var.project_id
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

module "load_balancer" {
  source             = "github.com/aviseu/terraform//modules/load_balancer?ref=v1.0.0"
  project_id         = var.project_id
  region             = var.region
  load_balancer_name = "jobs-lb"
  address_name       = "public-ip"
  default_backend    = "frontend"

  backends = {
    "frontend" : module.frontend.backend,
    "api" : module.api.backend
  }

  routes = {
    jobs-backoffice-viseu-me : {
      domain : "jobs-backoffice.viseu.me"
      certificate : "viseu-me-cloudflare-origin"
      paths : {
        "/" : "frontend",
        "/api/*" : "api"
      }
    }
  }
}
