variable "project_id" {
  description = "ID del proyecto GCP"
  type        = string
}

variable "region" {
  description = "Región GCP para los recursos"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "Zona GCP para los recursos zonales"
  type        = string
  default     = "us-central1-a"
}

variable "cockroachdb_instance_count" {
  description = "Número de instancias para el cluster de CockroachDB"
  type        = number
  default     = 3
}

variable "cockroachdb_machine_type" {
  description = "Tipo de máquina para las instancias de CockroachDB"
  type        = string
  default     = "e2-standard-2"
}

variable "cockroachdb_disk_size_gb" {
  description = "Tamaño del disco en GB para las instancias de CockroachDB"
  type        = number
  default     = 50
}

variable "gke_node_count" {
  description = "Número de nodos para el cluster GKE"
  type        = number
  default     = 2
}

variable "gke_machine_type" {
  description = "Tipo de máquina para los nodos de GKE"
  type        = string
  default     = "e2-standard-2"
}


variable "backend_repo" {
  description = "URL del repositorio GitHub para el backend"
  type        = string
}

variable "backend_image" {
  description = "Nombre de la imagen de contenedor para el backend"
  type        = string
  default     = "golang-backend"
}