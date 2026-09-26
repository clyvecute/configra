terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# 1. Enable Services
resource "google_project_service" "run_api" {
  service = "run.googleapis.com"
}

resource "google_project_service" "sql_api" {
  service = "sqladmin.googleapis.com"
}

# 2. Database (Cloud SQL Postgres)
resource "google_sql_database_instance" "instance" {
  name             = "configra-db"
  region           = var.region
  database_version = "POSTGRES_15"
  settings {
    tier = "db-f1-micro" # Free tier eligible-ish
    ip_configuration {
      ipv4_enabled = true # For simplicity, usually internal only
    }
  }
  deletion_protection  = false # For demo purposes
}

resource "google_sql_database" "database" {
  name     = "configra"
  instance = google_sql_database_instance.instance.name
}

resource "google_sql_user" "users" {
  name     = var.db_user
  instance = google_sql_database_instance.instance.name
  password = var.db_password
}

# 3. Cloud Run Service
resource "google_cloud_run_service" "default" {
  name     = "configra-api"
  location = var.region

  template {
    spec {
      containers {
        image = var.container_image
        env {
          name  = "DB_HOST"
          value = google_sql_database_instance.instance.public_ip_address
        }
        env {
          name  = "DB_USER"
          value = google_sql_user.users.name
        }
        env {
          name  = "DB_PASSWORD"
          value = google_sql_user.users.password
        }
        env {
          name  = "DB_NAME"
          value = google_sql_database.database.name
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
  
  depends_on = [google_project_service.run_api]
}

# 4. Public Access (IAM)
data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.default.location
  project     = google_cloud_run_service.default.project
  service     = google_cloud_run_service.default.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

output "url" {
  value = google_cloud_run_service.default.status[0].url
}
