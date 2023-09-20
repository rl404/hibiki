resource "google_cloud_run_v2_service" "server" {
  name     = var.cloud_run_name
  location = var.cloud_run_location
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    labels = {
      app = var.cloud_run_name
    }
    scaling {
      min_instance_count = 0
    }
    containers {
      name    = var.cloud_run_name
      image   = var.gcr_image_name
      command = ["./hibiki"]
      args    = ["server"]
      env {
        name  = "HIBIKI_CACHE_DIALECT"
        value = var.hibiki_cache_dialect
      }
      env {
        name  = "HIBIKI_CACHE_ADDRESS"
        value = var.hibiki_cache_address
      }
      env {
        name  = "HIBIKI_CACHE_PASSWORD"
        value = var.hibiki_cache_password
      }
      env {
        name  = "HIBIKI_CACHE_TIME"
        value = var.hibiki_cache_time
      }
      env {
        name  = "HIBIKI_DB_DIALECT"
        value = var.hibiki_db_dialect
      }
      env {
        name  = "HIBIKI_DB_ADDRESS"
        value = var.hibiki_db_address
      }
      env {
        name  = "HIBIKI_DB_NAME"
        value = var.hibiki_db_name
      }
      env {
        name  = "HIBIKI_DB_USER"
        value = var.hibiki_db_user
      }
      env {
        name  = "HIBIKI_DB_PASSWORD"
        value = var.hibiki_db_password
      }
      env {
        name  = "HIBIKI_PUBSUB_DIALECT"
        value = var.hibiki_pubsub_dialect
      }
      env {
        name  = "HIBIKI_PUBSUB_ADDRESS"
        value = var.hibiki_pubsub_address
      }
      env {
        name  = "HIBIKI_PUBSUB_PASSWORD"
        value = var.hibiki_pubsub_password
      }
      env {
        name  = "HIBIKI_MAL_CLIENT_ID"
        value = var.hibiki_mal_client_id
      }
      env {
        name  = "HIBIKI_CRON_UPDATE_LIMIT"
        value = var.hibiki_cron_update_limit
      }
      env {
        name  = "HIBIKI_CRON_FILL_LIMIT"
        value = var.hibiki_cron_fill_limit
      }
      env {
        name  = "HIBIKI_CRON_RELEASING_AGE"
        value = var.hibiki_cron_releasing_age
      }
      env {
        name  = "HIBIKI_CRON_FINISHED_AGE"
        value = var.hibiki_cron_finished_age
      }
      env {
        name  = "HIBIKI_CRON_NOT_YET_AGE"
        value = var.hibiki_cron_not_yet_age
      }
      env {
        name  = "HIBIKI_CRON_USER_MANGA_AGE"
        value = var.hibiki_cron_user_manga_age
      }
      env {
        name  = "HIBIKI_LOG_JSON"
        value = var.hibiki_log_json
      }
      env {
        name  = "HIBIKI_LOG_LEVEL"
        value = var.hibiki_log_level
      }
      env {
        name  = "HIBIKI_NEWRELIC_LICENSE_KEY"
        value = var.hibiki_newrelic_license_key
      }
    }
  }
}

resource "google_cloud_run_service_iam_binding" "noauth" {
  service  = google_cloud_run_v2_service.server.name
  location = google_cloud_run_v2_service.server.location
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}
