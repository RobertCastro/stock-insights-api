variable "project_id" {
  description = "ID del proyecto en GCP"
  type        = string
}

variable "region" {
  description = "Región de GCP donde se desplegarán los recursos"
  type        = string
}

variable "cluster_name" {
  description = "Nombre del cluster de GKE"
  type        = string
  default     = "stock-insights-cluster"
}

variable "db_name" {
  description = "Nombre de la base de datos"
  type        = string
  default     = "stockdb"
}

variable "db_user" {
  description = "Usuario de la base de datos"
  type        = string
  default     = "stockuser"
}

variable "db_password" {
  description = "Contraseña de la base de datos"
  type        = string
  sensitive   = true
}