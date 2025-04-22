variable "instance_count" {
  description = "Número de instancias para el cluster de CockroachDB"
  type        = number
}

variable "machine_type" {
  description = "Tipo de máquina para las instancias de CockroachDB"
  type        = string
}

variable "zone" {
  description = "Zona GCP para las instancias"
  type        = string
}

variable "disk_size_gb" {
  description = "Tamaño del disco en GB para cada instancia"
  type        = number
}

variable "subnetwork_id" {
  description = "ID de la subred donde se desplegarán las instancias"
  type        = string
}

variable "service_account_email" {
  description = "Email de la cuenta de servicio para las instancias"
  type        = string
}