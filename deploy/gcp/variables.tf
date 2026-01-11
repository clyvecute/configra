variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region"
  default     = "us-central1"
}

variable "db_user" {
  description = "Database User"
  default     = "configra_user"
}

variable "db_password" {
  description = "Database Password"
  sensitive   = true
}

variable "container_image" {
  description = "Docker image to deploy"
}
